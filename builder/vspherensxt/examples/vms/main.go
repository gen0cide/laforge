package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gen0cide/laforge/builder/vspherensxt/vsphere"
	"github.com/sirupsen/logrus"
)

func main() {
	httpClient := http.Client{}

	baseUrl, urlExists := os.LookupEnv("VSPHERE_URL")
	username, usernameExists := os.LookupEnv("VSPHERE_USERNAME")
	password, passwordExists := os.LookupEnv("VSPHERE_PASSWORD")
	if !urlExists || !usernameExists || !passwordExists {
		log.Fatalf("please set VSPHERE_URL (exists? %t), VSPHERE_USERNAME (exists? %t), and VSPHERE_PASSWORD (exists? %t)", urlExists, usernameExists, passwordExists)
	}
	vs := vsphere.VSphere{
		HttpClient: httpClient,
		BaseUrl:    baseUrl,
		Username:   username,
		Password:   password,
	}

	vsphere.InitializeGovmomi(&vs, baseUrl, username, password)

	ctx := context.Background()

	vmList, err := vs.ListVms(ctx)
	if err != nil {
		log.Fatalf("error while getting vm's: %v", err)
	}
	for _, vm := range vmList {
		logrus.WithFields(logrus.Fields{
			"identifier": vm.Vm.Value,
			"powerState": vm.Runtime.PowerState,
		}).Info(vm.Config.Name)
	}
}
