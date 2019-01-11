/* **********************************************************
 * Copyright 2018 VMware, Inc.  All rights reserved. -- VMware Confidential
 * **********************************************************/

package virtualmachineservice

import (
	"context"
	"fmt"
	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"time"
	"vmware.com/kubevsphere/pkg"
	"vmware.com/kubevsphere/pkg/lib"

	"github.com/kubernetes-incubator/apiserver-builder/pkg/builders"

	"vmware.com/kubevsphere/pkg/apis/vmoperator/v1beta1"
	clientSet "vmware.com/kubevsphere/pkg/client/clientset_generated/clientset"
	vmclientSet "vmware.com/kubevsphere/pkg/client/clientset_generated/clientset/typed/vmoperator/v1beta1"
	listers "vmware.com/kubevsphere/pkg/client/listers_generated/vmoperator/v1beta1"
	"vmware.com/kubevsphere/pkg/controller/sharedinformers"
)

const (
	ServiceOwnerRefKind    = "VirtualMachineService"
	ServiceOwnerRefVersion = pkg.VmOperatorKey
)

// +controller:group=vmoperator,version=v1beta1,kind=VirtualMachineService,resource=virtualmachineservices
type VirtualMachineServiceControllerImpl struct {
	builders.DefaultControllerFns

	informers *sharedinformers.SharedInformers

	// lister indexes properties about VirtualMachineService
	vmServiceLister listers.VirtualMachineServiceLister
	vmLister listers.VirtualMachineLister
	serviceLister corev1listers.ServiceLister
	endpointsLister corev1listers.EndpointsLister

	coreClientSet corev1client.CoreV1Interface
	clientSet clientSet.Interface
	vmServiceClientSet vmclientSet.VirtualMachineServiceInterface
}

func (c *VirtualMachineServiceControllerImpl) ServiceToVirtualMachineService(i interface{}) (string, error) {
	service, _ := i.(*corev1.Service)
	glog.Infof("Service update: %v", service.Name)
	if len(service.OwnerReferences) == 1 && service.OwnerReferences[0].Kind == ServiceOwnerRefKind {
		return service.Namespace + "/" + service.OwnerReferences[0].Name, nil
	} else {
		// The service is not owned
		return "", nil
	}
}

func (c *VirtualMachineServiceControllerImpl) EndpointsToVirtualMachineService(i interface{}) (string, error) {
	endpoints, _ := i.(*corev1.Endpoints)
	glog.V(4).Infof("Endpoints update: %v", endpoints.Name)
	if len(endpoints.OwnerReferences) == 1 && endpoints.OwnerReferences[0].Kind == ServiceOwnerRefKind {
		return endpoints.Namespace + "/" + endpoints.OwnerReferences[0].Name, nil
	} else {
		// The service is not owned
		return "", nil
	}
}

// Init initializes the controller and is called by the generated code
// Register watches for additional resource types here.
func (c *VirtualMachineServiceControllerImpl) Init(arguments sharedinformers.ControllerInitArguments) {

	c.informers = arguments.GetSharedInformers()

	clientSet, err := clientSet.NewForConfig(arguments.GetRestConfig())
	if err != nil {
		glog.Fatalf("Failed to create the virtual machine service client: %v", err)
	}
	c.clientSet = clientSet

	c.coreClientSet = arguments.GetSharedInformers().KubernetesClientSet.CoreV1()

	c.vmServiceClientSet = clientSet.VmoperatorV1beta1().VirtualMachineServices(corev1.NamespaceDefault)

	vmOperator := arguments.GetSharedInformers().Factory.Vmoperator().V1beta1()
	c.vmServiceLister = vmOperator.VirtualMachineServices().Lister()
	c.vmLister = vmOperator.VirtualMachines().Lister()

	services := arguments.GetSharedInformers().KubernetesFactory.Core().V1().Services()
	c.serviceLister = services.Lister()

	endpoints := arguments.GetSharedInformers().KubernetesFactory.Core().V1().Endpoints()
	c.endpointsLister = endpoints.Lister()

	arguments.Watch("Service", services.Informer(), c.ServiceToVirtualMachineService)
	arguments.Watch("Endpoint", endpoints.Informer(), c.EndpointsToVirtualMachineService)
}

