package k8s

import (
	"k8s.io/client-go/kubernetes"
	batchv1 "k8s.io/client-go/kubernetes/typed/batch/v1"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
)

func NewClientSet() (kubernetes.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// create the clientset
	clientset := kubernetes.NewForConfigOrDie(config)
	return clientset, nil
}

func NewCronJobsClient(clientset kubernetes.Interface, namespace string) batchv1.CronJobInterface {
	return clientset.BatchV1().CronJobs(namespace)
}

func NewSecretClient(clientset kubernetes.Interface, namespace string) corev1.SecretInterface {
	return clientset.CoreV1().Secrets(namespace)
}

func NewConfigMapClient(clientset kubernetes.Interface, namespace string) corev1.ConfigMapInterface {
	return clientset.CoreV1().ConfigMaps(namespace)
}
