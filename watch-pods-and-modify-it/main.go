package main

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
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
		labels := pod.ObjectMeta.Labels
		fmt.Println(labels)
		value, ok := labels["app"]
		if !ok {
			fmt.Printf("This Object %s Does not have label app ", pod.GetName())
			continue
		}
		if value == "skip" {
			fmt.Printf("This Object %s has App label however it is marked as Skip ", pod.GetName())
			continue
		}
		patch := []byte(`{"metadata":{"labels": {"manager":"op"}}}`)
		result, err := pods.Patch(ctx, pod.GetName(), types.StrategicMergePatchType, patch, v1.PatchOptions{})
		if err != nil {
			fmt.Printf("Error patching pod %s: %v\n", pod.GetName(), err)
			return
		}

		fmt.Printf("Patched pod result: %+v\n", result)
		time.Sleep(1 * time.Second)
		p, err := pods.Get(
			ctx,
			pod.GetName(),
			v1.GetOptions{},
		)
		if err != nil {
			panic(err)
		}
		labels = p.ObjectMeta.Labels
		fmt.Printf("existed labels is %+v\n", labels)
		newValue, ok := labels["app"]
		if !ok {
			panic(fmt.Sprintf("The object did not have the old labels: %s", pod.GetName()))
		}
		if newValue != value {
			panic(fmt.Sprintf("The object has the old label but it did change: %s Old : %s and New: %s", pod.GetName(), value, newValue))

		}
		value, ok = labels["manager"]
		if !ok {
			panic(fmt.Sprintf("The object did not have the desired label: %s", pod.GetName()))
		}
		if value != "op" {
			panic(fmt.Sprintf("The object did not have the new label", pod.GetName()))

		}

	}

}
