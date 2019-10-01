// +build !ignore_autogenerated

/* **********************************************************
 * Copyright 2019 VMware, Inc.  All rights reserved. -- VMware Confidential
 * **********************************************************/

// Code generated by deepcopy-gen. DO NOT EDIT.

package vmoperator

import (
	v1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachine) DeepCopyInto(out *VirtualMachine) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachine.
func (in *VirtualMachine) DeepCopy() *VirtualMachine {
	if in == nil {
		return nil
	}
	out := new(VirtualMachine)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachine) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineClass) DeepCopyInto(out *VirtualMachineClass) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineClass.
func (in *VirtualMachineClass) DeepCopy() *VirtualMachineClass {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineClass)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachineClass) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineClassHardware) DeepCopyInto(out *VirtualMachineClassHardware) {
	*out = *in
	out.Memory = in.Memory.DeepCopy()
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineClassHardware.
func (in *VirtualMachineClassHardware) DeepCopy() *VirtualMachineClassHardware {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineClassHardware)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineClassList) DeepCopyInto(out *VirtualMachineClassList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VirtualMachineClass, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineClassList.
func (in *VirtualMachineClassList) DeepCopy() *VirtualMachineClassList {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineClassList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachineClassList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineClassPolicies) DeepCopyInto(out *VirtualMachineClassPolicies) {
	*out = *in
	in.Resources.DeepCopyInto(&out.Resources)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineClassPolicies.
func (in *VirtualMachineClassPolicies) DeepCopy() *VirtualMachineClassPolicies {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineClassPolicies)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineClassResourceSpec) DeepCopyInto(out *VirtualMachineClassResourceSpec) {
	*out = *in
	out.Cpu = in.Cpu.DeepCopy()
	out.Memory = in.Memory.DeepCopy()
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineClassResourceSpec.
func (in *VirtualMachineClassResourceSpec) DeepCopy() *VirtualMachineClassResourceSpec {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineClassResourceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineClassResources) DeepCopyInto(out *VirtualMachineClassResources) {
	*out = *in
	in.Requests.DeepCopyInto(&out.Requests)
	in.Limits.DeepCopyInto(&out.Limits)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineClassResources.
func (in *VirtualMachineClassResources) DeepCopy() *VirtualMachineClassResources {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineClassResources)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineClassSpec) DeepCopyInto(out *VirtualMachineClassSpec) {
	*out = *in
	in.Hardware.DeepCopyInto(&out.Hardware)
	in.Policies.DeepCopyInto(&out.Policies)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineClassSpec.
func (in *VirtualMachineClassSpec) DeepCopy() *VirtualMachineClassSpec {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineClassSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineClassStatus) DeepCopyInto(out *VirtualMachineClassStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineClassStatus.
func (in *VirtualMachineClassStatus) DeepCopy() *VirtualMachineClassStatus {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineClassStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineCondition) DeepCopyInto(out *VirtualMachineCondition) {
	*out = *in
	in.LastProbeTime.DeepCopyInto(&out.LastProbeTime)
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineCondition.
func (in *VirtualMachineCondition) DeepCopy() *VirtualMachineCondition {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineCondition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineImage) DeepCopyInto(out *VirtualMachineImage) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineImage.
func (in *VirtualMachineImage) DeepCopy() *VirtualMachineImage {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineImage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachineImage) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineImageList) DeepCopyInto(out *VirtualMachineImageList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VirtualMachineImage, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineImageList.
func (in *VirtualMachineImageList) DeepCopy() *VirtualMachineImageList {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineImageList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachineImageList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineImageSpec) DeepCopyInto(out *VirtualMachineImageSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineImageSpec.
func (in *VirtualMachineImageSpec) DeepCopy() *VirtualMachineImageSpec {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineImageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineImageStatus) DeepCopyInto(out *VirtualMachineImageStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineImageStatus.
func (in *VirtualMachineImageStatus) DeepCopy() *VirtualMachineImageStatus {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineImageStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineList) DeepCopyInto(out *VirtualMachineList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VirtualMachine, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineList.
func (in *VirtualMachineList) DeepCopy() *VirtualMachineList {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachineList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineMetadata) DeepCopyInto(out *VirtualMachineMetadata) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineMetadata.
func (in *VirtualMachineMetadata) DeepCopy() *VirtualMachineMetadata {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineMetadata)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineNetworkInterface) DeepCopyInto(out *VirtualMachineNetworkInterface) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineNetworkInterface.
func (in *VirtualMachineNetworkInterface) DeepCopy() *VirtualMachineNetworkInterface {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineNetworkInterface)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachinePort) DeepCopyInto(out *VirtualMachinePort) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachinePort.
func (in *VirtualMachinePort) DeepCopy() *VirtualMachinePort {
	if in == nil {
		return nil
	}
	out := new(VirtualMachinePort)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineService) DeepCopyInto(out *VirtualMachineService) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineService.
func (in *VirtualMachineService) DeepCopy() *VirtualMachineService {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineService)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachineService) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineServiceList) DeepCopyInto(out *VirtualMachineServiceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]VirtualMachineService, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineServiceList.
func (in *VirtualMachineServiceList) DeepCopy() *VirtualMachineServiceList {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineServiceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *VirtualMachineServiceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineServicePort) DeepCopyInto(out *VirtualMachineServicePort) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineServicePort.
func (in *VirtualMachineServicePort) DeepCopy() *VirtualMachineServicePort {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineServicePort)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineServiceSpec) DeepCopyInto(out *VirtualMachineServiceSpec) {
	*out = *in
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]VirtualMachineServicePort, len(*in))
		copy(*out, *in)
	}
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineServiceSpec.
func (in *VirtualMachineServiceSpec) DeepCopy() *VirtualMachineServiceSpec {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineServiceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineServiceStatus) DeepCopyInto(out *VirtualMachineServiceStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineServiceStatus.
func (in *VirtualMachineServiceStatus) DeepCopy() *VirtualMachineServiceStatus {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineServiceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineSpec) DeepCopyInto(out *VirtualMachineSpec) {
	*out = *in
	in.Env.DeepCopyInto(&out.Env)
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make([]VirtualMachinePort, len(*in))
		copy(*out, *in)
	}
	if in.VmMetadata != nil {
		in, out := &in.VmMetadata, &out.VmMetadata
		*out = new(VirtualMachineMetadata)
		**out = **in
	}
	if in.NetworkInterfaces != nil {
		in, out := &in.NetworkInterfaces, &out.NetworkInterfaces
		*out = make([]VirtualMachineNetworkInterface, len(*in))
		copy(*out, *in)
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]VirtualMachineVolumes, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineSpec.
func (in *VirtualMachineSpec) DeepCopy() *VirtualMachineSpec {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineStatus) DeepCopyInto(out *VirtualMachineStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]VirtualMachineCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]VirtualMachineVolumeStatus, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineStatus.
func (in *VirtualMachineStatus) DeepCopy() *VirtualMachineStatus {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineVolumeStatus) DeepCopyInto(out *VirtualMachineVolumeStatus) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineVolumeStatus.
func (in *VirtualMachineVolumeStatus) DeepCopy() *VirtualMachineVolumeStatus {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineVolumeStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VirtualMachineVolumes) DeepCopyInto(out *VirtualMachineVolumes) {
	*out = *in
	if in.PersistentVolumeClaim != nil {
		in, out := &in.PersistentVolumeClaim, &out.PersistentVolumeClaim
		*out = new(v1.PersistentVolumeClaimVolumeSource)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VirtualMachineVolumes.
func (in *VirtualMachineVolumes) DeepCopy() *VirtualMachineVolumes {
	if in == nil {
		return nil
	}
	out := new(VirtualMachineVolumes)
	in.DeepCopyInto(out)
	return out
}
