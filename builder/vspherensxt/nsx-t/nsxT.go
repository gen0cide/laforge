package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	host     = "https://nsx01.cyberrange.rit.edu"
	username = os.Getenv("GOVMOMI_USERNAME")
	password = os.Getenv("GOVMOMI_PASSWORD")
)

type NSX struct {
	Client http.Client
	body   NSXcreateSegmentPayload
}

type SegmentComponents struct {
	Status              string `json:"type"` //typpe instance of segment - referred as status
	Transport_zone_path string `json:"transport_zone_path"`
	Resource_type       string `json:"resource_type"`
	Id                  string `json:"id"`
	Display_name        string `json:"display_name"`
}

type NSXsegmentchildren struct {
	Resource_type     string            `json:"resource_type"`
	Marked_for_delete string            `json:"marked_for_delete"`
	Segment           SegmentComponents `json:"Segment"`
}

//JSON body for the PATCH REST API call - to create a new segment

type NSXcreateSegmentPayload struct {
	Resource_type string               `json:"resource_type"`
	Children      []NSXsegmentchildren `json:"children"`
}

type NSXpostDHCPpayload struct { //interface that defines the JSON body for adding a DHCP profile

	Display_name_post string `json:"display_name"`
	Edge_cluster_id   string `json:"edge_cluster_id"`
}

type NSXgatewayAddress struct {
	Gateway_address string `json:"gateway_address"`
}
type NSXaddSubnetPayload struct {
	Display_name_subnet string              `json:"display_name"`
	Subnets             []NSXgatewayAddress `json:"subnets"`
	Connectivity_path   string              `json:"connectivity_path"`
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

//change the body type with the interface for each REST API call independently
func (nsx *NSX) postNSXsegment(url string, body NSXpostDHCPpayload) (response *http.Response, err error) {

	jsonReq, err := json.Marshal(body) //parse the body as a JSON request
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonReq))

	if err != nil {
		log.Fatal(err)
	}
	request.SetBasicAuth(username, password)
	request.Header.Add("Content-Type", "application/json") //the key-line to parse the REST API call
	response, err = nsx.Client.Do(request)                 //get the response of the request
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()
	//bodyBytes, _ := ioutil.ReadAll(response.Body)

	// Convert response body to string - This is to be done in the main method
	//bodyString := string(bodyBytes)

	return
}

func (nsx *NSX) patchNSXsegment(url string, body NSXcreateSegmentPayload) (response *http.Response, err error) {

	jsonReq, err := json.Marshal(body) //parse the body as a JSON request
	fmt.Printf("%s", jsonReq)
	request, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonReq))

	if err != nil {
		log.Fatal(err)
	}
	request.SetBasicAuth(username, password)
	request.Header.Add("Content-Type", "application/json") //the key-line to parse the REST API call
	//add everytime after setting the basic auth

	response, err = nsx.Client.Do(request) //get the response of the request
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	return
}

func (nsx *NSX) patchNSXsubnet(url string, body NSXaddSubnetPayload) (response *http.Response, err error) {

	jsonReq, err := json.Marshal(body) //parse the body as a JSON request
	fmt.Printf("%s", jsonReq)
	request, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(jsonReq))

	if err != nil {
		log.Fatal(err)
	}
	request.SetBasicAuth(username, password)
	request.Header.Add("Content-Type", "application/json") //the key-line to parse the REST API call
	//add everytime after setting the basic auth

	response, err = nsx.Client.Do(request) //get the response of the request
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	return
}

func (nsx *NSX) deleteNSXsegment(url string, body NSXcreateSegmentPayload) (response *http.Response, err error) {

	jsonReq, err := json.Marshal(body) //parse the body as a JSON request
	request, err := http.NewRequest(http.MethodDelete, url, bytes.NewBuffer(jsonReq))

	if err != nil {
		log.Fatal(err)
	}

	request.SetBasicAuth(username, password)
	request.Header.Add("Content-Type", "application/json") //the key-line to parse the REST API call
	//add everytime after setting the basic auth

	response, err = nsx.Client.Do(request) //get the response of the request
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	return
}

func main() {
	// userPassContext := security.NewUserPasswordSecurityContext(username, password)
	// security.AuthenticationHandler.Authenticate(userPassContext)
	// test := security.AuthenticationHandler{}
	// sessions

	// httpClient := http.Client{}
	// connector := client.NewRestConnector(host, httpClient)
	// infraClient := nsx_policy.NewDefaultInfraClient(connector)

	// basePath := ""
	// filter := ""
	// typeFilter := ""
	// infra, err := infraClient.Get(&basePath, &filter, &typeFilter)
	// if err != nil {
	// 	log.Fatalf("Error getting infra: %v\n", err)
	// }
	// fmt.Println(infra.DisplayName)

	nsxClient := NSX{
		Client: http.Client{},
	}

	json_param_newSegment := NSXcreateSegmentPayload{
		Resource_type: "Infra",
		Children: []NSXsegmentchildren{

			NSXsegmentchildren{
				Resource_type:     "ChildSegment",
				Marked_for_delete: "false",
				Segment: SegmentComponents{
					Status:              "DISCONNECTED",
					Transport_zone_path: "/infra/sites/default/enforcement-points/default/transport-zones/dd8e4438-6593-4521-a164-668a60dd28c3",
					Resource_type:       "Segment",
					Id:                  "test-seg2",
					Display_name:        "test-seg2",
				},
			},
		},
	}

	//use different response variable for each REST API call
	response, err := nsxClient.patchNSXsegment("https://nsx01.cyberrange.rit.edu/policy/api/v1/infra", json_param_newSegment)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.Status)

	//TEST POST
	//parse the JSON body using the defined interface
	json_param_postDHCP := NSXpostDHCPpayload{
		Display_name_post: "test-dhcp-server",
		Edge_cluster_id:   "96155f90-10a1-4aae-a11f-34a8cc7cca79",
	}

	response_postDHCP, err := nsxClient.postNSXsegment("https://nsx01.cyberrange.rit.edu/api/v1/dhcp/servers", json_param_postDHCP)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response_postDHCP.Status)

	json_param_subnet := NSXaddSubnetPayload{
		Display_name_subnet: "test-sub1",
		Subnets: []NSXgatewayAddress{
			NSXgatewayAddress{
				Gateway_address: "172.19.2.254/24",
			},
		},
		Connectivity_path: "/infra/tier-1s/CPTC",
	}

	response_addSubnet, err := nsxClient.patchNSXsubnet("https://nsx01.cyberrange.rit.edu/api/v1/dhcp/servers", json_param_subnet)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response_addSubnet.Status)

	//TEST adding subnet to network segment through REST API

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
