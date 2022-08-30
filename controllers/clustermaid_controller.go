/*
Copyright 2022.

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
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	v1 "cdx.foc/clusterMaid/api/v1"

	feat "cdx.foc/clusterMaid/pkg/feature"
)

type Clock interface {
	Now() time.Time
}

// ClusterMaidReconciler reconciles a ClusterMaid object
type ClusterMaidReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=webapp.cdx.foc,resources=clustermaids,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=webapp.cdx.foc,resources=clustermaids/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=webapp.cdx.foc,resources=clustermaids/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ClusterMaid object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *ClusterMaidReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx)
	var scanningActivities v1.ClusterMaid

	//log.Info("Running app: %s", time.Now().String())

	if err := r.Get(ctx, req.NamespacedName, &scanningActivities); err != nil {
		log.Error(err, "unable to fetch clustermaid crd.")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	cm := feat.ClusterManager{}

	cm.Execute(scanningActivities.Spec.Namespace)

	var duration time.Duration = 10000000000 // 10 seconds
	return ctrl.Result{
		RequeueAfter: duration,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ClusterMaidReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1.ClusterMaid{}).
		Complete(r)
}
