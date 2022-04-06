package main

import (
	"context"
	"fmt"

	sf "github.com/sa-/slicefunk"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Workload interface {
	List() ([]Workload, error)
	String() string
	PodTemplate() *v1.PodTemplateSpec
	Meta() *metav1.ObjectMeta
	Update() error
}

type deployment struct{ *appsv1.Deployment }
type statefulSet struct{ *appsv1.StatefulSet }
type daemonSet struct{ *appsv1.DaemonSet }
type replicaSet struct{ *appsv1.ReplicaSet }

func (w deployment) List() ([]Workload, error) {
	list, err := clientset.AppsV1().Deployments("").List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}
	return sf.Map(list.Items, func(w appsv1.Deployment) Workload { return deployment{&w} }), nil
}

func (w deployment) String() string {
	return fmt.Sprintf("Deployment %s/%s", w.Namespace, w.Name)
}

func (w deployment) PodTemplate() *v1.PodTemplateSpec {
	return &w.Spec.Template
}

func (w deployment) Meta() *metav1.ObjectMeta {
	return &w.ObjectMeta
}

func (w deployment) Update() error {
	_, err := clientset.AppsV1().Deployments(w.Namespace).Update(context.TODO(), w.Deployment, metav1.UpdateOptions{})
	return err
}

func (w statefulSet) List() ([]Workload, error) {
	list, err := clientset.AppsV1().StatefulSets("").List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}
	return sf.Map(list.Items, func(w appsv1.StatefulSet) Workload { return statefulSet{&w} }), nil
}

func (w statefulSet) String() string {
	return fmt.Sprintf("StatefulSet %s/%s", w.Namespace, w.Name)
}

func (w statefulSet) PodTemplate() *v1.PodTemplateSpec {
	return &w.Spec.Template
}

func (w statefulSet) Meta() *metav1.ObjectMeta {
	return &w.ObjectMeta
}

func (w statefulSet) Update() error {
	_, err := clientset.AppsV1().StatefulSets(w.Namespace).Update(context.TODO(), w.StatefulSet, metav1.UpdateOptions{})
	return err
}

func (w daemonSet) List() ([]Workload, error) {
	list, err := clientset.AppsV1().DaemonSets("").List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}
	return sf.Map(list.Items, func(w appsv1.DaemonSet) Workload { return daemonSet{&w} }), nil
}

func (w daemonSet) String() string {
	return fmt.Sprintf("DaemonSet %s/%s", w.Namespace, w.Name)
}

func (w daemonSet) PodTemplate() *v1.PodTemplateSpec {
	return &w.Spec.Template
}

func (w daemonSet) Meta() *metav1.ObjectMeta {
	return &w.ObjectMeta
}

func (w daemonSet) Update() error {
	_, err := clientset.AppsV1().DaemonSets(w.Namespace).Update(context.TODO(), w.DaemonSet, metav1.UpdateOptions{})
	return err
}

func (w replicaSet) List() ([]Workload, error) {
	list, err := clientset.AppsV1().ReplicaSets("").List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}
	return sf.Map(list.Items, func(w appsv1.ReplicaSet) Workload { return replicaSet{&w} }), nil
}

func (w replicaSet) String() string {
	return fmt.Sprintf("ReplicaSet %s/%s", w.Namespace, w.Name)
}

func (w replicaSet) PodTemplate() *v1.PodTemplateSpec {
	return &w.Spec.Template
}

func (w replicaSet) Meta() *metav1.ObjectMeta {
	return &w.ObjectMeta
}

func (w replicaSet) Update() error {
	_, err := clientset.AppsV1().ReplicaSets(w.Namespace).Update(context.TODO(), w.ReplicaSet, metav1.UpdateOptions{})
	return err
}
