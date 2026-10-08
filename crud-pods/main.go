package main

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

	// Delete existing pod
	fmt.Println("Checking if pod exists...")

	_, err = pods.Get(ctx, "test-pod-nginx", metav1.GetOptions{})

	if err == nil {
		fmt.Println("Pod already exists, deleting...")

		err = pods.Delete(
			ctx,
			"test-pod-nginx",
			metav1.DeleteOptions{},
		)
		if err != nil {
			panic(err)
		}

		// Wait until pod is actually gone
		for {
			_, err := pods.Get(
				ctx,
				"test-pod-nginx",
				metav1.GetOptions{},
			)

			if apierrors.IsNotFound(err) {
				fmt.Println("Pod completely deleted")
				break
			}

			if err != nil {
				panic(err)
			}

			time.Sleep(500 * time.Millisecond)
		}
	} else if !apierrors.IsNotFound(err) {
		panic(err)
	}

	// Create pod
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pod-nginx",
			Labels: map[string]string{
				"app": "backend",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx-pod",
					Image: "nginx",
					Ports: []corev1.ContainerPort{
						{
							ContainerPort: 80,
							HostPort:      8080,
						},
					},
				},
			},
		},
	}

	c, err := pods.Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"Pod created with ID: %s and name: %s\n",
		c.GetUID(),
		c.GetName(),
	)

	// Get pod
	fmt.Println("Searching for pod...")

	p, err := pods.Get(
		ctx,
		"test-pod-nginx",
		metav1.GetOptions{},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(p.GetName(), p.GetLabels())

	// Delete pod
	err = pods.Delete(
		ctx,
		c.GetName(),
		metav1.DeleteOptions{},
	)
	if err != nil {
		panic(err)
	}
}
