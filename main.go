package main

import (
	"log"
	"os"
	"time"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	for {
		versionConfig, err := os.ReadFile("/grimoire-config/tag")
		if err != nil {
			panic(err.Error())
		}
		grimoireVersion := string(versionConfig)
		
		log.Println("Grimoire version: " + grimoireVersion)
		log.Println("Applying job...")
		applyJob(clientset, grimoireVersion)
		log.Println("Job applied.")
		time.Sleep(10 * time.Second)
	}
}