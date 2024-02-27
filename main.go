package main

import (
	"context"
	"flag"
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

func GetConfig(path string) (*rest.Config, error) {
	kubeconfing := flag.String("kubeconfig", path, "")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfing)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return config, nil
}

func BuildClient(config *rest.Config) (*kubernetes.Clientset, error) {
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return client, nil
}

func GetJobs(client *kubernetes.Clientset) (*batchv1.JobList, error) {
	jobs, err := client.BatchV1().Jobs("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return jobs, nil
}

func ShowJobs(jobs *batchv1.JobList) {
	for _, job := range jobs.Items {
		fmt.Println(job.Name)
		fmt.Println(job.Spec)
		fmt.Println(job.Status)
	}
}

func CreateCronJobSpec(min string, hour string, dayMonth string, month string, dayWeek string) *batchv1.CronJob {
	cronJobSpec := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-cronjob",
			Namespace: "default",
		},
		Spec: batchv1.CronJobSpec{
			Schedule: fmt.Sprintf("%v %v %v %v %v", min, hour, dayMonth, month, dayWeek),
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:  "job-container",
									Image: "busybox",
									Command: []string{
										"/bin/sh",
										"-c",
										"echo Hello from the Kubernetes cluster",
									},
								},
							},
							RestartPolicy: corev1.RestartPolicyOnFailure,
						},
					},
				},
			},
		},
	}
	return cronJobSpec
}

func CreateJob(client *kubernetes.Clientset, jobSpec *batchv1.CronJob) (*batchv1.CronJob, error) {
	job, err := client.BatchV1().CronJobs("default").Create(context.Background(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Printf("Job %s created successfully\n", job.Name)
	return job, nil
}

func main() {
	path := "config"
	config, err := GetConfig(path)
	if err != nil {
		log.Fatalf("Error obtaining config: %v/n", err)
	}
	client, err := BuildClient(config)
	if err != nil {
		log.Fatalf("Error building client: %v/n", err)
	}
	jobs, err := GetJobs(client)
	if err != nil {
		log.Fatalf("Error obtaining jobs: %v/n", err)
	}
	ShowJobs(jobs)

	jobSpec := CreateCronJobSpec("0", "0", "*", "*", "4")

	_, er := CreateJob(client, jobSpec)
	if er != nil {
		log.Fatalf("Error creating CronJob: %v/n", err)
	}

}
