package main

import (
	"log"
	"testing"
)

func TestGetConfig(t *testing.T) {
	config, err := GetConfig()
	if err != nil {
		t.Errorf("GetConfig returned unexpected error %v", err)
	}
	if config == nil {
		t.Error("GetConfig returned nil config")
	}
}

func TestBuildClient(t *testing.T) {
	config, err := GetConfig()

	client, err := BuildClient(config)
	if err != nil {
		log.Fatalf("Unexpected Error: %v/n", err)
	}

	if client == nil {
		log.Fatal("Builder returned nil config")
	}
}

func TestCreateCronJobSpec(t *testing.T) {
	command := []string{
		"/bin/sh",
		"-c",
		"echo Hello from the Kubernetes cluster",
	}
	name := "my-cronjob"
	namespace := "default"
	jobSpec := CreateCronJobSpec(name, namespace, command, "0", "0", "*", "*", "4")
	if jobSpec == nil {
		t.Error("Unexpected error")
	}
}

func TestCreateJob(t *testing.T) {
	command := []string{
		"/bin/sh",
		"-c",
		"echo Hello from the Kubernetes cluster",
	}
	name := "my-cronjob"
	namespace := "default"

	config, _ := GetConfig()
	client, _ := BuildClient(config)
	jobSpec := CreateCronJobSpec(name, namespace, command, "0", "0", "*", "*", "4")
	cronJob, err := CreateJob(client, jobSpec)
	if err != nil {
		t.Errorf("Unexpected error: %v/n", err)
	}
	if cronJob == nil {
		t.Error("Unexpected behaviour. Job wasn't created")
	}
}

func TestDeleteJob(t *testing.T) {
	config, _ := GetConfig()
	client, _ := BuildClient(config)
	err := DeleteJob(client, "my-cronjob")
	if err != nil {
		t.Errorf("Error Deleting CronJob %v", err)
	}
}
