package main

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/hown3d/kainotomia/pkg/k8s/job"
	"gopkg.in/yaml.v3"
)

func main() {
	file, err := os.Open(filepath.Join(job.ConfigFilePath, job.ConfigFileName))
	if err != nil {
		log.Fatalf("opening config file: %v", err)
	}
	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("reading data from config file: %v", err)
	}
	cfg := new(job.Config)
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		log.Fatalf("unmarshaling job data: %v", err)
	}
	// TODO: create kubernetes config and load secret with token
}
