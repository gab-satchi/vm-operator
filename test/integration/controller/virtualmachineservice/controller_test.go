/* **********************************************************
 * Copyright 2018 VMware, Inc.  All rights reserved. -- VMware Confidential
 * **********************************************************/

package virtualmachineservice_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "gitlab.eng.vmware.com/iaas-platform/vm-operator/pkg/apis/vmoperator/v1alpha1"
	. "gitlab.eng.vmware.com/iaas-platform/vm-operator/pkg/client/clientset_generated/clientset/typed/vmoperator/v1alpha1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = XDescribe("VirtualMachineService controller", func() {
	var instance VirtualMachineService
	var expectedKey string
	var client VirtualMachineServiceInterface
	var before chan struct{}
	var after chan struct{}
	typ := "ClusterIP"
	port := VirtualMachineServicePort{
		Name:       "foo",
		Protocol:   "TCP",
		Port:       22,
		TargetPort: 22,
	}
	selector := map[string]string{"foo": "bar"}

	BeforeEach(func() {
		instance = VirtualMachineService{Spec: VirtualMachineServiceSpec{
			Type:     typ,
			Ports:    []VirtualMachineServicePort{port},
			Selector: selector,
		}}
		instance.Name = "instance-1"
		expectedKey = "virtualmachineservice-controller-test-handler/instance-1"
	})

	AfterEach(func() {
		_ = client.Delete(instance.Name, &metav1.DeleteOptions{})
	})

	XDescribe("when creating a new object", func() {
		It("invoke the reconcile method", func() {
			client = cs.VmoperatorV1alpha1().VirtualMachineServices("virtualmachineservice-controller-test-handler")
			before = make(chan struct{})
			after = make(chan struct{})

			actualKey := ""
			var actualErr error = nil

			// Setup test callbacks to be called when the message is reconciled
			controller.BeforeReconcile = func(key string) {
				controller.BeforeReconcile = nil
				defer close(before)
				actualKey = key
			}
			controller.AfterReconcile = func(key string, err error) {
				controller.AfterReconcile = nil
				defer close(after)
				actualKey = key
				actualErr = err
			}

			// Create an instance
			_, err := client.Create(&instance)
			Expect(err).ShouldNot(HaveOccurred())

			// Verify reconcile function is called against the correct key
			select {
			case <-before:
				Expect(actualKey).To(Equal(expectedKey))
				Expect(actualErr).ShouldNot(HaveOccurred())
			case <-time.After(time.Second * 2):
				Fail("reconcile never called")
			}

			select {
			case <-after:
				Expect(actualKey).To(Equal(expectedKey))
				Expect(actualErr).ShouldNot(HaveOccurred())
			case <-time.After(time.Second * 2):
				Fail("reconcile never finished")
			}
		})
	})
})