// Reconcile handles enqueued messages
func (c *VirtualMachineServiceControllerImpl) Reconcile(vmService *v1beta1.VirtualMachineService) error {
	glog.V(0).Infof("Running reconcile VirtualMachineService for %s\n", vmService.Name)

	startTime := time.Now()
	defer func() {
		glog.V(0).Infof("Finished reconcile VirtualMachineService %q (%v)", vmService.Name, time.Since(startTime))
	}()

	// We hold a Finalizer on the VM Service, so it must be present
	if !vmService.ObjectMeta.DeletionTimestamp.IsZero() {
		// This VM has been deleted, sync with backend
		glog.Infof("Deletion timestamp is non-zero")

		// Noop if our finalizer is not present
		//if u.ObjectMeta.Finalizers()
		if !lib.Contains(vmService.ObjectMeta.Finalizers, v1beta1.VirtualMachineServiceFinalizer) {
			glog.Infof("reconciling virtual machine service object %v causes a no-op as there is no finalizer.", vmService.Name)
			return nil
		}

		glog.Infof("reconciling virtual machine service object %v triggers delete.", vmService.Name)
		if err := c.processVmDeletion(vmService); err != nil {
			glog.Errorf("Error deleting virtual machine service object %v; %v", vmService.Name, err)
			return err
		}

		// Remove finalizer on successful deletion.
		glog.Infof("virtual machine service object %v deletion successful, removing finalizer.", vmService.Name)
		vmService.ObjectMeta.Finalizers = lib.Filter(vmService.ObjectMeta.Finalizers, v1beta1.VirtualMachineServiceFinalizer)
		if _, err := c.vmServiceClientSet.Update(vmService); err != nil {
			glog.Errorf("Error removing finalizer from virtual machine service object %v; %v", vmService.Name, err)
			return err
		}
		return nil
	}

	// vm service holds the latest vm service info from apiserver
	vmService, err := c.vmServiceLister.VirtualMachineServices(vmService.Namespace).Get(vmService.Name)
	if err != nil {
		glog.Infof("Unable to retrieve vm service %v from store: %v", vmService.Name, err)
		return err
	}

	_, err = c.processVmCreateOrUpdate(vmService)
	if err != nil {
		glog.Infof("Failed to process Create or Update for %s: %s", vmService.Name, err)
		return err
	}

	return err
}

func (c *VirtualMachineServiceControllerImpl) processVmDeletion(vmService *v1beta1.VirtualMachineService) error {
	glog.Infof("Process VM Service Deletion for vm service %s", vmService.Name)

	glog.V(4).Infof("Deleted VM Service%s", vmService.Name)
	return nil
}

// Process a level trigger for this VM Service.
func (c *VirtualMachineServiceControllerImpl) processVmCreateOrUpdate(vmService *v1beta1.VirtualMachineService) (*v1beta1.VirtualMachineService, error) {
	glog.Infof("Process VM Service Create or Update for vm service %s", vmService.Name)

	ns := vmService.Namespace
	ctx := context.TODO()

	_, err := c.coreClientSet.Services(ns).Get(vmService.Name, metav1.GetOptions{})
	var updated *v1beta1.VirtualMachineService
	switch {
	case err != nil:
		updated, err = c.processVmCreate(ctx, vmService)
	//case NotFound:
	//glog.Errorf("Failed to get service %s: %s", vmService.Name, err)
	default:
		updated, err = c.processVmUpdate(ctx, vmService)
	}

	return updated, err
}

func (c *VirtualMachineServiceControllerImpl) vmToEndpointAddress(vm *v1beta1.VirtualMachine) *corev1.EndpointAddress {
	return &corev1.EndpointAddress{
		IP:       vm.Status.VmIp,
		NodeName: &vm.Status.Host,
		TargetRef: &corev1.ObjectReference{
			Kind:            "VirtualMachine",
			Namespace:       vm.ObjectMeta.Namespace,
			Name:            vm.ObjectMeta.Name,
			UID:             vm.ObjectMeta.UID,
			ResourceVersion: vm.ObjectMeta.ResourceVersion,
		}}
}

