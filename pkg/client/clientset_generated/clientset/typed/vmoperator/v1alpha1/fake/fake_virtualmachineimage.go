/* **********************************************************
 * Copyright 2019 VMware, Inc.  All rights reserved. -- VMware Confidential
 * **********************************************************/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/vmware-tanzu/vm-operator/pkg/apis/vmoperator/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeVirtualMachineImages implements VirtualMachineImageInterface
type FakeVirtualMachineImages struct {
	Fake *FakeVmoperatorV1alpha1
}

var virtualmachineimagesResource = schema.GroupVersionResource{Group: "vmoperator.vmware.com", Version: "v1alpha1", Resource: "virtualmachineimages"}

var virtualmachineimagesKind = schema.GroupVersionKind{Group: "vmoperator.vmware.com", Version: "v1alpha1", Kind: "VirtualMachineImage"}

// Get takes name of the virtualMachineImage, and returns the corresponding virtualMachineImage object, and an error if there is any.
func (c *FakeVirtualMachineImages) Get(name string, options v1.GetOptions) (result *v1alpha1.VirtualMachineImage, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(virtualmachineimagesResource, name), &v1alpha1.VirtualMachineImage{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VirtualMachineImage), err
}

// List takes label and field selectors, and returns the list of VirtualMachineImages that match those selectors.
func (c *FakeVirtualMachineImages) List(opts v1.ListOptions) (result *v1alpha1.VirtualMachineImageList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(virtualmachineimagesResource, virtualmachineimagesKind, opts), &v1alpha1.VirtualMachineImageList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.VirtualMachineImageList{ListMeta: obj.(*v1alpha1.VirtualMachineImageList).ListMeta}
	for _, item := range obj.(*v1alpha1.VirtualMachineImageList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested virtualMachineImages.
func (c *FakeVirtualMachineImages) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(virtualmachineimagesResource, opts))
}

// Create takes the representation of a virtualMachineImage and creates it.  Returns the server's representation of the virtualMachineImage, and an error, if there is any.
func (c *FakeVirtualMachineImages) Create(virtualMachineImage *v1alpha1.VirtualMachineImage) (result *v1alpha1.VirtualMachineImage, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(virtualmachineimagesResource, virtualMachineImage), &v1alpha1.VirtualMachineImage{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VirtualMachineImage), err
}

// Update takes the representation of a virtualMachineImage and updates it. Returns the server's representation of the virtualMachineImage, and an error, if there is any.
func (c *FakeVirtualMachineImages) Update(virtualMachineImage *v1alpha1.VirtualMachineImage) (result *v1alpha1.VirtualMachineImage, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(virtualmachineimagesResource, virtualMachineImage), &v1alpha1.VirtualMachineImage{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VirtualMachineImage), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeVirtualMachineImages) UpdateStatus(virtualMachineImage *v1alpha1.VirtualMachineImage) (*v1alpha1.VirtualMachineImage, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(virtualmachineimagesResource, "status", virtualMachineImage), &v1alpha1.VirtualMachineImage{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VirtualMachineImage), err
}

// Delete takes name of the virtualMachineImage and deletes it. Returns an error if one occurs.
func (c *FakeVirtualMachineImages) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(virtualmachineimagesResource, name), &v1alpha1.VirtualMachineImage{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeVirtualMachineImages) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(virtualmachineimagesResource, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.VirtualMachineImageList{})
	return err
}

// Patch applies the patch and returns the patched virtualMachineImage.
func (c *FakeVirtualMachineImages) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.VirtualMachineImage, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(virtualmachineimagesResource, name, data, subresources...), &v1alpha1.VirtualMachineImage{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.VirtualMachineImage), err
}
