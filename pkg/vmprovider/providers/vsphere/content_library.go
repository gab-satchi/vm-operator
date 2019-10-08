// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package vsphere

import (
	"context"
	"io"
	"net/url"
	"strings"
	"time"

	"github.com/vmware/govmomi/ovf"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/vapi/library"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25/soap"
)

type ContentLibraryProvider struct {
	session *Session
}

type ContentDownloadHandler interface {
	GenerateDownloadUriForLibraryItem(ctx context.Context, restClient *rest.Client, item *library.Item, sess *Session) (DownloadUriResponse, error)
}

type ContentDownloadProvider struct {
	ApiWaitTimeSecs int
}

type TimerTaskResponse struct {
	TaskDone     bool
	Err          error
	returnValues map[string]string
}

const (
	libItemName        = "itemName"
	libFileName        = "fileName"
	libItemId          = "itemId"
	libSessionId       = "sessionId"
	libFileDownloadUrl = "fileUrl"
)

func NewContentLibraryProvider(ses *Session) *ContentLibraryProvider {
	contentLibProvider := &ContentLibraryProvider{
		session: ses,
	}
	return contentLibProvider
}

type DownloadUriResponse struct {
	FileUri           string
	DownloadSessionId string
}

// ParseAndRetrievePropsFromLibraryItem downloads the supported file from content library.
// parses the downloaded ovf and retrieves the properties defined under
// VirtualSystem.Product.Property and return them as a map.
func (cs *ContentLibraryProvider) ParseAndRetrievePropsFromLibraryItem(ctx context.Context, item *library.Item, clHandler ContentDownloadHandler) (map[string]string, error) {

	var downloadedFileContent io.ReadCloser

	err := cs.session.WithRestClient(ctx, func(c *rest.Client) error {
		// download ovf from the library item
		response, err := clHandler.GenerateDownloadUriForLibraryItem(ctx, c, item, cs.session)
		if err != nil {
			return err
		}

		if isInvalidResponse(response) {
			return errors.Errorf("error occurred downloading item %v", item.Name)
		}

		defer deleteLibraryItemDownloadSession(cs.session, c, ctx, response.DownloadSessionId)

		//read the file as string once it is prepared for download
		downloadedFileContent, err = ReadFileFromUrl(ctx, c, cs.session, response.FileUri)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil || downloadedFileContent == nil {
		log.Error(err, "error occurred when downloading file from library item", libItemName, item.Name)
		return nil, err
	}

	log.Info("downloaded library item", libItemName, item.Name)
	defer downloadedFileContent.Close()

	ovfProperties, err := ParseOvfAndFetchProperties(downloadedFileContent)
	if err != nil {
		return nil, err
	}

	return ovfProperties, nil

}

func isInvalidResponse(response DownloadUriResponse) bool {
	if response == (DownloadUriResponse{}) || response.DownloadSessionId == "" || response.FileUri == "" {
		return true
	}

	return false
}

func ParseOvfAndFetchProperties(fileContent io.ReadCloser) (map[string]string, error) {
	var env *ovf.Envelope
	properties := make(map[string]string)

	env, err := ovf.Unmarshal(fileContent)
	if err != nil {
		return nil, err
	}

	for _, product := range env.VirtualSystem.Product {
		for _, prop := range product.Property {
			if strings.HasPrefix(prop.Key, "vmware-system") {
				properties[prop.Key] = *prop.Default
			}
		}
	}
	return properties, nil
}

// GenerateDownloadUriForLibraryItem downloads the file from content library in 4 stages
// 1. create a download session
// 2. list the available files and downloads only the ovf files based on filename suffix
// 3. prepare the download session and fetch the url to be used for download
// 4. download the file
func (contentSession ContentDownloadProvider) GenerateDownloadUriForLibraryItem(ctx context.Context, c *rest.Client, item *library.Item, s *Session) (DownloadUriResponse, error) {

	var (
		fileDownloadUri    string
		fileToDownload     string
		downloadSessionId  string
		clApiSleepInterval = contentSession.ApiWaitTimeSecs
	)

	libManager := library.NewManager(c)
	var info *library.DownloadFile

	// create a download session for the file referred to by item id.
	session, err := libManager.CreateLibraryItemDownloadSession(ctx, library.Session{
		LibraryItemID: item.ID,
	})
	if err != nil {
		return DownloadUriResponse{}, err
	}

	downloadSessionId = session

	//list the files available for download in the library item
	files, err := libManager.ListLibraryItemDownloadSessionFile(ctx, session)
	if err != nil {
		return DownloadUriResponse{}, err
	}

	for _, file := range files {
		log.Info("Library Item file ", libFileName, file.Name)
		fileNameParts := strings.Split(file.Name, ".")
		if IsSupportedDeployType(fileNameParts[len(fileNameParts)-1]) {
			fileToDownload = file.Name
			break
		}
	}
	log.Info("download session created", libFileName, fileToDownload, libItemId, item.ID, libSessionId, session)

	_, err = libManager.PrepareLibraryItemDownloadSessionFile(ctx, session, fileToDownload)
	if err != nil {
		return DownloadUriResponse{}, err
	}

	log.Info("request posted to prepare file", libFileName, fileToDownload, libSessionId, session)

	// content library api to prepare a file for download guarantees eventual end state of either Error or Prepared
	// in order to avoid posting too many requests to the api we are setting a sleep of 'n' seconds between each retry

	ticker := time.NewTicker(time.Duration(clApiSleepInterval) * time.Second)

	work := func(t *time.Time) (TimerTaskResponse, error) {
		returnMap := map[string]string{}
		emptyStruct := TimerTaskResponse{}
		info, err = libManager.GetLibraryItemDownloadSessionFile(ctx, session, fileToDownload)
		if err != nil {
			return emptyStruct, err
		}

		if info.Status == "ERROR" {
			return TimerTaskResponse{}, errors.Errorf("Error occurred preparing file for download %v",
				info.ErrorMessage)
		}

		if info.Status == "PREPARED" {
			log.Info("Download file", libFileDownloadUrl, info.DownloadEndpoint.URI)
			fileDownloadUri = info.DownloadEndpoint.URI
			returnMap[libFileDownloadUrl] = fileDownloadUri
			return TimerTaskResponse{
				TaskDone:     true,
				returnValues: returnMap,
			}, nil
		}

		return TimerTaskResponse{}, nil
	}

	doneChannel := make(chan TimerTaskResponse)
	go RunTaskAtInterval(doneChannel, ticker, work)
	var finalResponse = <-doneChannel
	if finalResponse.Err != nil {
		return DownloadUriResponse{}, finalResponse.Err
	}

	return DownloadUriResponse{
		FileUri:           fileDownloadUri,
		DownloadSessionId: downloadSessionId,
	}, nil

}

//accepts a channel to send the output and function containing a unit of work which takes ticker as input and returns a timerTaskResponse struct
//the work is considered complete when
// 1. a struct is returned with the value of TaskDone set (or)
// 2. there is an error returned
func RunTaskAtInterval(doneCh chan TimerTaskResponse, ticker *time.Ticker, work func(t *time.Time) (TimerTaskResponse, error)) {
	var response TimerTaskResponse
	resultCh := make(chan TimerTaskResponse)
taskLoop:
	for {
		select {
		case t := <-ticker.C:
			go func() {
				routineResponse, err := work(&t)
				if err != nil {
					resultCh <- TimerTaskResponse{
						TaskDone: true,
						Err:      err,
					}
				}
				if routineResponse.TaskDone {
					resultCh <- routineResponse
				}
			}()
		case response = <-resultCh:
			log.Info("received in result channel", "error", response.Err, "values", response.returnValues)
			ticker.Stop()
			doneCh <- response
			break taskLoop
		}
	}
}

func ReadFileFromUrl(ctx context.Context, c *rest.Client, sess *Session, fileUri string) (io.ReadCloser, error) {

	p := soap.DefaultDownload

	var readerStream io.ReadCloser

	src, err := url.Parse(fileUri)
	if err != nil {
		return nil, err
	}

	readerStream, _, err = c.Download(ctx, src, &p)
	if err != nil {
		log.Error(err, "Error occurred when downloading file", libFileDownloadUrl, src, "fileuri", fileUri)
		return nil, err
	}

	return readerStream, nil
}

func deleteLibraryItemDownloadSession(s *Session, c *rest.Client, ctx context.Context, session string) {

	sessionError := library.NewManager(c).DeleteLibraryItemDownloadSession(ctx, session)

	if sessionError != nil {
		log.Error(sessionError, "Error occurred when deleting download session", libSessionId, session)
	}

	return
}