func (c *VirtualMachineServiceControllerImpl) vmServiceToService(vmService *v1beta1.VirtualMachineService) *corev1.Service {

	t := true
	om := metav1.ObjectMeta{
		Namespace: vmService.GetNamespace(),
		Name:      vmService.GetName(),
		Labels:    vmService.GetLabels(),
		Annotations: vmService.GetAnnotations(),
		OwnerReferences: []metav1.OwnerReference{
			metav1.OwnerReference{
				UID:                vmService.GetUID(),
				Name:               vmService.GetName(),
				Controller:         &t,
				BlockOwnerDeletion: &t,
				Kind:               ServiceOwnerRefKind,
				APIVersion:         ServiceOwnerRefVersion,
			},
		},
	}

	servicePorts := []corev1.ServicePort{}
	for _, vmPort := range vmService.Spec.Ports {
		sport := corev1.ServicePort{
			Name: vmPort.Name,
			Protocol: corev1.Protocol(vmPort.Protocol),
			Port: vmPort.Port,
			TargetPort: intstr.FromInt(int(vmPort.Port)),
		}
		servicePorts = append(servicePorts, sport)
	}

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "core/v1",
		},
		ObjectMeta: om,
		Spec: corev1.ServiceSpec{
			Selector: vmService.Spec.Selector,
			Type:     corev1.ServiceTypeClusterIP, // TODO: Pull this from VM Service
			Ports:    servicePorts,
		},
	}
}

func FindPort(vm *v1beta1.VirtualMachine, svcPort *corev1.ServicePort) (int, error) {
	portName := svcPort.TargetPort
	switch portName.Type {
	case intstr.String:
		name := portName.StrVal
		for _, port := range vm.Spec.Ports {
			if port.Name == name && port.Protocol == svcPort.Protocol {
				return int(port.Port), nil
			}
		}
	case intstr.Int:
		return portName.IntValue(), nil
	}

	return 0, fmt.Errorf("no suitable port for manifest: %s", vm.UID)
}

func addEndpointSubset(subsets []corev1.EndpointSubset, vm *v1beta1.VirtualMachine, epa corev1.EndpointAddress, epp *corev1.EndpointPort) ([]corev1.EndpointSubset) {
	ports := []corev1.EndpointPort{}
	if epp != nil {
		ports = append(ports, *epp)
	}

	subsets = append(subsets,
		corev1.EndpointSubset{
			Addresses: []corev1.EndpointAddress{epa},
			Ports:     ports,
		})

	return subsets
}

func (c *VirtualMachineServiceControllerImpl) updateService(ctx context.Context, vmService *v1beta1.VirtualMachineService, service *corev1.Service) error {
	glog.V(0).Infof("Updating service for VirtualMachineService for %s/%s\n", vmService.Namespace, vmService.Name)

	defer func() {
		glog.V(0).Infof("Finished syncing service for %q", vmService.Name)
	}()

	// See if there's actually an update here.
	currentService, err := c.serviceLister.Services(service.Namespace).Get(service.Name)
	if err != nil {
		if errors.IsNotFound(err) {
			currentService = service
		} else {
			glog.Errorf("Failed to list services: %s", err)
			return err
		}
	}

	createService := len(currentService.ResourceVersion) == 0

	newService := currentService.DeepCopy()
	newService.Labels = service.Labels
	if newService.Annotations == nil {
		newService.Annotations = make(map[string]string)
	}

	if createService {
		// No previous service, create one
		glog.V(0).Infof("Create service %s for %v/%v", newService, service.Namespace, service.Name)
		_, err = c.coreClientSet.Services(service.Namespace).Create(newService)
	} else {
		// Pre-existing
		glog.V(0).Infof("Update service %s for %v/%v", newService, service.Namespace, service.Name)
		_, err = c.coreClientSet.Services(service.Namespace).Update(newService)
	}

	return err
}

