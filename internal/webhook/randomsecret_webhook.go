package webhook

import (
	"context"
	"fmt"

	fancysecretsv1 "github.com/anthonymilazzo/secretsauce/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var logger = log.Log.WithName("randomsecret-resource")

// RandomSecretValidator validates RandomSecrets
type RandomSecretValidator struct {
	Client client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:webhook:path=/validate-fancysecrets-secretsauce-anthonymilazzo-com-v1-randomsecret,mutating=false,failurePolicy=fail,sideEffects=None,groups=fancysecrets.secretsauce.anthonymilazzo.com,resources=randomsecrets,verbs=create;update,versions=v1,name=vrandomsecret.kb.io,admissionReviewVersions=v1

func (r *RandomSecretValidator) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&fancysecretsv1.RandomSecret{}).
		WithValidator(r).
		Complete()
}

// ValidateCreate checks the creation of a RandomSecret object
func (r *RandomSecretValidator) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	logger.Info("validate create")
	return r.validateRandomSecretObject(ctx, obj)
}

// ValidateUpdate checks the update of a RandomSecret object
func (r *RandomSecretValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	logger.Info("validate update")
	return r.validateRandomSecretObject(ctx, newObj)
}

// ValidateDelete is a placeholder for when a RandomSecret object is deleted
func (r *RandomSecretValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	logger.Info("validate delete")
	return nil, nil
}

func (r *RandomSecretValidator) validateRandomSecretObject(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	log := log.FromContext(ctx)
	log.Info("### Begin Validation ###")

	randomSecret, ok := obj.(*fancysecretsv1.RandomSecret)
	if !ok {
		return nil, fmt.Errorf("expected object of type RandomSecret but got %T", obj)
	}

	return nil, r.validateRandomSecret(ctx, randomSecret)
}

func (r *RandomSecretValidator) validateRandomSecret(ctx context.Context, randomSecret *fancysecretsv1.RandomSecret) error {
	secretPolicies := &fancysecretsv1.SecretPolicyList{}
	if err := r.Client.List(ctx, secretPolicies); err != nil {
		return err
	}

	for _, secretPolicy := range secretPolicies.Items {
		if randomSecret.Spec.Length < secretPolicy.Spec.MinLength {
			return fmt.Errorf("RandomSecret length %d is less than SecretPolicy MinLength %d", randomSecret.Spec.Length, secretPolicy.Spec.MinLength)
		}
	}

	return nil
}
