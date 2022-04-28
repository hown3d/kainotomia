package job

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	batchClientV1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	coreClientV1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type Client struct {
	cronjobClient   batchClientV1.CronJobInterface
	secretClient    coreClientV1.SecretInterface
	configMapClient coreClientV1.ConfigMapInterface
	namespace       string
	jobImage        string
}

func NewClient(kubeclient kubernetes.Interface, namespace string, jobImage string) Client {
	return Client{
		cronjobClient:   kubeclient.BatchV1().CronJobs(namespace),
		secretClient:    kubeclient.CoreV1().Secrets(namespace),
		configMapClient: kubeclient.CoreV1().ConfigMaps(namespace),
		jobImage:        jobImage,
	}
}

func (c Client) Create(ctx context.Context, userID string, playlistID string, artists []string, token *oauth2.Token) error {
	secretName := tokenSecretName(userID)
	cmName := configmapName(playlistID)

	cronJobSpec := generateCronJob(jobSpec{
		jobImage:   c.jobImage,
		playlistID: playlistID,
		schedule:   everyFriday(),
		tokenSecretRef: corev1.LocalObjectReference{
			Name: secretName,
		},
		jobConfigMapRef: corev1.LocalObjectReference{
			Name: cmName,
		},
	})
	cj, err := c.cronjobClient.Create(ctx, cronJobSpec, v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating cronjob %v: %w", playlistID, err)
	}

	// jobConfig will be used by the job container to update the playlist
	jobConfig := Config{
		PlaylistID: playlistID,
		Artists:    artists,
	}
	err = c.storeJobConfig(ctx, cmName, jobConfig, cj.GetOwnerReferences())
	if err != nil {
		return fmt.Errorf("storing job config %v: %w", playlistID, err)
	}
	return nil
}

// Delete will delete the cronjob for the playlist
func (c Client) Delete(ctx context.Context, playlistID string) error {
	// token secret and configmap have their owner reference set to the cronjob, so they will be deleted by cascade.
	err := c.cronjobClient.Delete(ctx, playlistID, v1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("deleting cronjob for %v: %w", playlistID, err)
	}
	return err
}

func tokenSecretName(userID string) string {
	return fmt.Sprintf("%v-token", userID)
}

func configmapName(playlistID string) string {
	return fmt.Sprintf("%v-job-config", playlistID)
}

type jobSpec struct {
	jobImage        string
	playlistID      string
	schedule        string
	jobConfigMapRef corev1.LocalObjectReference
	tokenSecretRef  corev1.LocalObjectReference
}

func generateCronJob(s jobSpec) *batchv1.CronJob {
	jobConfigVolumeName := "job-config"
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
									VolumeMounts: []corev1.VolumeMount{
										{
											Name:      jobConfigVolumeName,
											MountPath: ConfigFilePath,
										},
									},
								},
							},
							Volumes: []corev1.Volume{
								{
									Name: jobConfigVolumeName,
									VolumeSource: corev1.VolumeSource{
										ConfigMap: &corev1.ConfigMapVolumeSource{
											LocalObjectReference: s.jobConfigMapRef,
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
