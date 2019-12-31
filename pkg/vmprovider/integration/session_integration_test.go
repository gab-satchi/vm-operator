// +build integration

// Copyright (c) 2019 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package integration

import (
	"context"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	vimTypes "github.com/vmware/govmomi/vim25/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/vmware-tanzu/vm-operator/pkg/apis/vmoperator/v1alpha1"
	vmoperatorv1alpha1 "github.com/vmware-tanzu/vm-operator/pkg/apis/vmoperator/v1alpha1"
	"github.com/vmware-tanzu/vm-operator/pkg/vmprovider/providers/vsphere"
	"github.com/vmware-tanzu/vm-operator/test/integration"
)

var (
	testNamespace = "test-namespace"
	testVMName    = "test-vm"
)

var _ = Describe("Sessions", func() {
	var (
		session *vsphere.Session
		err     error
		ctx     context.Context
	)
	BeforeEach(func() {
		ctx = context.Background()
		session, err = vsphere.NewSessionAndConfigure(context.TODO(), vSphereConfig, nil, nil)
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Query VM images", func() {

		Context("From Inventory - VMs", func() {

			BeforeEach(func() {
				//set source to use VM inventory
				vSphereConfig.ContentSource = ""
				err = session.ConfigureContent(context.TODO(), vSphereConfig.ContentSource)
				Expect(err).NotTo(HaveOccurred())
			})

			// TODO: The default govcsim setups 2 VM's per resource pool however we should create our own fixture for better
			// consistency and avoid failures when govcsim is updated.
			It("should list virtualmachines", func() {
				vms, err := session.ListVirtualMachines(context.TODO(), "*")
				Expect(err).NotTo(HaveOccurred())
				Expect(vms).ShouldNot(BeEmpty())
			})

			It("should get virtualmachine", func() {
				vm, err := session.GetVirtualMachine(context.TODO(), "DC0_H0_VM0")
				Expect(err).NotTo(HaveOccurred())
				Expect(vm.Name).Should(Equal("DC0_H0_VM0"))
			})
		})

		Context("From Content Library", func() {

			BeforeEach(func() {
				//set source to use CL
				vSphereConfig.ContentSource = integration.GetContentSourceID()
				err = session.ConfigureContent(context.TODO(), vSphereConfig.ContentSource)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should list virtualmachineimages from CL", func() {
				images, err := session.ListVirtualMachineImagesFromCL(context.TODO())
				Expect(err).NotTo(HaveOccurred())
				Expect(images).ShouldNot(BeEmpty())
				Expect(images[0].ObjectMeta.Name).Should(Equal("test-item"))
				Expect(images[0].Spec.Type).Should(Equal("ovf"))
			})

			It("should get virtualmachineimage from CL", func() {
				image, err := session.GetVirtualMachineImageFromCL(context.TODO(), "test-item")
				Expect(err).NotTo(HaveOccurred())
				Expect(image.ObjectMeta.Name).Should(Equal("test-item"))
				Expect(image.Spec.Type).Should(Equal("ovf"))
			})

			It("should not get virtualmachineimage from CL", func() {
				image, err := session.GetVirtualMachineImageFromCL(context.TODO(), "invalid")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).Should(Equal("failed to find image \"invalid\": no library items named: invalid"))
				Expect(image).Should(BeNil())
			})
		})
	})

	Describe("Clone VM", func() {

		BeforeEach(func() {
			//set source to use VM inventory
			vSphereConfig.ContentSource = ""
			err = session.ConfigureContent(context.TODO(), vSphereConfig.ContentSource)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("without specifying any networks in VM Spec", func() {

			It("should not override template networks", func() {
				imageName := "DC0_H0_VM0"
				vmClass := getVMClassInstance(testVMName, testNamespace)
				vm := getVirtualMachineInstance(testVMName, testNamespace, imageName, vmClass.Name)
				vmMetadata := map[string]string{}
				resVM, err := session.GetVirtualMachine(ctx, "DC0_H0_VM0")
				Expect(err).NotTo(HaveOccurred())
				nicChanges, err := session.GetNicChangeSpecs(ctx, vm, resVM)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(nicChanges)).Should(Equal(0))
				nicChanges, err = session.GetNicChangeSpecs(ctx, vm, nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(nicChanges)).Should(Equal(0))
				clonedVM, err := session.CloneVirtualMachine(ctx, vm, *vmClass, nil, vmMetadata, "foo")
				Expect(err).NotTo(HaveOccurred())
				Expect(clonedVM).ShouldNot(BeNil())
				// Existing NIF should not be changed.
				netDevices, err := clonedVM.GetNetworkDevices(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(netDevices)).Should(Equal(1))
				dev := netDevices[0].GetVirtualDevice()
				// For the vcsim env the source VM is attached to a distributed port group. Hence, the cloned VM
				// should also be attached to the same network.
				_, ok := dev.Backing.(*vimTypes.VirtualEthernetCardDistributedVirtualPortBackingInfo)
				Expect(ok).Should(BeTrue())

			})
		})

		Context("by specifying networks in VM Spec", func() {

			It("should override template networks", func() {
				imageName := "DC0_H0_VM0"
				vmClass := getVMClassInstance(testVMName, testNamespace)
				vm := getVirtualMachineInstance(testVMName+"change-net", testNamespace, imageName, vmClass.Name)
				// Add two network interfaces to the VM and attach to different networks
				vm.Spec.NetworkInterfaces = []vmoperatorv1alpha1.VirtualMachineNetworkInterface{
					{
						NetworkName: "VM Network",
					},
					{
						NetworkName:      "VM Network",
						EthernetCardType: "e1000",
					},
				}
				resVM, err := session.GetVirtualMachine(ctx, "DC0_H0_VM0")
				Expect(err).NotTo(HaveOccurred())
				nicChanges, err := session.GetNicChangeSpecs(ctx, vm, resVM)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(nicChanges)).Should(Equal(3))
				numAdd := 0
				for _, changeSpec := range nicChanges {
					Expect(changeSpec.GetVirtualDeviceConfigSpec().Operation).ShouldNot(Equal(vimTypes.VirtualDeviceConfigSpecOperationEdit))
					if changeSpec.GetVirtualDeviceConfigSpec().Operation == vimTypes.VirtualDeviceConfigSpecOperationAdd {
						numAdd += 1
						continue
					}
				}
				Expect(numAdd).Should(Equal(2))
				nicChanges, err = session.GetNicChangeSpecs(ctx, vm, nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(nicChanges)).Should(Equal(2))
				for _, changeSpec := range nicChanges {
					Expect(changeSpec.GetVirtualDeviceConfigSpec().Operation).Should(Equal(vimTypes.VirtualDeviceConfigSpecOperationAdd))
				}
				vmMetadata := map[string]string{}
				clonedVM, err := session.CloneVirtualMachine(ctx, vm, *vmClass, nil, vmMetadata, "foo")
				Expect(err).NotTo(HaveOccurred())
				netDevices, err := clonedVM.GetNetworkDevices(ctx)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(netDevices)).Should(Equal(2))
				// The interface type should be default vmxnet3
				dev1, ok := netDevices[0].(*vimTypes.VirtualVmxnet3)
				Expect(ok).Should(BeTrue())
				// TODO: enhance the test to verify the moref of the network matches the name of the network in spec.
				_, ok = dev1.Backing.(*vimTypes.VirtualEthernetCardNetworkBackingInfo)
				Expect(ok).Should(BeTrue())
				// The interface type should be e1000
				dev2, ok := netDevices[1].(*vimTypes.VirtualE1000)
				// TODO: enhance the test to verify the moref of the network matches the name of the network in spec.
				_, ok = dev2.Backing.(*vimTypes.VirtualEthernetCardNetworkBackingInfo)
				Expect(ok).Should(BeTrue())
			})
		})

		Context("when a default network is specified", func() {

			BeforeEach(func() {
				var err error
				// For the vcsim env the source VM is attached to a distributed port group. Hence, we are using standard
				// vswitch port group.
				vSphereConfig.Network = "VM Network"
				//Setup new session based on the default network
				session, err = vsphere.NewSessionAndConfigure(context.TODO(), vSphereConfig, nil, nil)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should override network from the template", func() {
				imageName := "DC0_H0_VM0"
				vmClass := getVMClassInstance(testVMName, testNamespace)
				vm := getVirtualMachineInstance(testVMName+"with-default-net", testNamespace, imageName, vmClass.Name)
				resVM, err := session.GetVirtualMachine(ctx, "DC0_H0_VM0")
				Expect(err).NotTo(HaveOccurred())
				nicChanges, err := session.GetNicChangeSpecs(ctx, vm, resVM)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(nicChanges)).Should(Equal(1))
				for _, changeSpec := range nicChanges {
					Expect(changeSpec.GetVirtualDeviceConfigSpec().Operation).Should(Equal(vimTypes.VirtualDeviceConfigSpecOperationEdit))
				}
				nicChanges, err = session.GetNicChangeSpecs(ctx, vm, nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(len(nicChanges)).Should(Equal(0))
				vmMetadata := map[string]string{}
				clonedVM, err := session.CloneVirtualMachine(context.TODO(), vm, *vmClass, nil, vmMetadata, "foo")
				Expect(err).NotTo(HaveOccurred())
				Expect(clonedVM).ShouldNot(BeNil())
				// Existing NIF should not be changed.
				netDevices, err := clonedVM.GetNetworkDevices(context.TODO())
				Expect(err).NotTo(HaveOccurred())
				Expect(len(netDevices)).Should(Equal(1))
				dev := netDevices[0].GetVirtualDevice()
				// TODO: enhance the test to verify the moref of the network matches the default network.
				_, ok := dev.Backing.(*vimTypes.VirtualEthernetCardNetworkBackingInfo)
				Expect(ok).Should(BeTrue())

			})

			It("should not override networks specified in VM Spec ", func() {
				imageName := "DC0_H0_VM0"
				vmClass := getVMClassInstance(testVMName, testNamespace)
				vm := getVirtualMachineInstance(testVMName+"change-default-net", testNamespace, imageName, vmClass.Name)
				// Add two network interfaces to the VM and attach to different networks
				vm.Spec.NetworkInterfaces = []vmoperatorv1alpha1.VirtualMachineNetworkInterface{
					{
						NetworkName: "DC0_DVPG0",
					},
					{
						NetworkName:      "DC0_DVPG0",
						EthernetCardType: "e1000",
					},
				}
				vmMetadata := map[string]string{}
				clonedVM, err := session.CloneVirtualMachine(context.TODO(), vm, *vmClass, nil, vmMetadata, "foo")
				Expect(err).NotTo(HaveOccurred())
				netDevices, err := clonedVM.GetNetworkDevices(context.TODO())
				Expect(err).NotTo(HaveOccurred())
				Expect(len(netDevices)).Should(Equal(2))
				// The interface type should be default vmxnet3
				dev1, ok := netDevices[0].(*vimTypes.VirtualVmxnet3)
				Expect(ok).Should(BeTrue())
				// TODO: enhance the test to verify the moref of the network matches the name of the network in spec.
				_, ok = dev1.Backing.(*vimTypes.VirtualEthernetCardDistributedVirtualPortBackingInfo)
				Expect(ok).Should(BeTrue())
				// The interface type should be e1000
				dev2, ok := netDevices[1].(*vimTypes.VirtualE1000)
				// TODO: enhance the test to verify the moref of the network matches the name of the network in spec.
				_, ok = dev2.Backing.(*vimTypes.VirtualEthernetCardDistributedVirtualPortBackingInfo)
				Expect(ok).Should(BeTrue())
			})
		})

		Context("from Content-library", func() {

			BeforeEach(func() {
				//set source to use CL
				vSphereConfig.ContentSource = integration.GetContentSourceID()
				err = session.ConfigureContent(context.TODO(), vSphereConfig.ContentSource)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should clone VM", func() {
				imageName := "test-item"

				vmClass := getVMClassInstance(testVMName, testNamespace)
				vmName := "CL_DeployedVM"
				vm := getVirtualMachineInstance(vmName, testNamespace, imageName, vmClass.Name)

				vmMetadata := map[string]string{}
				clonedVM, err := session.CloneVirtualMachine(context.TODO(), vm, *vmClass, nil, vmMetadata, "foo")
				Expect(err).NotTo(HaveOccurred())
				Expect(clonedVM.Name).Should(Equal(vmName))
			})
		})
	})

	Context("Session creation with invalid global extraConfig", func() {
		BeforeEach(func() {
			err = os.Setenv("JSON_EXTRA_CONFIG", "invalid-json")
			Expect(err).NotTo(HaveOccurred())
		})
		AfterEach(func() {
			err = os.Setenv("JSON_EXTRA_CONFIG", "")
			Expect(err).NotTo(HaveOccurred())
		})
		It("Should fail", func() {
			session, err = vsphere.NewSessionAndConfigure(context.TODO(), vSphereConfig, nil, nil)
			Expect(err.Error()).To(MatchRegexp("Unable to parse value of 'JSON_EXTRA_CONFIG' environment variable"))
		})
	})

	Describe("Clone VM with global metadata", func() {
		const (
			localKey  = "localK"
			localVal  = "localV"
			globalKey = "globalK"
			globalVal = "globalV"
		)

		JustBeforeEach(func() {
			//set source to use VM inventory

			session, err = vsphere.NewSessionAndConfigure(context.TODO(), vSphereConfig, nil, nil)
			vSphereConfig.ContentSource = ""
			err = session.ConfigureContent(context.TODO(), vSphereConfig.ContentSource)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("with vm metadata and global extraConfig", func() {
			BeforeEach(func() {
				err = os.Setenv("JSON_EXTRA_CONFIG", "{\""+globalKey+"\":\""+globalVal+"\"}")
				Expect(err).NotTo(HaveOccurred())
			})
			AfterEach(func() {
				err = os.Setenv("JSON_EXTRA_CONFIG", "")
				Expect(err).NotTo(HaveOccurred())
			})
			It("should copy all the values into the VM", func() {
			})

			Context("with global extraConfig", func() {
				It("should copy the values into the VM", func() {
					imageName := "DC0_H0_VM0"
					vmClass := getVMClassInstance(testVMName, testNamespace)
					vm := getVirtualMachineInstance(testVMName+"-extraConfig", testNamespace, imageName, vmClass.Name)
					vm.Spec.VmMetadata.Transport = "ExtraConfig"
					vmMetadata := map[string]string{localKey: localVal}

					clonedVM, err := session.CloneVirtualMachine(context.TODO(), vm, *vmClass, nil, vmMetadata, "foo")
					Expect(err).NotTo(HaveOccurred())
					Expect(clonedVM).ShouldNot(BeNil())

					keysFound := map[string]bool{localKey: false, globalKey: false}
					// Add all the default keys
					for k := range vsphere.DefaultExtraConfig {
						keysFound[k] = false
					}
					mo, err := clonedVM.ManagedObject(context.TODO())
					for _, option := range mo.Config.ExtraConfig {
						key := option.GetOptionValue().Key
						keysFound[key] = true
						if key == localKey {
							Expect(option.GetOptionValue().Value).Should(Equal(localVal))
						} else if key == globalKey {
							Expect(option.GetOptionValue().Value).Should(Equal(globalVal))
						} else if defaultVal, ok := vsphere.DefaultExtraConfig[key]; ok {
							Expect(option.GetOptionValue().Value).Should(Equal(defaultVal))
						}
					}
					for k, v := range keysFound {
						Expect(v).Should(BeTrue(), "Key %v not found in VM", k)
					}
				})
			})
			Context("without vm metadata or global extraConfig", func() {
				It("should copy the default values into the VM", func() {
					imageName := "DC0_H0_VM0"
					vmClass := getVMClassInstance(testVMName, testNamespace)
					vm := getVirtualMachineInstance(testVMName+"-default-extraConfig", testNamespace, imageName, vmClass.Name)
					clonedVM, err := session.CloneVirtualMachine(context.TODO(), vm, *vmClass, nil, nil, "foo")
					Expect(err).NotTo(HaveOccurred())
					Expect(clonedVM).ShouldNot(BeNil())

					keysFound := map[string]bool{}
					// Add all the default keys
					for k := range vsphere.DefaultExtraConfig {
						keysFound[k] = false
					}
					mo, err := clonedVM.ManagedObject(context.TODO())
					for _, option := range mo.Config.ExtraConfig {
						key := option.GetOptionValue().Key
						keysFound[key] = true
						if defaultVal, ok := vsphere.DefaultExtraConfig[key]; ok {
							Expect(option.GetOptionValue().Value).Should(Equal(defaultVal))
						}
					}
					for k, v := range keysFound {
						Expect(v).Should(BeTrue(), "Key %v not found in VM", k)
					}
				})
			})
		})

		Describe("Resource Pool", func() {
			var rpName string
			var rpSpec *v1alpha1.ResourcePoolSpec

			BeforeEach(func() {
				rpName = "test-folder"
				rpSpec = &vmoperatorv1alpha1.ResourcePoolSpec{
					Name: rpName,
				}
				rpMoId, err := session.CreateResourcePool(context.TODO(), rpSpec)
				Expect(err).NotTo(HaveOccurred())
				Expect(rpMoId).To(Not(BeEmpty()))
			})

			AfterEach(func() {
				// RP would already be deleted after the deletion test. But DeleteResourcePool handles delete of an RP if it's already deleted.
				Expect(session.DeleteResourcePool(context.TODO(), rpSpec.Name)).To(Succeed())
			})

			Context("Create a ResourcePool, verify it exists and delete it", func() {
				JustBeforeEach(func() {
				})

				It("create is tested in setup", func() {
					// Create is tested in JustBeforeEach
				})

				It("Verifies if a ResourcePool exists", func() {
					exists, err := session.DoesResourcePoolExist(context.TODO(), integration.DefaultNamespace, rpSpec.Name)
					Expect(exists).To(BeTrue())
					Expect(err).NotTo(HaveOccurred())
				})
				It("delete is tested in teardown", func() {
					// Delete is tested in JustAfterEach
				})
			})

			Context("Create two resource pools with the duplicate names", func() {
				It("second resource pool should fail to create", func() {
					// Try to create another ResourcePool with the same spec.
					rpMoId, err := session.CreateResourcePool(context.TODO(), rpSpec)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("ServerFaultCode: DuplicateName"))
					Expect(rpMoId).To(BeEmpty())
				})
			})
			Context("Delete a Resource Pool that doesn't exist", func() {
				It("should succeed", func() {
					Expect(session.DeleteResourcePool(context.TODO(), "nonexistent-resourcepool")).To(Succeed())
				})
			})
		})

		Describe("Folder", func() {
			var folderName string
			var folderSpec *v1alpha1.FolderSpec

			BeforeEach(func() {
				folderName = "test-folder"
				folderSpec = &vmoperatorv1alpha1.FolderSpec{
					Name: folderName,
				}
			})

			Context("Create a Folder, verify it exists and delete it", func() {
				JustBeforeEach(func() {
					folderMoId, err := session.CreateFolder(context.TODO(), folderSpec)
					Expect(err).NotTo(HaveOccurred())
					Expect(folderMoId).To(Not(BeEmpty()))

				})
				JustAfterEach(func() {
					Expect(session.DeleteFolder(context.TODO(), folderName)).To(Succeed())
				})

				It("create is tested in setup", func() {
					// Create is tested in JustBeforeEach
				})

				It("Verifies if a Folder exists", func() {
					exists, err := session.DoesFolderExist(context.TODO(), integration.DefaultNamespace, folderName)
					Expect(exists).To(BeTrue())
					Expect(err).NotTo(HaveOccurred())
				})
				It("delete is tested in teardown", func() {
					// Delete is tested in JustAfterEach
				})
			})

			Context("Create two folders with the duplicate names", func() {
				It("Second folder should fail to create", func() {
					folderMoId1, err := session.CreateFolder(context.TODO(), folderSpec)
					Expect(err).NotTo(HaveOccurred())
					Expect(folderMoId1).To(Not(BeEmpty()))

					// Try to crete another folder with the same spec.
					folderMoId2, err := session.CreateFolder(context.TODO(), folderSpec)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("ServerFaultCode: DuplicateName"))
					Expect(folderMoId2).To(BeEmpty())
				})
			})
			Context("Delete a Folder that doesnt exist", func() {
				It("should succeed", func() {
					Expect(session.DeleteFolder(context.TODO(), folderSpec.Name)).To(Succeed())
				})
			})
		})

		Describe("Clone VM gracefully fails", func() {
			Context("Should fail gracefully", func() {
				var savedDatastoreAttribute string
				vm := &vmoperatorv1alpha1.VirtualMachine{
					ObjectMeta: metav1.ObjectMeta{
						Name: "TestVM",
					},
				}

				BeforeEach(func() {
					savedDatastoreAttribute = vSphereConfig.Datastore
				})

				AfterEach(func() {
					vSphereConfig.Datastore = savedDatastoreAttribute
					vSphereConfig.ContentSource = ""
					vSphereConfig.StorageClassRequired = false
				})

				It("with existing content source, empty datastore and empty profile id", func() {
					vSphereConfig.Datastore = ""
					vSphereConfig.ContentSource = integration.GetContentSourceID()
					session, err := vsphere.NewSessionAndConfigure(context.TODO(), vSphereConfig, nil, nil)
					Expect(err).NotTo(HaveOccurred())
					clonedVM, err :=
						session.CloneVirtualMachine(context.TODO(), vm, v1alpha1.VirtualMachineClass{}, nil, nil, "")
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError("cannot clone VM when neither storage class or datastore is specified"))
					Expect(clonedVM).Should(BeNil())
				})
				It("with existing content source but mandatory profile id is not set", func() {
					vSphereConfig.ContentSource = integration.GetContentSourceID()
					vSphereConfig.StorageClassRequired = true
					session, err = vsphere.NewSessionAndConfigure(context.TODO(), vSphereConfig, nil, nil)
					Expect(err).NotTo(HaveOccurred())
					clonedVM, err :=
						session.CloneVirtualMachine(context.TODO(), vm, v1alpha1.VirtualMachineClass{}, nil, nil, "")
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError("storage class is required but not specified"))
					Expect(clonedVM).Should(BeNil())
				})
				It("without content source and missing mandatory profile ID", func() {
					vSphereConfig.StorageClassRequired = true
					session, err = vsphere.NewSessionAndConfigure(context.TODO(), vSphereConfig, nil, nil)
					Expect(err).NotTo(HaveOccurred())
					clonedVM, err :=
						session.CloneVirtualMachine(context.TODO(), vm, v1alpha1.VirtualMachineClass{}, nil, nil, "")
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError("storage class is required but not specified"))
					Expect(clonedVM).Should(BeNil())
				})
			})
		})
	})
})