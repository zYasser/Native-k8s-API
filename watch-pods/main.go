package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	pods := clientset.CoreV1().Pods("default")
	watcher, err := pods.Watch(ctx, v1.ListOptions{})
	if err != nil {
		panic(err)
	}
	defer watcher.Stop()
	for event := range watcher.ResultChan() {
		pod, ok := event.Object.(*corev1.Pod)
		if !ok {
			continue
		}

		fmt.Printf("Event: %-8s Pod: %-30s Phase: %s\n",
			event.Type,
			pod.Name,
			pod.Status.Phase,
		)

	}

}
