package graph

//go:generate go run github.com/99designs/gqlgen generate
import "github.com/gen0cide/laforge/graphql/graph/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver Specify all the options that are able to be resolved here
type Resolver struct {
	enviroments           []*model.Environment
}

// TestEnvironment ..
func TestEnvironment(testStatus *model.Status ) (*model.Environment,[]*model.ProvisionedNetwork,[]*model.ProvisionedHost,[]*model.ProvisionedStep) {
	// Test Simple Configs
	testCidr := "129.21.434.0/23"
	testSSHPort := "22"
	testRDPPort := "3519"
	testDNSServerA := "8.8.8.8"
	testDNAServerB := "8.8.4.4"
	testNTPServerA := "129.6.15.28"
	testNTPServerB := "129.6.15.29"
	testFloat := 1.5
	testTag := model.Tag{
		ID:          "test",
		Name:        "This is a Test Tag",
		Description: new(string),
	}
	textConfig := model.ConfigMap{
		Key:   "TestKey",
		Value: "TestValue",
	}
	testUser := model.User{
		ID:    "frybin",
		Name:  "Fred Rybin",
		UUID:  "dec27c98-876f-4d51-93e6-ac22bdc3cf55",
		Email: "fxr7632@rit.edu",
	}
	testVars := model.VarsMap{
		Key:   "Testing",
		Value: "True",
	}

	// Test Networks
	testNetwork1 := model.Network{
		ID:         "7fbc05c9-2e62-45d9-947f-43eae7660aba",
		Name:       "Test Network 1",
		Cidr:       "10.0.0.0/24",
		VdiVisible: false,
		Tags:       []*model.Tag{&testTag},
	}
	testNetwork2 := model.Network{
		ID:         "149e4a5a-7a18-4ea4-96a8-2b9d649d00dd",
		Name:       "Test Network 1",
		Cidr:       "10.0.10.0/24",
		VdiVisible: false,
		Vars:       []*model.VarsMap{&testVars},
	}
	testNetwork3 := model.Network{
		ID:         "0f37382e-cf72-428b-8e28-bd445023117c",
		Name:       "Test Network 1",
		Cidr:       "10.0.100.0/24",
		VdiVisible: true,
	}

	//Step Definitions
	testDNSRecord := model.DNSRecord{
		ID:       "6d178a21-09f5-4692-ad3d-f868712fb254",
		Name:     "my",
		Values:   []*string{},
		Type:     "A",
		Zone:     "test.com",
		Vars:     []*model.VarsMap{},
		Tags:     []*model.Tag{},
		Disabled: false,
	}
	testCommand := model.Command{
		ID:           "1c675cf8-a995-4669-9fa5-3e8a7fc193a6",
		Name:         "Show Directory Content",
		Description:  "Show Current Directory Contents",
		Program:      "ls",
		Args:         []*string{},
		IgnoreErrors: false,
		Cooldown:     0,
		Timeout:      0,
		Disabled:     false,
		Vars:         []*model.VarsMap{},
		Tags:         []*model.Tag{},
		Maintainer:   &testUser,
	}
	testScript := model.Script{
		ID:           "930cec53-5d82-4685-9ef5-d4d93c76b1a1",
		Name:         "Generate 'testUser' user",
		Language:     "shell",
		Description:  "Generates a user named testUser",
		Source:       "local",
		SourceType:   "./adduser-testUser.sh",
		Cooldown:     0,
		Timeout:      0,
		IgnoreErrors: false,
		Args:         []*string{},
		Disabled:     false,
		Vars:         []*model.VarsMap{},
		Tags:         []*model.Tag{},
		AbsPath:      "",
		Maintainer:   &testUser,
		Findings:     []*model.Finding{{
			Name:        "default local admin password",
			Description: "",
			Severity:    model.FindingSeverityMediumSeverity,
			Difficulty:  model.FindingDifficultyNoviceDifficulty,
			Maintainer:  &model.User{},
			Tags:        []*model.Tag{},
			Host:        &model.Host{},
		}},
	}
	testFileDelete := model.FileDelete{
		ID:   "e7f7e4b5-2854-458a-89af-72b53cd678c8",
		Path: "/tmp/test.conf",
	}
	testFileDownload := model.FileDownload{
		ID:          "d2da9b50-1e96-4f9f-a068-155c0616813e",
		SourceType:  "local",
		Source:      "./test.conf",
		Destination: "/tmp/test.conf",
		Templete:    false,
		Mode:        "0777",
		Disabled:    false,
		Md5:         "387337cd9c394c04fb5a33c176cf8715",
		AbsPath:     "",
		Tags:        []*model.Tag{},
	}
	testFileExtract := model.FileExtract{
		ID:          "ba02a4da-3ee9-4b34-8a81-e1e7b2e92c57",
		Source:      "/tmp/test.zip",
		Destination: "/tmp",
		Type:        "zip",
	}

	// Test Hosts
	testWindowsHost := model.Host{
		ID:               "0bee3405-7482-4b88-8c59-a94059e107b9",
		Hostname:         "test.windows.host",
		Os:               "wk12",
		LastOctet:        12,
		AllowMacChanges:  false,
		ExposedTCPPorts:  []*string{&testRDPPort},
		ExposedUDPPorts:  []*string{},
		OverridePassword: "TestPass123",
		Vars:             []*model.VarsMap{},
		UserGroups:       []*string{},
		DependsOn:        []*model.Host{},
		Disk: 			  &model.Disk{
			Size: 50,
		},
		Maintainer:       &testUser,
		Tags:             []*model.Tag{&testTag},
		FileDownloads:    []*model.FileDownload{&testFileDownload},
		FileExtracts:     []*model.FileExtract{&testFileExtract},
	}
	testLinuxHost1 := model.Host{
		ID:               "1e1095ca-c42c-4ef7-a9ef-d58d54e9bcb7",
		Hostname:         "test1.linux.host",
		Os:               "Ubuntu18",
		LastOctet:        10,
		AllowMacChanges:  true,
		ExposedTCPPorts:  []*string{&testSSHPort},
		ExposedUDPPorts:  []*string{},
		OverridePassword: "TestPass12",
		Vars:             []*model.VarsMap{},
		UserGroups:       []*string{},
		DependsOn:        []*model.Host{},
		Disk: 			  &model.Disk{
			Size: 50,
		},
		Maintainer:       &testUser,
		Tags:             []*model.Tag{&testTag},
		DNSRecords:       []*model.DNSRecord{&testDNSRecord},
		Commands:         []*model.Command{&testCommand},
		Scripts:          []*model.Script{&testScript},
	}
	testLinuxHost2 := model.Host{
		ID:               "ef151b9c-e8ba-4df9-81ff-90ca678b67f8",
		Hostname:         "test2.linux.host",
		Os:               "Ubuntu18",
		LastOctet:        15,
		AllowMacChanges:  true,
		ExposedTCPPorts:  []*string{&testSSHPort},
		ExposedUDPPorts:  []*string{},
		OverridePassword: "TestPass12",
		Vars:             []*model.VarsMap{},
		UserGroups:       []*string{},
		DependsOn:        []*model.Host{&testLinuxHost1},
		Disk: 			  &model.Disk{
			Size: 50,
		},
		Maintainer:       &testUser,
		Tags:             []*model.Tag{},
		Commands:         []*model.Command{&testCommand},
		Scripts:          []*model.Script{&testScript},
		FileDeletes:      []*model.FileDelete{&testFileDelete},
		FileDownloads:    []*model.FileDownload{&testFileDownload},
	}

	// Test Provisioned Hosts
	testProHost1 := model.ProvisionedHost{
		ID:                 "ae3ced80-6385-4421-8a75-9a5ea464c62b",
		SubnetIP:           "10.0.0.12",
		Status:             testStatus,
		Host:               &testWindowsHost,
		Heartbeat: &model.AgentStatus{
			ClientID: "ae3ced80-6385-4421-8a75-9a5ea464c62b",
			Hostname: "test.windows.host",
			UpTime:   708,
			BootTime: 1608158324,
			NumProcs: 10,
			Os:       "Windows",
			HostID:   "???",
			TotalMem: 16426635264,
			FreeMem:  8437784576,
			UsedMem:  3781406720,
		},
	}
	
	testProHost2 := model.ProvisionedHost{
		ID:                 "10fd9876-fc42-453d-ab93-c02bfea696bb",
		SubnetIP:           "10.0.10.10",
		Status:             testStatus,
		Host:               &testLinuxHost1,
		Heartbeat: &model.AgentStatus{
			ClientID: "10fd9876-fc42-453d-ab93-c02bfea696bb",
			Hostname: "test1.linux.host",
			UpTime:   708,
			BootTime: 1608158324,
			NumProcs: 10,
			Os:       "Linux",
			HostID:   "???",
			Load1:    &testFloat,
			Load5:    &testFloat,
			Load15:   &testFloat,
			TotalMem: 16426635264,
			FreeMem:  8437784576,
			UsedMem:  3781406720,
		},
	}

	testProHost3 := model.ProvisionedHost{
		ID:                 "51a37e0b-0223-4da5-a613-aac4128f80de",
		SubnetIP:           "10.0.10.15",
		Status:             testStatus,
		Host:               &testLinuxHost2,
		Heartbeat: &model.AgentStatus{
			ClientID: "51a37e0b-0223-4da5-a613-aac4128f80de",
			Hostname: "test2.linux.host",
			UpTime:   708,
			BootTime: 1608158324,
			NumProcs: 10,
			Os:       "Linux",
			HostID:   "???",
			Load1:    &testFloat,
			Load5:    &testFloat,
			Load15:   &testFloat,
			TotalMem: 16426635264,
			FreeMem:  8437784576,
			UsedMem:  3781406720,
		},
	}

	testProHost4 := model.ProvisionedHost{
		ID:                 "91e3c993-e076-4989-8445-2ba895fc1edd",
		SubnetIP:           "10.0.100.12",
		Status:             testStatus,
		Host:               &testWindowsHost,
		Heartbeat: &model.AgentStatus{
			ClientID: "91e3c993-e076-4989-8445-2ba895fc1edd",
			Hostname: "test.windows.host",
			UpTime:   708,
			BootTime: 1608158324,
			NumProcs: 10,
			Os:       "Windows",
			HostID:   "???",
			TotalMem: 16426635264,
			FreeMem:  8437784576,
			UsedMem:  3781406720,
		},
	}
	
	testProHost5 := model.ProvisionedHost{
		ID:                 "13d4a47c-9ce8-48a2-b6d3-6f5af66b748a",
		SubnetIP:           "10.0.100.10",
		Status:             testStatus,
		Host:               &testLinuxHost1,
		Heartbeat: &model.AgentStatus{
			ClientID: "13d4a47c-9ce8-48a2-b6d3-6f5af66b748a",
			Hostname: "test1.linux.host",
			UpTime:   708,
			BootTime: 1608158324,
			NumProcs: 10,
			Os:       "Linux",
			HostID:   "???",
			Load1:    &testFloat,
			Load5:    &testFloat,
			Load15:   &testFloat,
			TotalMem: 16426635264,
			FreeMem:  8437784576,
			UsedMem:  3781406720,
		},
	}
	testProHost6 := model.ProvisionedHost{
		ID:                 "47f3e7ae-9a87-45c8-900f-8e0043abf4bb",
		SubnetIP:           "10.0.100.15",
		Status:             testStatus,
		Host:               &testLinuxHost2,
		Heartbeat: &model.AgentStatus{
			ClientID: "47f3e7ae-9a87-45c8-900f-8e0043abf4bb",
			Hostname: "test2.linux.host",
			UpTime:   708,
			BootTime: 1608158324,
			NumProcs: 10,
			Os:       "Linux",
			HostID:   "???",
			Load1:    &testFloat,
			Load5:    &testFloat,
			Load15:   &testFloat,
			TotalMem: 16426635264,
			FreeMem:  8437784576,
			UsedMem:  3781406720,
		},
	}

	// Test Provisioned Steps

	testProHost1Steps := []*model.ProvisionedStep{{
		ID:              "a83983d8-625a-49e3-baa3-d7e1748f5f16",
		ProvisionType:   "FileDownload",
		StepNumber:      1,
		ProvisionedHost: &testProHost1,
		Status:          testStatus,
		FileDownload:    &testFileDownload,
		},
		{
		ID:              "6800c01a-066f-471c-8068-25cfc4e1d7bb",
		ProvisionType:   "FileExtract",
		StepNumber:      2,
		ProvisionedHost: &testProHost1,
		Status:          testStatus,
		FileExtract:     &testFileExtract,
	}}

	testProHost2Steps := []*model.ProvisionedStep{{
		ID:              "edb60b33-2fe2-4418-989f-88d7d8c64ca9",
		ProvisionType:   "DNSRecord",
		StepNumber:      1,
		ProvisionedHost: &testProHost2,
		Status:          testStatus,
		DNSRecord:       &testDNSRecord,
	},{
		ID:              "8470d6c7-3ded-4845-ae73-ce511825d740",
		ProvisionType:   "Command",
		StepNumber:      2,
		ProvisionedHost: &testProHost2,
		Status:          testStatus,
		Command:         &testCommand,
	},{
		ID:              "e0fe9c68-e504-4e02-9411-f216c6a6f3ca",
		ProvisionType:   "Script",
		StepNumber:      3,
		ProvisionedHost: &testProHost2,
		Status:          testStatus,
		Script:          &testScript,
	}}

	testProHost3Steps := []*model.ProvisionedStep{{
		ID:              "93e99598-5808-47c8-adf0-5bd90520674f",
		ProvisionType:   "Command",
		StepNumber:      1,
		ProvisionedHost: &model.ProvisionedHost{},
		Status:          testStatus,
		Command:         &testCommand,
	},{
		ID:              "3c2b2a12-8857-43e1-b79c-d863c8a73b3c",
		ProvisionType:   "Script",
		StepNumber:      2,
		ProvisionedHost: &model.ProvisionedHost{},
		Status:          testStatus,
		Script:          &testScript,
	},{
		ID:              "34f62ca7-f79d-4272-8c86-20548ff48af0",
		ProvisionType:   "FileDownload",
		StepNumber:      3,
		ProvisionedHost: &model.ProvisionedHost{},
		Status:          testStatus,
		FileDownload:    &testFileDownload,
	},{
		ID:              "243b5875-d1a7-4236-8168-b3e3b7e86af5",
		ProvisionType:   "FileDelete",
		StepNumber:      4,
		ProvisionedHost: &model.ProvisionedHost{},
		Status:          testStatus,
		FileDelete:      &testFileDelete,
	}}

	testProHost4Steps := []*model.ProvisionedStep{{
		ID:              "f2f7c2f2-d546-4343-8f3e-5955682ca0e4",
		ProvisionType:   "FileDownload",
		StepNumber:      1,
		ProvisionedHost: &testProHost4,
		Status:          testStatus,
		FileDownload:    &testFileDownload,
	},{
		ID:              "da7fc650-d46e-4470-922a-67bba8519ba3",
		ProvisionType:   "FileExtract",
		StepNumber:      2,
		ProvisionedHost: &testProHost4,
		Status:          testStatus,
		FileExtract:     &testFileExtract,
	}}

	testProHost5Steps := []*model.ProvisionedStep{{
		ID:              "a78261db-9699-47fc-818f-047501d774eb",
		ProvisionType:   "DNSRecord",
		StepNumber:      1,
		ProvisionedHost: &testProHost5,
		Status:          testStatus,
		DNSRecord:       &testDNSRecord,
	},{
		ID:              "34b9557e-440d-4669-a412-ca890e321587",
		ProvisionType:   "Command",
		StepNumber:      2,
		ProvisionedHost: &testProHost5,
		Status:          testStatus,
		Command:         &testCommand,
	},{
		ID:              "4b04f5c9-8a78-481d-86ba-6f2e143f1e42",
		ProvisionType:   "Script",
		StepNumber:      3,
		ProvisionedHost: &testProHost5,
		Status:          testStatus,
		Script:          &testScript,
	}}

	testProHost6Steps := []*model.ProvisionedStep{{
		ID:              "bf6335f1-c2bb-4349-8b85-df3dad2a371c",
		ProvisionType:   "Command",
		StepNumber:      1,
		ProvisionedHost: &model.ProvisionedHost{},
		Status:          testStatus,
		Command:         &testCommand,
	},{
		ID:              "396ade50-4ca0-48d0-af4d-9cf17137b323",
		ProvisionType:   "Script",
		StepNumber:      2,
		ProvisionedHost: &model.ProvisionedHost{},
		Status:          testStatus,
		Script:          &testScript,
	},{
		ID:              "5847e03d-614c-4764-8dd3-ffbe9adf9e1c",
		ProvisionType:   "FileDownload",
		StepNumber:      3,
		ProvisionedHost: &model.ProvisionedHost{},
		Status:          testStatus,
		FileDownload:    &testFileDownload,
	},{
		ID:              "3126c79e-dc31-4973-a079-1e9e0708daa2",
		ProvisionType:   "FileDelete",
		StepNumber:      4,
		ProvisionedHost: &model.ProvisionedHost{},
		Status:          testStatus,
		FileDelete:      &testFileDelete,
	}}

	// Test Provisioned Networks
	testProNet1 := model.ProvisionedNetwork{
		ID:               "04f058f8-6d27-42f5-9ef8-2b77115c7b1f",
		Name:             "Test Network 1",
		Cidr:             testNetwork1.Cidr,
		Vars:             []*model.VarsMap{&testVars},
		Tags:             []*model.Tag{},
		ProvisionedHosts: []*model.ProvisionedHost{&testProHost1},
		Status:           testStatus,
		Network:          &testNetwork1,
	}
	testProNet2 := model.ProvisionedNetwork{
		ID:               "00435bab-6e8b-44c2-8bf0-99f62c8ee8bd",
		Name:             "Test Network 2",
		Cidr:             testNetwork2.Cidr,
		Vars:             []*model.VarsMap{},
		Tags:             []*model.Tag{&testTag},
		ProvisionedHosts: []*model.ProvisionedHost{&testProHost2, &testProHost3},
		Status:           testStatus,
		Network:          &testNetwork2,
	}
	testProNet3 := model.ProvisionedNetwork{
		ID:               "d6540443-c510-4cc3-a2a9-fb872e299bd8",
		Name:             "Test Network 3",
		Cidr:             testNetwork3.Cidr,
		Vars:             []*model.VarsMap{},
		Tags:             []*model.Tag{},
		ProvisionedHosts: []*model.ProvisionedHost{&testProHost4, &testProHost5, &testProHost6},
		Status:           testStatus,
		Network:          &testNetwork3,
	}

	// Test Teams
	testTeam1 := model.Team{
		ID:                  "37cc72ca-b6c8-4f17-9b89-3d686f348f19",
		TeamNumber:          0,
		Config:              []*model.ConfigMap{&textConfig},
		Revision:            1,
		Maintainer:          &testUser,
		Tags:                []*model.Tag{},
		ProvisionedNetworks: []*model.ProvisionedNetwork{&testProNet1, &testProNet2, &testProNet3},
	}
	testTeam2 := model.Team{
		ID:                  "62d20b3e-46f2-40ed-a8ae-2c265ef3e621",
		TeamNumber:          1,
		Config:              []*model.ConfigMap{},
		Revision:            1,
		Maintainer:          &testUser,
		Tags:                []*model.Tag{&testTag},
		ProvisionedNetworks: []*model.ProvisionedNetwork{&testProNet1, &testProNet2, &testProNet3},
	}
	testTeam3 := model.Team{
		ID:                  "3e3f397e-1452-45ce-abc2-14f6a9a7792f",
		TeamNumber:          2,
		Config:              []*model.ConfigMap{},
		Revision:            1,
		Maintainer:          &testUser,
		Tags:                []*model.Tag{&testTag},
		ProvisionedNetworks: []*model.ProvisionedNetwork{&testProNet1, &testProNet2, &testProNet3},
	}

	// Test environment
	testComp := model.Competition{
		ID:           "adcd9812-76f5-456b-83b7-ec619f5ad6bb",
		RootPassword: "Password123",
		Config:       []*model.ConfigMap{},
		DNS:          &model.DNS{
			ID:         "5b4f9fb3-880b-412c-9117-8988212056f2",
			Type:       "bind",
			RootDomain: "test.comp",
			DNSServers: []*string{&testDNSServerA, &testDNAServerB},
			NTPServer:  []*string{&testNTPServerA, &testNTPServerB},
		},
	}
	testBuild := model.Build{
		ID:         "dda5e8a1-33ef-4f7f-baa5-56b2535211d5",
		Revision:   1,
		Tags:       []*model.Tag{},
		Config:     []*model.ConfigMap{},
		Maintainer: &testUser,
		Teams:      []*model.Team{&testTeam1, &testTeam2, &testTeam3},
	}
	testEnvironment := model.Environment{
		ID:              "a3f73ee0-da71-4aa6-9280-18ad1a1a8d16",
		CompetitionID:   "adcd9812-76f5-456b-83b7-ec619f5ad6bb",
		Name:            "Test Environment",
		Description:     "A mock environment for the Frontend",
		Builder:         "nativeVMWare",
		TeamCount:       3,
		AdminCIDRs:      []*string{&testCidr},
		ExposedVDIPorts: []*string{&testSSHPort, &testRDPPort},
		Tags:            []*model.Tag{&testTag},
		Config:          []*model.ConfigMap{&textConfig},
		Maintainer:      &testUser,
		Networks:        []*model.Network{&testNetwork1, &testNetwork2, &testNetwork3},
		Hosts:           []*model.Host{&testWindowsHost,&testLinuxHost1,&testLinuxHost2},
		Build:           &testBuild,
		Competition:     &testComp,
	}
	testTeam1.Environment = &testEnvironment
	testTeam2.Environment = &testEnvironment
	testTeam3.Environment = &testEnvironment
	testTeam1.Build = &testBuild
	testTeam2.Build = &testBuild
	testTeam3.Build = &testBuild
	testProNet1.Build = &testBuild
	testProNet2.Build = &testBuild
	testProNet3.Build = &testBuild
	testProHost1.ProvisionedNetwork = &testProNet1
	testProHost1.ProvisionedSteps = testProHost1Steps
	testProHost2.ProvisionedNetwork = &testProNet2
	testProHost2.ProvisionedSteps = testProHost2Steps
	testProHost3.ProvisionedNetwork = &testProNet2
	testProHost3.ProvisionedSteps = testProHost3Steps
	testProHost4.ProvisionedNetwork = &testProNet3
	testProHost4.ProvisionedSteps = testProHost4Steps
	testProHost5.ProvisionedNetwork = &testProNet3
	testProHost5.ProvisionedSteps = testProHost5Steps
	testProHost6.ProvisionedNetwork = &testProNet3
	testProHost6.ProvisionedSteps = testProHost6Steps
	allTestSteps := []*model.ProvisionedStep{}
	allTestSteps = append(allTestSteps,testProHost1Steps...)
	allTestSteps = append(allTestSteps,testProHost2Steps...)
	allTestSteps = append(allTestSteps,testProHost3Steps...)
	allTestSteps = append(allTestSteps,testProHost4Steps...)
	allTestSteps = append(allTestSteps,testProHost5Steps...)
	allTestSteps = append(allTestSteps,testProHost6Steps...)
	return 	&testEnvironment, 
			[]*model.ProvisionedNetwork{&testProNet1,&testProNet2,&testProNet3}, 
			[]*model.ProvisionedHost{&testProHost1,&testProHost2,&testProHost3,&testProHost4,&testProHost5,&testProHost6},
			allTestSteps
}