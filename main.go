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
var podTemplateAnnotation = "k8s.smp.io/last-reroll"
var listOptions = metav1.ListOptions{
	LabelSelector: workloadLabel,
}
var clientset *kubernetes.Clientset = nil

func main() {
	flag.StringVar(&workloadLabel, "workload-label", workloadLabel, "")
	flag.StringVar(&podTemplateAnnotation, "pod-template-annotation", podTemplateAnnotation, "")
	flag.Parse()

	clientset = getClient()
	for {
		do(deployment{nil})
		do(statefulSet{nil})
		do(daemonSet{nil})
		do(replicaSet{nil})
		time.Sleep(1 * time.Minute)
	}
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

func getClient() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	return clientset
}
