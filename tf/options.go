package tf

// func NewEmptyOptions() (*options.TerragruntOptions, error) {
// 	binname := "terraform"
// 	if runtime.GOOS == "windows" {
// 		binname = "terraform.exe"
// 	}
// 	tfexepath, err := exec.LookPath(binname)
// 	if err != nil {
// 		return nil, err
// 	}

// 	m := make(map[string]string)
// 	for _, e := range os.Environ() {
// 		if i := strings.Index(e, "="); i >= 0 {
// 			m[e[:i]] = e[i+1:]
// 		}
// 	}

// 	return &options.TerragruntOptions{
// 		TerraformPath:    tfexepath,
// 		Env:              m,
// 		NonInteractive:   true,
// 		AutoInit:         true,
// 		AutoRetry:        true,
// 		MaxRetryAttempts: 5,
// 		Sleep:            10 * time.Second,
// 	}, err
// }
