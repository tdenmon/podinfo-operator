/*
Copyright 2024.

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

package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appv1 "tdenmon/angi-takehome/api/v1"
)

var _ = Describe("PodInfo Controller", func() {
	const (
		PodInfoName      = "test-podinfo"
		PodInfoNamespace = "default"
	)

	Context("When reconciling a resource", func() {
		ctx := context.Background()

		typeNamespacedName := types.NamespacedName{
			Name:      PodInfoName,
			Namespace: PodInfoNamespace,
		}
		podinfo := &appv1.PodInfo{}

		BeforeEach(func() {
			By("creating the custom resource for the Kind PodInfo")
			err := k8sClient.Get(ctx, typeNamespacedName, podinfo)
			if err != nil && errors.IsNotFound(err) {
				resource := &appv1.PodInfo{
					ObjectMeta: metav1.ObjectMeta{
						Name:      PodInfoName,
						Namespace: PodInfoNamespace,
					},
					Spec: appv1.PodInfoSpec{
						ReplicaCount: 3,
						ResourceInfo: appv1.Resources{
							CpuRequest:  resource.MustParse("100m"),
							MemoryLimit: resource.MustParse("64Mi"),
						},
						ImageInfo: appv1.Image{
							ImageRepository: "ghcr.io/stefanprodan/podinfo",
							ImageTag:        "latest",
						},
						UiInfo: appv1.Ui{
							UiColor:   "#336699",
							UiMessage: "Test string",
						},
						RedisInfo: appv1.Redis{
							RedisEnabled: true,
						},
					},
				}
				Expect(k8sClient.Create(ctx, resource)).To(Succeed())
			}
		})

		AfterEach(func() {
			resource := &appv1.PodInfo{}
			err := k8sClient.Get(ctx, typeNamespacedName, resource)
			Expect(err).NotTo(HaveOccurred())

			By("Cleanup the specific resource instance PodInfo")
			Expect(k8sClient.Delete(ctx, resource)).To(Succeed())
		})
		It("should successfully reconcile the resource", func() {
			By("Reconciling the created resource")
			controllerReconciler := &PodInfoReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
		})
		It("should have the correct replica count", func() {
			Eventually(func() bool {
				err := k8sClient.Get(ctx, typeNamespacedName, podinfo)
				return err == nil
			}).Should(BeTrue())
			Expect(podinfo.Spec.ReplicaCount).Should(Equal(int32(3)))
		})
		It("should have the correct resource info", func() {
			Eventually(func() bool {
				err := k8sClient.Get(ctx, typeNamespacedName, podinfo)
				return err == nil
			}).Should(BeTrue())
			expectedResourceInfo := appv1.Resources{
				CpuRequest:  resource.MustParse("100m"),
				MemoryLimit: resource.MustParse("64Mi"),
			}
			Expect(podinfo.Spec.ResourceInfo).Should(Equal(expectedResourceInfo))
		})
		It("should have the correct image info", func() {
			Eventually(func() bool {
				err := k8sClient.Get(ctx, typeNamespacedName, podinfo)
				return err == nil
			}).Should(BeTrue())
			expectedImageInfo := appv1.Image{
				ImageRepository: "ghcr.io/stefanprodan/podinfo",
				ImageTag:        "latest",
			}
			Expect(podinfo.Spec.ImageInfo).Should(Equal(expectedImageInfo))
		})
		It("should have the correct ui info", func() {
			Eventually(func() bool {
				err := k8sClient.Get(ctx, typeNamespacedName, podinfo)
				return err == nil
			}).Should(BeTrue())
			expectedUiInfo := appv1.Ui{
				UiColor:   "#336699",
				UiMessage: "Test string",
			}
			Expect(podinfo.Spec.UiInfo).Should(Equal(expectedUiInfo))
		})
		It("should have the correct redis info", func() {
			Eventually(func() bool {
				err := k8sClient.Get(ctx, typeNamespacedName, podinfo)
				return err == nil
			}).Should(BeTrue())
			expectedRedisInfo := appv1.Redis{
				RedisEnabled: true,
			}
			Expect(podinfo.Spec.RedisInfo).Should(Equal(expectedRedisInfo))
		})
		It("should update podinfo deployment on replica count change", func() {
			Eventually(func() bool {
				err := k8sClient.Get(ctx, typeNamespacedName, podinfo)
				return err == nil
			}).Should(BeTrue())

			controllerReconciler := &PodInfoReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			podinfo.Spec.ReplicaCount = 2
			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(k8sClient.Update(ctx, podinfo)).Should(HaveOccurred())
			Expect(podinfo.Spec.ReplicaCount).Should(Equal(int32(2)))
		})
		It("should update podinfo deployment on resource changes", func() {
			Eventually(func() bool {
				err := k8sClient.Get(ctx, typeNamespacedName, podinfo)
				return err == nil
			}).Should(BeTrue())

			controllerReconciler := &PodInfoReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			podinfo.Spec.ResourceInfo = appv1.Resources{
				CpuRequest:  resource.MustParse("50m"),
				MemoryLimit: resource.MustParse("32Mi"),
			}
			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(k8sClient.Update(ctx, podinfo)).Should(HaveOccurred())
			Expect(podinfo.Spec.ResourceInfo.MemoryLimit).Should(Equal(resource.MustParse("32Mi")))
			Expect(podinfo.Spec.ResourceInfo.CpuRequest).Should(Equal(resource.MustParse("50m")))
		})
		It("should update podinfo deployment on ui changes", func() {
			Eventually(func() bool {
				err := k8sClient.Get(ctx, typeNamespacedName, podinfo)
				return err == nil
			}).Should(BeTrue())

			controllerReconciler := &PodInfoReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}

			podinfo.Spec.UiInfo = appv1.Ui{
				UiColor:   "#432123",
				UiMessage: "Another string",
			}
			_, err := controllerReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespacedName,
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(k8sClient.Update(ctx, podinfo)).Should(HaveOccurred())
			Expect(podinfo.Spec.UiInfo.UiColor).Should(Equal("#432123"))
			Expect(podinfo.Spec.UiInfo.UiMessage).Should(Equal("Another string"))
		})
	})
})
