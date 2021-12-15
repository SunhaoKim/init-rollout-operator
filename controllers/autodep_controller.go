/*
Copyright 2021.

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
	"context"

	appsv1alpha1 "init_rollout_operator/api/v1alpha1"
	"init_rollout_operator/resources"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// AutodepReconciler reconciles a Autodep object
type AutodepReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

const (
	deploymentfinalizer = "initrolloutoperator"
)

//+kubebuilder:rbac:groups=apps.autodep.com,resources=autodeps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps.autodep.com,resources=autodeps/status,verbs=create;get;update;patch
//+kubebuilder:rbac:groups=apps.autodep.com,resources=autodeps/finalizers,verbs=delete;update
// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Autodep object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *AutodepReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	autodep := &appsv1alpha1.Autodep{}
	_ = r.Log.WithValues("this operator is auto deploy deployment and service", req.NamespacedName)
	//检查autodep对象
	err := r.Get(ctx, req.NamespacedName, autodep)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		//忽略掉 not-found 错误，它们不能通过重新排队修复（要等待新的通知）
		//在删除一个不存在的对象时，可能会报这个错误。
		r.Log.Error(err, "failed get autodep")
		return ctrl.Result{}, err
	}
	dep := resources.DeploymentForbackend(autodep)
	//预删除逻辑
	if autodep.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(autodep, deploymentfinalizer) {
			controllerutil.AddFinalizer(autodep, deploymentfinalizer)
			if err := r.Update(ctx, autodep); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		//将删除deployment
		if controllerutil.ContainsFinalizer(autodep, deploymentfinalizer) {
			if err := r.Delete(ctx, dep); err != nil {
				return ctrl.Result{}, err
			}
			//对预删除逻辑进行更新
			controllerutil.RemoveFinalizer(autodep, deploymentfinalizer)
			if err := r.Update(ctx, autodep); err != nil {
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
	//检查deployment,不存在及创建
	err = r.Get(ctx, types.NamespacedName{Namespace: autodep.Namespace, Name: dep.Name}, dep)
	if err != nil && errors.IsNotFound(err) {
		r.Log.Info("deployment not found just create", dep.Namespace, dep.Name)
		if err := r.Create(ctx, dep); err != nil {
			r.Log.Error(err, "failed create deployment")
			return ctrl.Result{}, err
		}
		r.Log.V(1).Info("create deployment success", dep.Namespace, dep.Name)
	}
	// your logic here
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AutodepReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.Autodep{}).
		Complete(r)
}
