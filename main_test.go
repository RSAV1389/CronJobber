package main

import (
	"testing"
)

//func TestGetConfig(t *testing.T) {
//	path := "config"
//	config, err := GetConfig(path)
//	if err != nil {
//		t.Errorf("GetConfig returned unexpected error %v", err)
//	}
//	if config == nil {
//		t.Error("GetConfig returned nil config")
//	}
//}

//func TestGetConfig_Error(t *testing.T) {
//	path := "helldivers2/not/such/a/good/game/as/they/saying"
//	config, err := GetConfig(path)
//	if err == nil {
//		t.Error("GetConfig didn't return expected error")
//	}
//	if config != nil {
//		t.Error("GetConfig returned config BUT SHOUlDN't")
//	}
//}

//func TestBuildClient(t *testing.T) {
//	path := "config"
//	config, err := GetConfig(path)
//
//	client, err := BuildClient(config)
//	if err != nil {
//		log.Fatalf("Unexpected Error: %v/n", err)
//	}
//
//	if client == nil {
//		log.Fatal("Builder returned nil config")
//	}
//}

//	func TestBuildClient_Error(t *testing.T) {
//		path := "helldivers2/not/such/a/good/game/as/they/saying"
//		config, err := GetConfig(path)
//
//		client, err := BuildClient(config)
//		if err == nil {
//			log.Fatalf("Builder didn't return expected error")
//		}
//
//		if client != nil {
//			log.Fatal("Builder somehow returned client")
//		}
//	}
func TestCreateCronJobSpec(t *testing.T) {
	jobSpec := CreateCronJobSpec("0", "0", "*", "*", "4")
	if jobSpec == nil {
		t.Error("Unexpected error")
	}
}

func TestCreateJob(t *testing.T) {
	path := "config"
	config, _ := GetConfig(path)
	client, _ := BuildClient(config)
	jobSpec := CreateCronJobSpec("0", "0", "*", "*", "4")
	cronJob, err := CreateJob(client, jobSpec)
	if err != nil {
		t.Errorf("Unexpected error: %v/n", err)
	}
	if cronJob == nil {
		t.Error("Unexpected behaviour. Job wasn't created")
	}
}
