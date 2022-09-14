package main

import (
	"context"
	"log"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	jobBaseName = "deploy-wand-job"
	namespaceName = "deploy-wand"
	grimoire = "grimoire"
	busybox = "busybox"
)

func applyJob(clientset *kubernetes.Clientset, grimoireVersion string) {
	jobName := jobBaseName + "-" + grimoireVersion
	var ttl int32 = 200
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name: jobName,
			Namespace: namespaceName,
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  busybox,
							Image: busybox, // + ":" + grimoireVersion,
							Command: []string{
								"sh",
								"-c",
								"busybox | head -1",
							},
						},
					},
					RestartPolicy: corev1.RestartPolicyNever,
				},
			},
			TTLSecondsAfterFinished: &ttl,
		},
	}

	log.Printf("Creating job: %v", job)
	_, err := clientset.BatchV1().Jobs(namespaceName).Get(context.TODO(), jobName, metav1.GetOptions{})
	if err != nil {
		_, err = clientset.BatchV1().Jobs(namespaceName).Create(context.TODO(), job, metav1.CreateOptions{})
		if err != nil {
			log.Println("Failed to create job " + jobName)
			panic(err.Error())
		}
		log.Println("Job " + jobName + " is created.")
	} else {
		log.Println("Job " + jobName + " is already created.")
	}
}