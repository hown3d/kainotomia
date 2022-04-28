package job

import (
	"context"
	"fmt"

	"gopkg.in/yaml.v3"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ConfigFilePath string = "/kainotomia/job"
	ConfigFileName string = "config.yml"
)

type Config struct {
	PlaylistID string
	Artists    []string
}

func (c Client) storeJobConfig(ctx context.Context, name string, jobconfig Config, references []v1.OwnerReference) error {
	data, err := yaml.Marshal(jobconfig)
	if err != nil {
		return fmt.Errorf("marshaling job config: %w", err)
	}
	_, err = c.configMapClient.Create(ctx, &corev1.ConfigMap{
		ObjectMeta: v1.ObjectMeta{
			Name: name,
		},
		Data: map[string]string{
			ConfigFileName: string(data),
		},
	}, v1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating job config map %v: %w", jobconfig.PlaylistID, err)
	}
	return nil
}
