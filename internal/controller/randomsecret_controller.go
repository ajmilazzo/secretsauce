/*
Copyright 2023 Anthony Milazzo.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	fancysecretsv1 "github.com/anthonymilazzo/secretsauce/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RandomSecretReconciler reconciles a RandomSecret object
type RandomSecretReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=fancysecrets.secretsauce.anthonymilazzo.com,resources=randomsecrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=fancysecrets.secretsauce.anthonymilazzo.com,resources=randomsecrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=fancysecrets.secretsauce.anthonymilazzo.com,resources=randomsecrets/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;patch;delete

// Reconcile moves the current state of the cluster closer to the desired state.
// It compares the state specified by the RandomSecret object against the actual cluster state.
func (r *RandomSecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	log.Info("### Begin Reconciling RandomSecret ###")

	var randomSecret fancysecretsv1.RandomSecret
	if err := r.Get(ctx, req.NamespacedName, &randomSecret); err != nil {
		// Handle not found error or requeue upon other errors
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.ensureSecretExists(ctx, &randomSecret); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager initializes the controller with the Manager.
func (r *RandomSecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&fancysecretsv1.RandomSecret{}).
		Complete(r)
}

// ensureSecretExists either creates or updates the secret.
func (r *RandomSecretReconciler) ensureSecretExists(ctx context.Context, randomSecret *fancysecretsv1.RandomSecret) error {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      randomSecret.Spec.SecretName,
			Namespace: randomSecret.Namespace,
		},
	}

	_, err := ctrl.CreateOrUpdate(ctx, r.Client, secret, func() error {
		if secret.Data == nil {
			secret.Data = make(map[string][]byte)
		}

		value, exists := secret.Data["value"]
		// If the length changes, we need to update the secret
		correctLength := int32(len(string(value))) == randomSecret.Spec.Length
		if !exists || !correctLength {
			return setRandomData(secret, randomSecret.Spec.Length)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return r.Status().Update(ctx, randomSecret)
}
