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
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	utils "github.com/rikatz/falco-network-operator/utils"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	macroOutboundRule = "- macro: outbound\n" +
		"  condition: syscall.type=connect and evt.dir=< and (fd.typechar=4 or fd.typechar=6)\n"
)

// NetworkPolicyReconciler reconciles a NetworkPolicy object
type NetworkPolicyReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=networking.kubernetes.io,resources=networkpolicies,verbs=get;list;watch
//+kubebuilder:rbac:groups=networking.kubernetes.io,resources=networkpolicies/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the NetworkPolicy object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.2/pkg/reconcile
func (r *NetworkPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("networkpolicy", req.NamespacedName)

	var networkPolicies networkingv1.NetworkPolicyList

	if err := r.List(ctx, &networkPolicies, client.InNamespace(req.Namespace)); err != nil {
		log.Error(err, "unable to fetch NetworkPolicy")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	ns := req.Namespace
	netpolName := req.Name

	log.V(1).Info("falco rule update", "namespace", ns, "rule", netpolName)

	var falcoRule strings.Builder

	var egressNetPols []networkingv1.NetworkPolicy

	if len(networkPolicies.Items) > 0 {
		for _, np := range networkPolicies.Items {
			if len(np.Spec.Egress) > 0 {
				egressNetPols = append(egressNetPols, np)
			}
		}
	}

	if len(egressNetPols) > 0 {
		falcoRule.WriteString(macroOutboundRule)
		// TODO: Add some better output here, like the pod name, podnamespace, etc
		namespaceRule := "- rule: namespace unexpected outbound connection\n" +
			"  desc: A pod in namespace attempted to connect to a non configured Network Policy\n" +
			"  output: \"Outbound network traffic from Pod on unexpected IP/Port\"\n" +
			"  priority: WARNING\n"

		falcoRule.WriteString(namespaceRule)
		namespaceCondition := fmt.Sprintf("  condition: outbound and k8s.ns.name = \"%s\" and not (%s_outbound_rules)\n", ns, ns)
		falcoRule.WriteString(namespaceCondition)

		ruleMacroHeader := fmt.Sprintf("- macro: %s_outbound_rules\n", ns)

		for k, np := range egressNetPols {
			falcoRule.WriteString(ruleMacroHeader)
			var orString string
			if k > 0 {
				falcoRule.WriteString("  append: true\n")
				orString = "or "
			}

			condition := fmt.Sprintf("  condition: %s(%s)\n", orString, utils.NetPol2FalcoCond(np))
			falcoRule.WriteString(condition)
		}
	}
	fmt.Print(falcoRule.String())

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NetworkPolicyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkingv1.NetworkPolicy{}).
		Complete(r)
}
