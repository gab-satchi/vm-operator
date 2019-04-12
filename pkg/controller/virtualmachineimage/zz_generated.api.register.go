/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file was autogenerated by apiregister-gen. Do not edit it manually!

package virtualmachineimage

import (
	"github.com/kubernetes-incubator/apiserver-builder-alpha/pkg/controller"
	"github.com/vmware-tanzu/vm-operator/pkg/controller/sharedinformers"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog"
)

// VirtualMachineImageController implements the controller.VirtualMachineImageController interface
type VirtualMachineImageController struct {
	queue *controller.QueueWorker

	// Handles messages
	controller *VirtualMachineImageControllerImpl

	Name string

	BeforeReconcile func(key string)
	AfterReconcile  func(key string, err error)

	Informers *sharedinformers.SharedInformers
}

// NewController returns a new VirtualMachineImageController for responding to VirtualMachineImage events
func NewVirtualMachineImageController(config *rest.Config, si *sharedinformers.SharedInformers) *VirtualMachineImageController {
	q := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "VirtualMachineImage")

	queue := &controller.QueueWorker{q, 10, "VirtualMachineImage", nil}
	c := &VirtualMachineImageController{queue, nil, "VirtualMachineImage", nil, nil, si}

	// For non-generated code to add events
	uc := &VirtualMachineImageControllerImpl{}
	var ci sharedinformers.Controller = uc

	// Call the Init method that is implemented.
	// Support multiple Init methods for backwards compatibility
	if i, ok := ci.(sharedinformers.LegacyControllerInit); ok {
		i.Init(config, si, c.LookupAndReconcile)
	} else if i, ok := ci.(sharedinformers.ControllerInit); ok {
		i.Init(&sharedinformers.ControllerInitArgumentsImpl{si, config, c.LookupAndReconcile})
	}

	c.controller = uc

	queue.Reconcile = c.reconcile
	if c.Informers.WorkerQueues == nil {
		c.Informers.WorkerQueues = map[string]*controller.QueueWorker{}
	}
	c.Informers.WorkerQueues["VirtualMachineImage"] = queue
	si.Factory.Vmoperator().V1alpha1().VirtualMachineImages().Informer().
		AddEventHandler(&controller.QueueingEventHandler{q, nil, false})
	return c
}

func (c *VirtualMachineImageController) GetName() string {
	return c.Name
}

func (c *VirtualMachineImageController) LookupAndReconcile(key string) (err error) {
	return c.reconcile(key)
}

func (c *VirtualMachineImageController) reconcile(key string) (err error) {
	var namespace, name string

	if c.BeforeReconcile != nil {
		c.BeforeReconcile(key)
	}
	if c.AfterReconcile != nil {
		// Wrap in a function so err is evaluated after it is set
		defer func() { c.AfterReconcile(key, err) }()
	}

	namespace, name, err = cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return
	}

	u, err := c.controller.Get(namespace, name)
	if errors.IsNotFound(err) {
		klog.Infof("Not doing work for VirtualMachineImage %v because it has been deleted", key)
		// Set error so it is picked up by AfterReconcile and the return function
		err = nil
		return
	}
	if err != nil {
		klog.Errorf("Unable to retrieve VirtualMachineImage %v from store: %v", key, err)
		return
	}

	// Set error so it is picked up by AfterReconcile and the return function
	err = c.controller.Reconcile(u)

	return
}

func (c *VirtualMachineImageController) Run(stopCh <-chan struct{}) {
	for _, q := range c.Informers.WorkerQueues {
		q.Run(stopCh)
	}
	controller.GetDefaults(c.controller).Run(stopCh)
}
