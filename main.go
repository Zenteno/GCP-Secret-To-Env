package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"gopkg.in/yaml.v3"
)

var config = flag.String("config", "config.yml", "Path of file that contains variables to inject")
var profile = flag.String("profile", "default", "Profile to load variables from")

type Secret struct {
	Variable string `yaml:"VARIABLE"`
	SecretId string `yaml:"SECRET_ID"`
}

func main() {
	flag.Parse()
	if _profile := os.Getenv("SECRET_PROFILE"); _profile != "" {
		*profile = _profile
	}
	if _config := os.Getenv("CONFIG_FILE"); _config != "" {
		*config = _config
	}
	variables := map[string][]Secret{}
	file, err := ioutil.ReadFile(*config)
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
	if _, ok := variables[*profile]; !ok {
		log.Fatalf("Profile %s is not present on config file\n", *profile)
	}
	for _, v := range variables[*profile] {
		accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
			Name: v.SecretId,
		}
		result, err := client.AccessSecretVersion(ctx, accessRequest)
		if err != nil {
			log.Fatalf("failed to access secret version: %v", err)
		}
		fmt.Printf("export %s=\"%s\"\n", v.Variable, result.Payload.GetData())
	}

}
