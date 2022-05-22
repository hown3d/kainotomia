package k8s

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func NewFakeClienSet(obj ...runtime.Object) kubernetes.Interface {
	return fake.NewSimpleClientset(obj...)
}

func ToObjects[K runtime.Object](objs []K) []runtime.Object {
	runtimeObjs := make([]runtime.Object, len(objs))
	for _, obj := range objs {
		runtimeObjs = append(runtimeObjs, obj)
	}
	return runtimeObjs
}
