package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gen0cide/laforge/builder/vsphere"
)

func main() {
	httpClient := http.Client{}

	baseUrl, urlExists := os.LookupEnv("VSPHERE_URL")
	username, usernameExists := os.LookupEnv("VSPHERE_USERNAME")
	password, passwordExists := os.LookupEnv("VSPHERE_PASSWORD")
	templateId, templateIdExists := os.LookupEnv("VSPHERE_TEMPLATE_ID")
	if !urlExists || !usernameExists || !passwordExists || !templateIdExists {
		log.Fatalf("please set VSPHERE_URL (exists? %t), VSPHERE_USERNAME (exists? %t), VSPHERE_PASSWORD (exists? %t), VSPHERE_TEMPLATE_ID (exists? %t)", urlExists, usernameExists, passwordExists, templateIdExists)
	}
	vshpere := vsphere.VSphere{
		Client: httpClient,
		BaseUrl: baseUrl,	
		Username: username,
		Password: password,
	}

	template, err := vshpere.GetTemplate(templateId)
	if err != nil {
		log.Fatalf("error while getting resource pools: %v", err)
	}
	fmt.Printf("%s [%s]\n", template.GuestOS, template.Identifier)
}