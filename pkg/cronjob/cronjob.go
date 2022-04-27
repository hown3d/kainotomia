package cronjob

import (
	"context"
	"fmt"

	config "github.com/hown3d/kainotomia/pkg/jobconfig"
	"golang.org/x/oauth2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	batchClientV1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	coreClientV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Client struct {
	cronjobClient batchClientV1.CronJobInterface
	secretsClient coreClientV1.SecretInterface
	args          Options
}

type Options struct {
	JobImage string
}

type Option func(*Options)

func WithJobImage(jobImage string) Option {
	return func(args *Options) {
		args.JobImage = jobImage
	}
}

func New(kubeclient *kubernetes.Clientset, namespace string, opts ...Option) Client {
	args := &Options{
		JobImage: "quay.io/hown3d/kainotomia",
	}
	for _, opt := range opts {
		opt(args)
	}
	return Client{
		cronjobClient: kubeclient.BatchV1().CronJobs(namespace),
		secretsClient: kubeclient.CoreV1().Secrets(namespace),
	}
}

func (c Client) Create(ctx context.Context, name string, token oauth2.Token) error {
	cj, err := c.cronjobClient.Create(ctx, generateCronJob(jobSpec{
		jobImage:   c.args.JobImage,
		playlistID: name,
		schedule:   everyFriday(),
		tokenSecretRef: corev1.LocalObjectReference{
			Name: tokenSecretName(name),
		},
		cmRef: corev1.LocalObjectReference{
			Name: configmapName(name),
		},
	}), v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating cronjob %v: %w", name, err)
	}
	err = config.StoreToken(ctx, c.secretsClient, name, token, cj.GetOwnerReferences())
	if err != nil {
		return fmt.Errorf("storing to")
	}
	return nil
}

func tokenSecretName(name string) string {
	return fmt.Sprintf("%v-token", name)
}

func configmapName(name string) string {
	return fmt.Sprintf("%v-job-config", name)
}

type jobSpec struct {
	jobImage       string
	playlistID     string
	schedule       string
	cmRef          corev1.LocalObjectReference
	tokenSecretRef corev1.LocalObjectReference
}

func generateCronJob(s jobSpec) *batchv1.CronJob {
	return &batchv1.CronJob{
		ObjectMeta: v1.ObjectMeta{
			Name: s.playlistID,
		},
		Spec: batchv1.CronJobSpec{
			Schedule: s.schedule,
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "job",
									Image: s.jobImage,
									EnvFrom: []corev1.EnvFromSource{
										{
											SecretRef: &corev1.SecretEnvSource{
												LocalObjectReference: s.tokenSecretRef,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func everyFriday() string {
	return fmt.Sprintf("59 23 * * 5")
}
