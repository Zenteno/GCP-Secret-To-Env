package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"gopkg.in/yaml.v3"
)

var file = flag.String("variables", "variables.yml", "Path of file that contains variables to inject")

type Secret struct {
	Variable string `yaml:"VARIABLE"`
	SecretId string `yaml:"SECRET_ID"`
}

func main() {
	flag.Parse()
	variables := []Secret{}
	file, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatalf("Could not read the file due to this %s error \n", err)
	}
	err = yaml.Unmarshal(file, &variables)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
	}
	defer client.Close()

	for _, v := range variables {
		accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
			Name: v.SecretId,
		}
		result, err := client.AccessSecretVersion(ctx, accessRequest)
		if err != nil {
			log.Fatalf("failed to access secret version: %v", err)
		}
		fmt.Printf("%s=\"%s\"\n", v.Variable, result.Payload.GetData())
	}

}
