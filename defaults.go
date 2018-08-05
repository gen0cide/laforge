package laforge

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/sethvargo/go-password/password"
)

var (
	// ExampleObjects holds a map of example objects
	ExampleObjects = map[string]interface{}{
		"ami":         defaultAMI(),
		"command":     defaultCommand(),
		"dns_record":  defaultDNSRecord(),
		"identity":    defaultIdentity(),
		"network":     defaultNetwork(),
		"remote_file": defaultRemoteFile(),
		"script":      defaultScript(),
		"host":        defaultHost(),
		"environment": defaultEnvironment(),
	}
)

// ExampleObjectByName retrieves a templated example object of type (name)
func ExampleObjectByName(name string) ([]byte, error) {
	obj, found := ExampleObjects[name]
	if !found {
		return []byte{}, fmt.Errorf("requested object type is unrecognized: %s", name)
	}
	exampledef, err := RenderHCLv2Object(obj)
	if err != nil {
		return []byte{}, errors.WithMessage(errors.WithStack(err), "failed to render the example object")
	}
	return exampledef, nil
}

func defaultDNSRecord() *DNSRecord {
	return &DNSRecord{
		ID:         "example_dns_record_config",
		Name:       "www",
		Type:       "CNAME",
		Value:      "foo01.bar.com",
		Disabled:   true,
		OnConflict: defaultOnConflict(),
	}
}

func defaultRemoteFile() *RemoteFile {
	return &RemoteFile{
		ID:          "example_remote_file_config",
		SourceType:  "local",
		Source:      "/a/nonexist/path/auth_keys",
		Destination: "/root/.ssh/authorized_keys",
		Perms:       "0600",
		Disabled:    true,
		OnConflict:  defaultOnConflict(),
	}
}

func defaultEnvironment() *Environment {
	return &Environment{
		ID:          "fake_environment",
		Name:        "fake_env_for_demo_purposes",
		Description: "not a real environment, please configure!",
		Type:        "native",
		Config: map[string]string{
			"combine_scripts": "true",
		},
		Vars: map[string]string{
			"no_users": "true",
		},
		Tags: map[string]string{
			"example": "yes",
		},
		Maintainer: defaultMaintainer(),
		OnConflict: defaultOnConflict(),
		HostByNetwork: map[string][]*Host{
			"testnet": []*Host{
				defaultHost(),
			},
		},
	}
}

func defaultHost() *Host {
	return &Host{
		ID:           "example_host_config",
		Hostname:     "example01",
		Description:  "example host configuration, do not use!",
		OS:           "ubuntu",
		LastOctet:    187,
		InstanceSize: "xlarge",
		Disk: Disk{
			Size: 250,
		},
		ProvisionSteps: []string{
			"example_a_record",
			"example_cname_record",
			"example_service_setup_script",
			"example_file_a_config",
			"example_file_b_config",
			"example_useradd_script",
			"example_restart_command",
		},
		ExposedTCPPorts: []string{
			"9910",
			"1-1024",
		},
		ExposedUDPPorts: []string{
			"53",
		},
		OverridePassword: "veryweak123",
		IO:               defaultIO(),
		Maintainer:       defaultMaintainer(),
		OnConflict:       defaultOnConflict(),
		Vars: map[string]string{
			"env_specific_val": "yes",
		},
		Tags: map[string]string{
			"used_for": "linux_servers",
		},
	}
}

func defaultScript() *Script {
	return &Script{
		ID:           "example_script_config",
		Name:         "example_setup",
		Description:  "this script is a NOOP basic example of how to write a script config",
		Maintainer:   defaultMaintainer(),
		Language:     "shell",
		SourceType:   "local",
		Source:       "../relative/path/to/config/is/valid.sh",
		Cooldown:     10,
		IgnoreErrors: true,
		Args: []string{
			"--setup",
			"--delete-afterwards",
		},
		IO:       defaultIO(),
		Disabled: true,
		Vars: map[string]string{
			"restrict_to_flavor": "ubuntu",
			"use_sudo":           "true",
		},
		Tags: map[string]string{
			"used_for": "linux_servers",
		},
		OnConflict: defaultOnConflict(),
	}
}

func defaultIdentity() *Identity {
	return &Identity{
		ID:          "example_identity_config",
		Firstname:   "Bruce",
		Lastname:    "Wayne",
		Email:       "bruce@wayneenterprises.com",
		Password:    "changeme123",
		Description: "CEO of Wayne Enterprises",
		AvatarFile:  "",
		Vars: map[string]string{
			"title":       "CEO",
			"ou":          "Executives",
			"should_do_x": "true",
		},
		OnConflict: defaultOnConflict(),
	}
}

func defaultNetwork() *Network {
	return &Network{
		ID:         "example_network_config",
		Name:       "testnet",
		CIDR:       "192.168.22.0/24",
		VDIVisible: false,
		Vars: map[string]string{
			"aws":      "enabled",
			"vagrant":  "disabled",
			"use_dhcp": "false",
		},
		Tags: map[string]string{
			"difficulty": "hard",
			"location":   "offsite",
		},
		OnConflict: defaultOnConflict(),
	}
}

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

func defaultAMI() *AMI {
	return &AMI{
		ID:          "example_ami_config",
		Name:        "Example AMI Config",
		Description: "example AMI config block - should not be used!",
		Provider:    "aws",
		Username:    "ubuntu",
		Config: map[string]string{
			"ami_id":        "ami40abc1140",
			"region":        "us-west-2",
			"requires_sudo": "true",
			"access_method": "ssh_key",
			"ssh_key":       "<<< fake data here >>>",
		},
		Tags: map[string]string{
			"based_on": "ubuntu16.04",
			"arch":     "amd64",
		},
		Maintainer: defaultMaintainer(),
	}
}

func defaultCommand() *Command {
	return &Command{
		ID:          "example_command_config",
		Name:        "example command configuration block",
		Description: "for training purposes only!",
		Program:     "/usr/bin/ruby",
		Args: []string{
			"-e",
			"puts 'foo'.gsub(/foo/,'bar')",
		},
		IgnoreErrors: false,
		Cooldown:     0,
		IO:           defaultIO(),
		Disabled:     true,
		OnConflict:   defaultOnConflict(),
		Maintainer:   defaultMaintainer(),
	}
}

func defaultOnConflict() OnConflict {
	return OnConflict{
		Do:     "default",
		Append: false,
	}
}

func defaultIO() IO {
	return IO{
		Stdin:  "",
		Stdout: "/dev/null",
		Stderr: "/tmp/example.log",
	}
}

func defaultMaintainer() *User {
	return &User{
		ID:    "gen0cide",
		Name:  "Alex Levinson",
		Email: "alexl@uber.com",
	}
}
