package main

import (
	"flag"
	"log"
	"strconv"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var workloadLabel = "k8s.smp.io/reroll-every"
var podTemplateAnnotation = "k8s.smp.io/rev"
var listOptions = metav1.ListOptions{
	LabelSelector: workloadLabel,
}
var clientset *kubernetes.Clientset = nil

func main() {
	flag.StringVar(&workloadLabel, "workload-label", workloadLabel, "")
	flag.StringVar(&podTemplateAnnotation, "pod-template-annotation", podTemplateAnnotation, "")
	flag.Parse()

	clientset = getClient()
	do(deployment{nil})
	do(statefulSet{nil})
}

func do(baseItem Workload) {
	workloads, err := baseItem.List()
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()

	for _, w := range workloads {
		schedule, err := time.ParseDuration(w.Meta().Labels[workloadLabel])
		if err != nil {
			log.Printf("Invalid duration on %s: %s", w, err)
			continue
		}
		rev, err := strconv.ParseInt(w.PodTemplate().Annotations[podTemplateAnnotation], 10, 64)
		if err != nil {
			rev = 0
		}

		lastRevTime := time.Unix(rev, 0)
		if now.Before(lastRevTime.Add(schedule)) {
			continue
		}

		log.Printf("Rerolling %s", w)
		podTemplate := w.PodTemplate()
		if podTemplate.Annotations == nil {
			podTemplate.Annotations = make(map[string]string)
		}
		podTemplate.Annotations[podTemplateAnnotation] = strconv.FormatInt(now.Unix(), 10)
		err = w.Update()
		if err != nil {
			log.Printf("Failed to update %s: %s", w, err)
		}
	}
}

/*
func doDeployments(clientset *kubernetes.Clientset) {
	now := time.Now()

	list, err := clientset.AppsV1().Deployments("").List(context.TODO(), listOptions)
	if err != nil {
		log.Fatal(err)
	}

	for _, deploy := range list.Items {
		schedule, err := time.ParseDuration(deploy.Labels[workloadLabel])
		if err != nil {
			log.Printf("Invalid duration on Deployment %s/%s: %s", deploy.Namespace, deploy.Name, err)
			continue
		}
		rev, err := strconv.ParseInt(deploy.Spec.Template.Annotations[podTemplateAnnotation], 10, 64)
		if err != nil {
			log.Printf("Invalid annotation on Deployment pod template %s/%s: %s", deploy.Namespace, deploy.Name, err)
			continue
		}

		lastRevTime := time.Unix(rev, 0)
		if now.Before(lastRevTime.Add(schedule)) {
			continue
		}

		log.Printf("Rerolling Deployment %s/%s", deploy.Namespace, deploy.Name)
		deploy.Spec.Template.Annotations[podTemplateAnnotation] = strconv.FormatInt(now.Unix(), 10)
		_, err = clientset.AppsV1().Deployments(deploy.Namespace).Update(context.TODO(), &deploy, metav1.UpdateOptions{})
		if err != nil {
			log.Printf("Failed to update Deployment %s/%s: %s", deploy.Namespace, deploy.Name, err)
		}
	}
}*/

func getClient() *kubernetes.Clientset {
	// config, err := rest.InClusterConfig()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	config := &rest.Config{Host: "localhost:8001"}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return clientset
}
