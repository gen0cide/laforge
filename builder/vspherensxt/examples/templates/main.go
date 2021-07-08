package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gen0cide/laforge/builder/vspherensxt/vsphere"
)

func main() {
	httpClient := http.Client{}

	baseUrl, urlExists := os.LookupEnv("VSPHERE_URL")
	username, usernameExists := os.LookupEnv("VSPHERE_USERNAME")
	password, passwordExists := os.LookupEnv("VSPHERE_PASSWORD")
	libraryName, libraryNameExists := os.LookupEnv("VSPHERE_CONTENT_LIBRARY")
	templateName, templateNameExists := os.LookupEnv("VSPHERE_TEMPLATE_NAME")
	if !urlExists || !usernameExists || !passwordExists || !templateNameExists || !libraryNameExists {
		log.Fatalf("please set VSPHERE_URL (exists? %t), VSPHERE_USERNAME (exists? %t), VSPHERE_PASSWORD (exists? %t), VSPHERE_CONTENT_LIBRARY (exists? %t), VSPHERE_TEMPLATE_NAME (exists? %t)", urlExists, usernameExists, passwordExists, libraryNameExists, templateNameExists)
	}
	vshpere := vsphere.VSphere{
		HttpClient: httpClient,
		BaseUrl:    baseUrl,
		Username:   username,
		Password:   password,
	}

	templateId, err := vshpere.GetTemplateIDByName(libraryName, templateName)
	if err != nil {
		log.Fatalf("error while searching content library for \"%s\": %v", templateName, err)
	}
	template, err := vshpere.GetTemplate(templateId)
	if err != nil {
		log.Fatalf("error while getting resource pools: %v", err)
	}
	fmt.Printf("%s [%s]\n", template.GuestOS, template.Identifier)
}
