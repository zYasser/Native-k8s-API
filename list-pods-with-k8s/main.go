package main

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main(){
config, err := rest.InClusterConfig()
if err != nil {
	panic(err)
}
clientset, err := kubernetes.NewForConfig(config)
ctx :=context.Background()
result , err:= clientset.CoreV1().Pods("default").List(ctx , v1.ListOptions{})
if err != nil{
	panic(err)
}
fmt.Printf("Result kind :{} " , result.Kind)
for _, value :=  range result.Items{
	fmt.Printf("Kind : {} , Meta Data {}" , value.Kind, value.ObjectMeta.Name)
}
}