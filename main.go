package main

import (
	"context"
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os/user"
	"path/filepath"
)

func GetConfig() (*rest.Config, error) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(usr.HomeDir, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", path)
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

func GetCronJobs(client *kubernetes.Clientset) (*batchv1.CronJobList, error) {
	jobs, err := client.BatchV1().CronJobs("").List(context.Background(), metav1.ListOptions{})
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

func ShowCronJobs(jobs *batchv1.CronJobList) {
	for _, job := range jobs.Items {
		fmt.Println(job.Name)
		fmt.Println(job.Spec)
		fmt.Println(job.Status)
	}
}

func CreateCronJobSpec(name string, namespace string, command []string, min string, hour string, dayMonth string, month string, dayWeek string) *batchv1.CronJob {
	cronJobSpec := &batchv1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: batchv1.CronJobSpec{
			Schedule: fmt.Sprintf("%v %v %v %v %v", min, hour, dayMonth, month, dayWeek),
			JobTemplate: batchv1.JobTemplateSpec{
				Spec: batchv1.JobSpec{
					Template: corev1.PodTemplateSpec{
						Spec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Name:    "job-container",
									Image:   "busybox",
									Command: command,
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

func CreateCronJob(client *kubernetes.Clientset, jobSpec *batchv1.CronJob) (*batchv1.CronJob, error) {
	job, err := client.BatchV1().CronJobs("default").Create(context.Background(), jobSpec, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Printf("Job %s created successfully\n", job.Name)
	return job, nil
}

func DeleteCronJob(client *kubernetes.Clientset, name string) error {
	err := client.BatchV1().CronJobs("default").Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		log.Fatalf("Error Deleting CronJob %v", err)
		return err
	}
	log.Print("Deleting job seems successful")
	return nil
}

func SuspendCronJob(client *kubernetes.Clientset, name string) error {
	cronJob, err := client.BatchV1().CronJobs("default").Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		log.Fatalf("Error Getting CronJob %v", err)
		return err
	}

	cronJob.Spec.Suspend = new(bool)
	_, err = client.BatchV1().CronJobs("default").Update(context.Background(), cronJob, metav1.UpdateOptions{})
	if err != nil {
		log.Fatalf("Error Suspending CronJob %v", err)
		return err
	}

	log.Print("Suspending job seems successful")
	return nil
}

func main() {
	config, err := GetConfig()
	if err != nil {
		log.Fatalf("Error obtaining config: %v/n", err)
	}

	client, err := BuildClient(config)
	if err != nil {
		log.Fatalf("Error building client: %v/n", err)
	}

	//jobs, err := GetJobs(client)
	//if err != nil {
	//	log.Fatalf("Error obtaining jobs: %v/n", err)
	//}

	//ShowJobs(jobs)

	cronJ, err := GetCronJobs(client)
	ShowCronJobs(cronJ)

	//command := []string{
	//	"/bin/sh",
	//	"-c",
	//	"echo Hello from the Kubernetes cluster",
	//}
	//
	//name := "my-cronjob"
	//namespace := "default"
	//
	//jobSpec := CreateCronJobSpec(name, namespace, command, "0", "0", "*", "*", "4")
	//
	//_, er := CreateCronJob(client, jobSpec)
	//
	//if er != nil {
	//	log.Fatalf("Error creating CronJob: %v/n", err)
	//}
	//
	//ero := DeleteJob(client, "my-cronjob")
	//if ero != nil {
	//	log.Fatalf("Error Deleting CronJob: %v/n", err)
	//}

}
