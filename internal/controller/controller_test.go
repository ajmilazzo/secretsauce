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

	AfterEach(func() {
		_ = k8sClient.Get(ctx, randomSecretKey, createdRandomSecret)
		_ = k8sClient.Delete(ctx, createdRandomSecret)

		_ = k8sClient.Get(ctx, secretKey, createdSecret)
		_ = k8sClient.Delete(ctx, createdSecret)
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
	})
})
