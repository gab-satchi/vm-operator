/* **********************************************************
 * Copyright 2019 VMware, Inc.  All rights reserved. -- VMware Confidential
 * **********************************************************/

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	vmoperatorv1alpha1 "github.com/vmware-tanzu/vm-operator/pkg/apis/vmoperator/v1alpha1"
	clientset "github.com/vmware-tanzu/vm-operator/pkg/client/clientset_generated/clientset"
	internalinterfaces "github.com/vmware-tanzu/vm-operator/pkg/client/informers_generated/externalversions/internalinterfaces"
	v1alpha1 "github.com/vmware-tanzu/vm-operator/pkg/client/listers_generated/vmoperator/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// VirtualMachineInformer provides access to a shared informer and lister for
// VirtualMachines.
type VirtualMachineInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.VirtualMachineLister
}

type virtualMachineInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewVirtualMachineInformer constructs a new informer for VirtualMachine type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewVirtualMachineInformer(client clientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredVirtualMachineInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredVirtualMachineInformer constructs a new informer for VirtualMachine type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredVirtualMachineInformer(client clientset.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.VmoperatorV1alpha1().VirtualMachines(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.VmoperatorV1alpha1().VirtualMachines(namespace).Watch(options)
			},
		},
		&vmoperatorv1alpha1.VirtualMachine{},
		resyncPeriod,
		indexers,
	)
}

func (f *virtualMachineInformer) defaultInformer(client clientset.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredVirtualMachineInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *virtualMachineInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&vmoperatorv1alpha1.VirtualMachine{}, f.defaultInformer)
}

func (f *virtualMachineInformer) Lister() v1alpha1.VirtualMachineLister {
	return v1alpha1.NewVirtualMachineLister(f.Informer().GetIndexer())
}
