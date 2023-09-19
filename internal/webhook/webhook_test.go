package webhook

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fancysecretsv1 "github.com/anthonymilazzo/secretsauce/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("RandomSecret webhook", func() {
	const (
		RandomSecretName      = "test-random-secret"
		RandomSecretNamespace = "default"
		SecretName            = "test-secret"
		SecretPolicyName      = "test-secret-policy"
		SecretPolicyNamespace = "default"
		MinLength             = 10
	)

	var (
		ctx          = context.Background()
		randomSecret *fancysecretsv1.RandomSecret
	)

	BeforeEach(func() {
		secretPolicy := &fancysecretsv1.SecretPolicy{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "fancysecrets.secretsauce.anthonymilazzo.com/v1",
				Kind:       "SecretPolicy",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      SecretPolicyName,
				Namespace: SecretPolicyNamespace,
			},
			Spec: fancysecretsv1.SecretPolicySpec{
				MinLength: MinLength,
			},
		}
		Expect(k8sClient.Create(ctx, secretPolicy)).Should(Succeed())
	})

	// Clean up the SecretPolicy at the end
	AfterEach(func() {
		secretPolicyKey := types.NamespacedName{Name: SecretPolicyName, Namespace: SecretPolicyNamespace}
		createdSecretPolicy := &fancysecretsv1.SecretPolicy{}
		_ = k8sClient.Get(ctx, secretPolicyKey, createdSecretPolicy)
		_ = k8sClient.Delete(ctx, createdSecretPolicy)
	})

	Context("Secret creation", func() {
		BeforeEach(func() {
			randomSecret = &fancysecretsv1.RandomSecret{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "fancysecrets.secretsauce.anthonymilazzo.com/v1",
					Kind:       "RandomSecret",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      RandomSecretName,
					Namespace: RandomSecretNamespace,
				},
				Spec: fancysecretsv1.RandomSecretSpec{
					SecretName: SecretName,
				},
			}
		})

		// Clean up the RandomSecret at the end
		AfterEach(func() {
			randomSecretKey := types.NamespacedName{Name: RandomSecretName, Namespace: RandomSecretNamespace}
			createdRandomSecret := &fancysecretsv1.RandomSecret{}
			_ = k8sClient.Get(ctx, randomSecretKey, createdRandomSecret)
			_ = k8sClient.Delete(ctx, createdRandomSecret)
		})

		When("creating RandomSecret with length < minLength", func() {
			BeforeEach(func() {
				randomSecret.Spec.Length = 3
			})
			It("Should fail to create the secret", func() {
				Expect(k8sClient.Create(ctx, randomSecret)).ShouldNot(Succeed())
			})
		})

		When("creating RandomSecret with length > minLength", func() {
			BeforeEach(func() {
				randomSecret.Spec.Length = 20
			})
			It("Should successfully create the secret", func() {
				Expect(k8sClient.Create(ctx, randomSecret)).Should(Succeed())
			})
		})
	})
})
