package nsx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	nsx_policy "github.com/vmware/vsphere-automation-sdk-go/services/nsxt"
)

var (
	host = "https://nsx01.cyberrange.rit.edu"
	// username = os.Getenv("GOVMOMI_USERNAME")
	// password = os.Getenv("GOVMOMI_PASSWORD")
)

type NSX struct {
	Client http.Client
	body   io.Reader
}

func (nsx *NSX) getNSXsegment(url string) (response *http.Response, err error) {
	//defining the http client
	// c :.Client{
	// 	CheckRedirect: redirectPolicyFunc,
	// }

	// response, err = nsx.Client.Get(url) //url = https://nsx01.cyberrange.rit.edu/api/v1/spec/openapi/nsx_api.json

	// building the REST API as an http request
	// this is the URL to be called in main - "https://nsx01.cyberrange.rit.edu/api/v1/spec/openapi/nsx_api.json"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	// default header - verify its requirement - if not required then remove it
	//request.Header.Add("If-None-Match", `W/"wyzzy"`)

	//add the request body from Postman - as JSON object or raw-form or JavaScript code depending on the query
	// request.Body.Add("")
	response, err = nsx.Client.Do(request) //get the response of the request
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	return
}

func (nsx *NSX) putNSXsegment(url string, body io.Reader) (response *http.Response, err error) {

	jsonReq, err := json.Marshal(body) //parse the body as a JSON request
	request, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonReq))

	if err != nil {
		log.Fatal(err)
	}
	response, err = nsx.Client.Do(request) //get the response of the request
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	//bodyBytes, _ := ioutil.ReadAll(response.Body)

	// Convert response body to string - This is to be done in the main method
	//bodyString := string(bodyBytes)

	return
}

func (nsx *NSX) patchNSXsegment(url string, body io.Reader) (response *http.Response, err error) {

	jsonReq, err := json.Marshal(body) //parse the body as a JSON request
	request, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonReq))

	if err != nil {
		log.Fatal(err)
	}
	response, err = nsx.Client.Do(request) //get the response of the request
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	return
}

func (nsx *NSX) deleteNSXsegment(url string, body io.Reader) (response *http.Response, err error) {

	jsonReq, err := json.Marshal(body) //parse the body as a JSON request
	request, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(jsonReq))

	if err != nil {
		log.Fatal(err)
	}
	response, err = nsx.Client.Do(request) //get the response of the request
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	return
}

// func (nsx *NSX) putNSXsegment(url string) (response *http.Response, err error) {
// 	// c := &http.Client{
// 	// 	CheckRedirect: redirectPolicyFunc,
// 	// }

// 	request, err := http.NewRequest(http.MethodPut, url, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	response, err = nsx.Client.Do(request)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	body, err := ioutil.ReadAll(response.Body)

// 	defer response.Body.Close()
// 	return
// }

// func (nsx *NSX) patchNSXsegment(url string, body io.Reader) (response *http.Response, err error) {
// 	// c := &http.Client{
// 	// 	CheckRedirect: redirectPolicyFunc,
// 	// }
// 	request, err := http.NewRequest(http.MethodPatch, url, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	response, err = nsx.Client.Do(request)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	body, err := ioutil.ReadAll(response.Body)

// 	defer response.Body.Close()

// 	return
// }

// func (nsx *NSX) deleteNSXsegment(url string) (response *http.Response, err error) {
// 	request, err := http.NewRequest(http.MethodDelete, url, nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	response, err = nsx.Client.Do(request)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer response.Body.Close()
//     bodyBytes, err := ioutil.ReadAll(response.Body)
// 	return
// }
func main() {
	// userPassContext := security.NewUserPasswordSecurityContext(username, password)
	// security.AuthenticationHandler.Authenticate(userPassContext)
	// test := security.AuthenticationHandler{}
	// sessions

	httpClient := http.Client{}
	connector := client.NewRestConnector(host, httpClient)
	infraClient := nsx_policy.NewDefaultInfraClient(connector)

	basePath := ""
	filter := ""
	typeFilter := ""
	infra, err := infraClient.Get(&basePath, &filter, &typeFilter)
	if err != nil {
		log.Fatalf("Error getting infra: %v\n", err)
	}
	fmt.Println(infra.DisplayName)

	nsxClient := NSX{
		Client: http.Client{},
	}

	//blue print for testing and adidng network and consuming teh REST API Call

	//Test GET
	response, err := nsxClient.getNSXsegment("some/url")
	fmt.Println(response.Status)

	//Test PUT
	// response, err := nsxClient.getNSXsegment("API end point URL",json_body)
	// fmt.Println(response.Status)

	//Test PATCH
	// response, err := nsxClient.getNSXsegment("API end point URL",json_body)
	// fmt.Println(response.Status)

	//Test DELETE
	// response, err := nsxClient.getNSXsegment("API end point URL",json_body)
	// fmt.Println(response.Status)

	// cfg := nsxt.Configuration{
	// 	BasePath:             "/api/v1",
	// 	Host:                 host,
	// 	Scheme:               "https",
	// 	UserName:             username,
	// 	Password:             password,
	// 	Insecure:             true,
	// }

	// client, err := nsxt.NewAPIClient(&cfg)
	// if err != nil {
	// 	log.Fatalf("error initializing NSX-T client: %v\n", err)
	// }

	// fmt.Println("NSX-T initialized...")

	// ctx := context.Background()
	// // client.Context

	// infra, res, err := client.PolicyApi.ReadInfra(ctx)
	// if err != nil {
	// 	log.Fatalf("Error reading infra: %v\n", err)
	// }
	// fmt.Printf("HTTP Status: %s\n", res.Status)
	// fmt.Println(infra)
}
