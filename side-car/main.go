package main

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"

	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func generateServiceFromDeployment(deployment *apps.Deployment) *corev1.Service {
	fmt.Printf("[INFO] Generating Service for Deployment: %s\n", deployment.GetName())

	template := deployment.Spec.Template
	deploymentLabels := template.ObjectMeta.Labels

	fmt.Printf("[DEBUG] Deployment %s template labels: %+v\n", deployment.GetName(), deploymentLabels)

	serviceLabel := map[string]string{
		"app": "service",
	}

	service := &corev1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:   fmt.Sprintf("%s-service", deployment.GetName()),
			Labels: serviceLabel,
		},
		Spec: corev1.ServiceSpec{
			Selector: deploymentLabels,
			Type:     corev1.ServiceTypeNodePort,
			Ports: []corev1.ServicePort{
				corev1.ServicePort{
					Protocol: corev1.ProtocolTCP,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: template.Spec.Containers[0].Ports[0].ContainerPort,
					},
					Port: 80,
				},
			},
		},
	}

	fmt.Printf("[DEBUG] Generated Service: name=%s, type=%s, selector=%+v, targetPort=%d\n",
		service.GetName(),
		service.Spec.Type,
		service.Spec.Selector,
		service.Spec.Ports[0].TargetPort.IntVal,
	)

	return service
}

func main() {
	fmt.Println("[INFO] Starting Kubernetes Service controller...")

	config, err := rest.InClusterConfig()
	if err != nil {
		fmt.Printf("[ERROR] Failed to load in-cluster configuration: %v\n", err)
		panic(err)
	}
	fmt.Println("[INFO] Loaded in-cluster Kubernetes configuration")

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("[ERROR] Failed to create Kubernetes client: %v\n", err)
		panic(err)
	}
	fmt.Println("[INFO] Kubernetes client created successfully")

	ctx := context.Background()
	deploy := clientset.AppsV1().Deployments("default")

	fmt.Println("[INFO] Starting watch on Pods in namespace default...")
	watcher, err := deploy.Watch(ctx, v1.ListOptions{})
	if err != nil {
		fmt.Printf("[ERROR] Failed to start Pod watch: %v\n", err)
		panic(err)
	}
	defer watcher.Stop()
	fmt.Println("[INFO] Pod watch started successfully")

	for event := range watcher.ResultChan() {
		fmt.Printf("[DEBUG] Received event: type=%s, objectType=%T\n", event.Type, event.Object)

		d, ok := event.Object.(*apps.Deployment)
		if !ok {
			fmt.Printf("[DEBUG] Ignoring event because object is not a Deployment: %T\n", event.Object)
			continue
		}

		fmt.Printf("[INFO] Deployment event received: name=%s, eventType=%s\n", d.GetName(), event.Type)

		if event.Type != watch.Added {
			fmt.Printf("[DEBUG] Ignoring Deployment %s because event type is %s, not Added\n", d.GetName(), event.Type)
			continue
		}

		labels := d.ObjectMeta.Labels
		fmt.Printf("[DEBUG] Deployment %s labels: %+v\n", d.GetName(), labels)

		value, ok := labels["side-car"]
		if !ok {
			fmt.Printf("[WARN] Deployment %s does not have label side-car\n", d.GetName())
			continue
		}

		fmt.Printf("[INFO] Deployment %s has side-car=%s\n", d.GetName(), value)

		if value == "skip" {
			fmt.Printf("[INFO] Deployment %s has side-car label but is marked to skip\n", d.GetName())
			continue
		}

		serviceI := clientset.CoreV1().Services("default")

		fmt.Printf("[INFO] Creating Service for Deployment %s...\n", d.GetName())
		newService := generateServiceFromDeployment(d)

		fmt.Printf("[DEBUG] Service creation request: %+v\n", newService)

		createdService, err := serviceI.Create(ctx, newService, v1.CreateOptions{})
		if err != nil {
			fmt.Printf("[ERROR] Failed to create Service for Deployment %s: %v\n", d.GetName(), err)
		} else {
			fmt.Printf("[INFO] Successfully created Service: name=%s, namespace=%s\n",
				createdService.GetName(),
				createdService.GetNamespace(),
			)
		}

		// patch := []byte(`{"metadata":{"labels": {"manager":"op"}}}`)
		// result, err := pods.Patch(ctx, pod.GetName(), types.StrategicMergePatchType, patch, v1.PatchOptions{})
		// if err != nil {
		// 	fmt.Printf("Error patching pod %s: %v\n", pod.GetName(), err)
		// 	return
		// }

		// fmt.Printf("Patched pod result: %+v\n", result)
		// time.Sleep(1 * time.Second)
		// p, err := pods.Get(
		// 	ctx,
		// 	pod.GetName(),
		// 	v1.GetOptions{},
		// )
		// if err != nil {
		// 	panic(err)
		// }
		// labels = p.ObjectMeta.Labels
		// fmt.Printf("existed labels is %+v\n", labels)
		// newValue, ok := labels["app"]
		// if !ok {
		// 	panic(fmt.Sprintf("The object did not have the old labels: %s", pod.GetName()))
		// }
		// if newValue != value {
		// 	panic(fmt.Sprintf("The object has the old label but it did change: %s Old : %s and New: %s", pod.GetName(), value, newValue))

		// }
		// value, ok = labels["manager"]
		// if !ok {
		// 	panic(fmt.Sprintf("The object did not have the desired label: %s", pod.GetName()))
		// }
		// if value != "op" {
		// 	panic(fmt.Sprintf("The object did not have the new label", pod.GetName()))

		// }

	}

	fmt.Println("[WARN] Pod watch channel closed; controller stopped receiving events")
}
