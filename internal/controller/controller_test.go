package controller

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	fancysecretsv1 "github.com/anthonymilazzo/secretsauce/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("RandomSecret controller", func() {

	const (
		RandomSecretName      = "test-random-secret"
		RandomSecretNamespace = "default"
		RandomSecretLength    = 10
		SecretName            = "test-secret"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	var (
		randomSecret        *fancysecretsv1.RandomSecret
		ctx                 = context.Background()
		randomSecretKey     = types.NamespacedName{Name: RandomSecretName, Namespace: RandomSecretNamespace}
		secretKey           = types.NamespacedName{Name: SecretName, Namespace: RandomSecretNamespace}
		createdRandomSecret *fancysecretsv1.RandomSecret
		createdSecret       *corev1.Secret
	)

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
				Length:     int32(RandomSecretLength),
				SecretName: SecretName,
			},
		}
		createdRandomSecret = &fancysecretsv1.RandomSecret{}
		createdSecret = &corev1.Secret{}
	})

	Context("RandomSecret Tests", func() {
		It("should create the RandomSecret", func() {
			Expect(k8sClient.Create(ctx, randomSecret)).Should(Succeed())
		})

		It("should create the Secret", func() {
			_ = k8sClient.Create(ctx, randomSecret)

			Expect(k8sClient.Get(ctx, randomSecretKey, createdRandomSecret)).Should(Succeed())

			Eventually(func() error {
				return k8sClient.Get(ctx, secretKey, createdSecret)
			}, timeout, interval).Should(Succeed())
		})

		It("should have the correct length", func() {
			_ = k8sClient.Create(ctx, randomSecret)

			Expect(k8sClient.Get(ctx, randomSecretKey, createdRandomSecret)).Should(Succeed())

			Eventually(func() error {
				return k8sClient.Get(ctx, secretKey, createdSecret)
			}, timeout, interval).Should(Succeed())

			Expect(createdSecret.Data["value"]).Should(HaveLen(RandomSecretLength))
		})

		It("should update the Secret", func() {
			_ = k8sClient.Create(ctx, randomSecret)

			Expect(k8sClient.Get(ctx, randomSecretKey, createdRandomSecret)).Should(Succeed())

			Eventually(func() error {
				return k8sClient.Get(ctx, secretKey, createdSecret)
			}, timeout, interval).Should(Succeed())

			// Save the value of the secret
			value := createdSecret.Data["value"]

			// Update the RandomSecret
			Expect(k8sClient.Get(ctx, randomSecretKey, createdRandomSecret)).Should(Succeed())
			createdRandomSecret.Spec.Length = 20
			Expect(k8sClient.Update(ctx, createdRandomSecret)).Should(Succeed())

			// Wait for the Secret to be updated
			Eventually(func() []byte {
				_ = k8sClient.Get(ctx, secretKey, createdSecret)
				return createdSecret.Data["value"]
			}, timeout, interval).ShouldNot(Equal(value))

			// Verify the length has changed
			Expect(len(string(createdSecret.Data["value"]))).Should(Equal(20))
		})

		It("should delete the Secret", func() {
			_ = k8sClient.Create(ctx, randomSecret)

			Expect(k8sClient.Get(ctx, randomSecretKey, createdRandomSecret)).Should(Succeed())

			Eventually(func() error {
				return k8sClient.Get(ctx, secretKey, createdSecret)
			}, timeout, interval).Should(Succeed())

			// Delete the RandomSecret
			Expect(k8sClient.Delete(ctx, createdRandomSecret)).Should(Succeed())

			// Wait for the Secret to be deleted
			Eventually(func() error {
				return k8sClient.Get(ctx, secretKey, createdSecret)
			}, timeout, interval).ShouldNot(Succeed())
		})
	})
})
