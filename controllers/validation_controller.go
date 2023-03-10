/*
Copyright 2023.

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

package controllers

import (

	//yaml "gopkg.in/yaml.v3"
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	//"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	//"github.com/openstack-k8s-operators/lib-common/modules/storage"

	validationv1alpha1 "github.com/matbu/validation-operator/api/v1alpha1"
)

// ValidationReconciler reconciles a Validation object
type ValidationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=validation.redhat.com,resources=validations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=validation.redhat.com,resources=validations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=validation.redhat.com,resources=validations/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Validation object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *ValidationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	//_ = log.FromContext(ctx)

	// Fetch the Validation instance
	validation := &validationv1alpha1.Validation{}
	err := r.Get(ctx, req.NamespacedName, validation)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Check if the dep already exists, if not create a new one
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: validation.Name, Namespace: validation.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new dep
		dep := r.deploymentForValidation(validation)

		err = r.Create(ctx, dep)
		if err != nil {
			//fmt.Println(err.Error())
			return ctrl.Result{}, err
		}

		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		//fmt.Println(err.Error())

		return ctrl.Result{}, err
	}

	// Ensure the deployment size is the same as the spec
	size := validation.Spec.Size
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err = r.Update(ctx, found)
		if err != nil {
			//fmt.Println(err.Error())
			//log.Error(err, "Failed to update Deployment", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)
			return ctrl.Result{}, err
		}
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}

	// Update the Validation status with the pod names
	// List the pods for this validation's deployment
	podList := &corev1.PodList{}
	listOpts := []client.ListOption{
		client.InNamespace(validation.Namespace),
		client.MatchingLabels(labelsForValidation(validation.Name)),
	}
	if err = r.List(ctx, podList, listOpts...); err != nil {
		//fmt.Println(err.Error())
		//log.Error(err, "Failed to list pods", "validation.Namespace", validation.Namespace, "validation.Name", validation.Name)
		return ctrl.Result{}, err
	}
	podNames := getPodNames(podList.Items)

	// Update status.Nodes if needed
	if !reflect.DeepEqual(podNames, validation.Status.Nodes) {
		validation.Status.Nodes = podNames
		err := r.Status().Update(ctx, validation)
		if err != nil {
			//fmt.Println(err.Error())
			//log.Error(err, "Failed to update validation status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil

}

// func (r *ValidationReconciler) getValidationInstance(ctx context.Context, req ctrl.Request) (*validationv1alpha1.Validation, error) {
// 	// Fetch the Validation instance
// 	instance := &validationv1alpha1.Validation{}
// 	err := r.Get(ctx, req.NamespacedName, instance)
// 	if err != nil {
// 		if errors.IsNotFound(err) {
// 			// Request object not found, could have been deleted after reconcile request.
// 			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
// 			// Return and don't requeue
// 			fmt.Println("Validation resource not found. Ignoring since object must be deleted")
// 			//log.Info("Validation resource not found. Ignoring since object must be deleted")
// 			return &validationv1alpha1.Validation{}, nil
// 		}
// 		// Error reading the object - requeue the request.
// 		fmt.Println(err.Error())
// 		//log.Error(err, "Failed to get Validation")
// 		return &validationv1alpha1.Validation{}, err
// 	}
//
// 	return instance, nil
// }

// jobForValidation returns a Validation Job object
func (r *ValidationReconciler) deploymentForValidation(instance *validationv1alpha1.Validation) *appsv1.Deployment {
	ls := labelsForValidation(instance.Name)

	replicas := instance.Spec.Size
	args := instance.Spec.Args
	command := instance.Spec.Command
	// if len(args) == 0 {
	// 	if len(instance.Spec.Validation) == 0 {
	// 		instance.Spec.Validation = "validation.yaml"
	// 	}
	// 	args = []string{"validation", "run", instance.Spec.Validation}
	// }

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      instance.Name,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicy(instance.Spec.RestartPolicy),
					Containers: []corev1.Container{{
						ImagePullPolicy: "Always",
						Image:           instance.Spec.Image,
						Name:            instance.Spec.Name,
						Command:         command,
						Args:            args,
					}},
				},
			},
		},
	}

	// Set Validation instance as the owner and controller
	ctrl.SetControllerReference(instance, dep, r.Scheme)
	return dep
}

// labelsForValidation returns the labels for selecting the resources
// belonging to the given valdation CR name.
func labelsForValidation(name string) map[string]string {
	return map[string]string{"app": "validation", "validation_cr": name}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(pods []corev1.Pod) []string {
	var podNames []string
	for _, pod := range pods {
		podNames = append(podNames, pod.Name)
	}
	return podNames
}

// SetupWithManager sets up the controller with the Manager.
func (r *ValidationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&validationv1alpha1.Validation{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
