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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appv1 "tdenmon/angi-takehome/api/v1"
)

// PodInfoReconciler reconciles a PodInfo object
type PodInfoReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=app.podinfo.angi.takehome,resources=podinfoes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.podinfo.angi.takehome,resources=podinfoes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.podinfo.angi.takehome,resources=podinfoes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodInfo object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *PodInfoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	var podinfo appv1.PodInfo
	errGet := r.Get(ctx, req.NamespacedName, &podinfo)
	if errGet != nil {
		log.Error(errGet, "Error getting podinfo")
		return ctrl.Result{}, client.IgnoreNotFound(errGet)
	}

	if podinfo.Spec.RedisInfo.RedisEnabled {
		redisDeploy := NewRedisDeploy(&podinfo)
		_, errCreate := ctrl.CreateOrUpdate(ctx, r.Client, redisDeploy, func() error {
			return ctrl.SetControllerReference(&podinfo, redisDeploy, r.Scheme)
		})

		if errCreate != nil {
			log.Error(errCreate, "Error creating redis deployment")
			return ctrl.Result{}, nil
		}

		redisService := NewRedisService(&podinfo)
		_, errCreate = ctrl.CreateOrUpdate(ctx, r.Client, redisService, func() error {
			return ctrl.SetControllerReference(&podinfo, redisService, r.Scheme)
		})

		if errCreate != nil {
			log.Error(errCreate, "Error creating redis service")
			return ctrl.Result{}, nil
		}
	}

	err := r.Status().Update(ctx, &podinfo)
	if err != nil {
		return ctrl.Result{}, err
	}
	// Make Deployment(?) for Redis if enabled
	// Make Service for Redis, if enabled
	// Make Deployment with PodInfo pods, pointed at Redis service via env var
	// Make Service for PodInfo Deployment
	// CreateOrUpdate all of the above
	// Check for errors
	// Update podinfo object
	// ???
	// Profit!

	return ctrl.Result{}, nil
}

func NewRedisDeploy(podinfo *appv1.PodInfo) *appsv1.Deployment {
	labels := map[string]string{
		"app": podinfo.Name,
	}

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podinfo.Name,
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "redis",
							Image: "redis:latest",
							Ports: []corev1.ContainerPort{
								{
									Name:          "redis",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 6379,
								},
							},
						},
					},
				},
			},
		},
	}
}

func NewRedisService(podinfo *appv1.PodInfo) *corev1.Service {
	labels := map[string]string{
		"app": podinfo.Name,
	}

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podinfo.Name,
			Labels:    labels,
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "redis",
					Protocol:   corev1.ProtocolTCP,
					Port:       6379,
					TargetPort: intstr.FromInt(6379),
				},
			},
			Selector:  labels,
			ClusterIP: "",
		},
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodInfoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1.PodInfo{}).
		Complete(r)
}
