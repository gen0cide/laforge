package laforge

import "github.com/sethvargo/go-password/password"

func defaultCompetition(name string) *Competition {
	res, err := password.Generate(24, 5, 5, false, false)
	if err != nil {
		Logger.Errorf("Error generating random password: %v", err)
	}
	return &Competition{
		ID:           name,
		RootPassword: res,
		DNS: &DNS{
			ID:         "default",
			Type:       "disabled",
			RootDomain: "CHANGEME.LOCAL",
			DNSServers: []string{
				"127.0.0.1",
				"127.0.0.2",
			},
			NTPServers: []string{
				"127.0.0.1",
				"127.0.0.2",
			},
		},
		Remote: &Remote{
			ID:            "default",
			Type:          "disabled",
			Region:        "us-west-2",
			Key:           "AWS_API_KEY",
			Secret:        "AWS_API_SECRET",
			StateBucket:   "S3_BUCKET_FOR_STATE",
			StorageBucket: "S3_BUCKET_FOR_STORAGE",
		},
	}
}