func (c *VirtualMachineServiceControllerImpl) updateEndpoints(ctx context.Context, vmService *v1beta1.VirtualMachineService, service *corev1.Service) error {
	glog.V(0).Infof("Updating endponts for VirtualMachineService for %s/%s\n", vmService.Namespace, vmService.Name)

	defer func() {
		glog.V(0).Infof("Finished syncing endpoints for %q", vmService.Name)
	}()

	vms, err := c.vmLister.VirtualMachines(vmService.Namespace).List(labels.Set(service.Spec.Selector).AsSelectorPreValidated())
	if err != nil {
		return err
	}

	glog.V(0).Infof("VMs: %s match labels %s", vms, service.Spec.Selector)

	// Determine if any endpoints match
	subsets := []corev1.EndpointSubset{}

	for _, vm := range vms {
		glog.V(0).Infof("Resolving ports for vm %s/%s", vm.Namespace, vm.Name)
		// Handle multiple VM interfaces
		if len(vm.Status.VmIp) == 0 {
			glog.V(0).Infof("Failed to find an IP for vm %s/%s", vm.Namespace, vm.Name)
			continue
		}

		// Ignore if all requried values aren't present
		if vm.Status.Host == "" {
			glog.Infof("Skipping vm %s/%s due to empty host", vm.Namespace, vm.Name)
			continue
		}

		epa := *c.vmToEndpointAddress(vm)

		// TODO: Headless support
		for i := range service.Spec.Ports {
			glog.V(0).Infof("Resolving service ports: %d", i)
			servicePort := &service.Spec.Ports[i]

			portName := servicePort.Name
			portProto := servicePort.Protocol

			glog.Infof("Port for VM %s: %s %s", vm.Name, portName, portProto)

			portNum, err := FindPort(vm, servicePort)
			if err != nil {
				glog.V(4).Infof("Failed to find port for service %s/%s: %v", service.Namespace, service.Name, err)
				continue
			}

			epp := &corev1.EndpointPort{Name: portName, Port: int32(portNum), Protocol: portProto}
			subsets = addEndpointSubset(subsets, vm, epa, epp)
		}

		// See if there's actually an update here.
		currentEndpoints, err := c.endpointsLister.Endpoints(service.Namespace).Get(service.Name)
		if err != nil {
			if errors.IsNotFound(err) {
				currentEndpoints = &corev1.Endpoints{
					ObjectMeta: metav1.ObjectMeta{
						Name:   service.Name,
						Labels: service.Labels,
					},
				}
			} else {
				glog.Errorf("Failed to list services: %s", err)
				return err
			}
		}

		createEndpoints := len(currentEndpoints.ResourceVersion) == 0

		newEndpoints := currentEndpoints.DeepCopy()
		newEndpoints.Subsets = subsets
		newEndpoints.Labels = service.Labels
		if newEndpoints.Annotations == nil {
			newEndpoints.Annotations = make(map[string]string)
		}

		if createEndpoints {
			// No previous endpoints, create them
			glog.V(0).Infof("Create endpoints %s for %v/%v", newEndpoints, vmService.Namespace, vmService.Name)
			_, err = c.coreClientSet.Endpoints(service.Namespace).Create(newEndpoints)
		} else {
			// Pre-existing
			glog.V(0).Infof("Update endpoints %s for %v/%v", newEndpoints, vmService.Namespace, vmService.Name)
			_, err = c.coreClientSet.Endpoints(service.Namespace).Update(newEndpoints)
		}
		if err != nil {
			glog.Errorf("Failed to create endpoints: %s", err)
			if createEndpoints && errors.IsForbidden(err) {
				// A request is forbidden primarily for two reasons:
				// 1. namespace is terminating, endpoint creation is not allowed by default.
				// 2. policy is misconfigured, in which case no service would function anywhere.
				// Given the frequency of 1, we log at a lower level.
				glog.V(5).Infof("Forbidden from creating endpoints: %v", err)
			}
			return err
		}
		return nil

	}

	return nil
}

// Process a create event for a new VM.
func (c *VirtualMachineServiceControllerImpl) processVmCreate(ctx context.Context, vmService *v1beta1.VirtualMachineService) (*v1beta1.VirtualMachineService, error) {
	glog.Infof("Creating VM Service: %s", vmService.Name)
	// Create Service
	service := c.vmServiceToService(vmService)

	err := c.updateService(ctx, vmService, service)
	if err != nil {
		glog.Errorf("Failed to update Service for %s/%s: %s", vmService.Namespace, vmService.Name, err)
		return nil, err
	}

	err = c.updateEndpoints(ctx, vmService, service)
	if err != nil {
		glog.Errorf("Failed to update Endpoints for %s/%s: %s", vmService.Namespace, vmService.Name, err)
		return nil, err
	}

	return nil, nil
}

// Process an update event for an existing VM.
func (c *VirtualMachineServiceControllerImpl) processVmUpdate(ctx context.Context, vmService *v1beta1.VirtualMachineService) (*v1beta1.VirtualMachineService, error) {
	glog.Infof("Updating VM Service: %s", vmService.Name)
	// Ensure Service and Endpoints are correct
	// Determine if Service matches any VMs
	vms, err := c.vmLister.VirtualMachines(vmService.Namespace).List(labels.Set(vmService.Spec.Selector).AsSelectorPreValidated())
	if err != nil {
		// Since we're getting stuff from a local cache, it is
		// basically impossible to get this error.
		return nil, err
	}

	for _, vm := range vms {
		glog.Infof("VM %s/%s matched with labels: %s", vm.Namespace, vm.Name, vm.ObjectMeta.Labels)
	}

	return nil, nil
}

func (c *VirtualMachineServiceControllerImpl) Get(namespace, name string) (*v1beta1.VirtualMachineService, error) {
	return c.vmServiceLister.VirtualMachineServices(namespace).Get(name)
}