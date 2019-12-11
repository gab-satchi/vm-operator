/* **********************************************************
 * Copyright 2019 VMware, Inc.  All rights reserved. -- VMware Confidential
 * **********************************************************/

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	internalinterfaces "github.com/vmware-tanzu/vm-operator/pkg/client/informers_generated/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// VirtualMachines returns a VirtualMachineInformer.
	VirtualMachines() VirtualMachineInformer
	// VirtualMachineClasses returns a VirtualMachineClassInformer.
	VirtualMachineClasses() VirtualMachineClassInformer
	// VirtualMachineImages returns a VirtualMachineImageInformer.
	VirtualMachineImages() VirtualMachineImageInformer
	// VirtualMachineServices returns a VirtualMachineServiceInformer.
	VirtualMachineServices() VirtualMachineServiceInformer
	// VirtualMachineSetResourcePolicies returns a VirtualMachineSetResourcePolicyInformer.
	VirtualMachineSetResourcePolicies() VirtualMachineSetResourcePolicyInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// VirtualMachines returns a VirtualMachineInformer.
func (v *version) VirtualMachines() VirtualMachineInformer {
	return &virtualMachineInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// VirtualMachineClasses returns a VirtualMachineClassInformer.
func (v *version) VirtualMachineClasses() VirtualMachineClassInformer {
	return &virtualMachineClassInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// VirtualMachineImages returns a VirtualMachineImageInformer.
func (v *version) VirtualMachineImages() VirtualMachineImageInformer {
	return &virtualMachineImageInformer{factory: v.factory, tweakListOptions: v.tweakListOptions}
}

// VirtualMachineServices returns a VirtualMachineServiceInformer.
func (v *version) VirtualMachineServices() VirtualMachineServiceInformer {
	return &virtualMachineServiceInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}

// VirtualMachineSetResourcePolicies returns a VirtualMachineSetResourcePolicyInformer.
func (v *version) VirtualMachineSetResourcePolicies() VirtualMachineSetResourcePolicyInformer {
	return &virtualMachineSetResourcePolicyInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}
