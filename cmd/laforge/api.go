package main

import (
	"fmt"
	"os"
)

func writeProvisionedHostConfig() {
	config := map[string]string{
		"HOST_ID":           "",
		"REMOTE_ADDR":       "",
		"SSH_HOSTNAME":      "",
		"SSH_PORT":          "",
		"SSH_USER":          "",
		"SSH_IDENTITY_FILE": "",
	}

	config = validateAPIConfigs(config)

}

func validateAPIConfigs(base map[string]string) map[string]string {
	ret := map[string]string{}

	for k := range base {
		key := fmt.Sprintf("LAFORGE_API_%s", k)
		val := os.Getenv(key)
		if val == "" {
			fmt.Printf("[LAFORGE API ERROR] MISSING API ENV PARAMETER: %s\n", key)
			os.Exit(1)
		}
		ret[k] = val
	}

	return ret
}
