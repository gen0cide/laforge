package core

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/juju/utils/filepath"
	"github.com/masterzen/winrm"
	"github.com/masterzen/winrm/soap"
	"github.com/pkg/errors"
)

// WinRMClient is a type to connection to Windows hosts remotely over the WinRM protocol
type WinRMClient struct {
	Config *WinRMAuthConfig
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// Kind implements the Sheller interface
func (w *WinRMClient) Kind() string {
	return "winrm"
}

// SetIO implements the Sheller interface
func (w *WinRMClient) SetIO(stdout io.Writer, stderr io.Writer, stdin io.Reader) error {
	w.Stdin = stdin
	w.Stdout = stdout
	w.Stderr = stderr
	return nil
}

// SetConfig implements the Sheller interface
func (w *WinRMClient) SetConfig(c *WinRMAuthConfig) error {
	if c == nil {
		return errors.New("nil auth config provised")
	}
	w.Config = c
	return nil
}

// LaunchInteractiveShell implements the Sheller interface
func (w *WinRMClient) LaunchInteractiveShell() error {
	endpoint := winrm.NewEndpoint(
		w.Config.RemoteAddr,
		w.Config.Port,
		w.Config.HTTPS,
		w.Config.SkipVerify,
		[]byte{},
		[]byte{},
		[]byte{},
		0,
	)

	if w.Stderr == nil {
		w.Stderr = os.Stderr
	}

	if w.Stdin == nil {
		w.Stdin = os.Stdin
	}

	if w.Stdout == nil {
		w.Stdout = os.Stdout
	}

	params := winrm.DefaultParameters
	params.Timeout = "PT24H"
	client, err := winrm.NewClientWithParameters(endpoint, w.Config.User, w.Config.Password, params)
	if err != nil {
		return errors.WithMessage(err, "could not create winrm client")
	}

	shell, err := client.CreateShell()
	if err != nil {
		panic(err)
	}
	var cmd *winrm.Command
	cmd, err = shell.Execute("powershell -NoProfile -ExecutionPolicy Bypass")
	if err != nil {
		panic(err)
	}

	go io.Copy(cmd.Stdin, os.Stdin)
	go io.Copy(os.Stdout, cmd.Stdout)
	go io.Copy(os.Stderr, cmd.Stderr)

	cmd.Wait()
	shell.Close()

	return nil
}

// ExecuteNonInteractive allows you to execute commands in a non-interactive session (note: standard command shell, not powershell)
func (w *WinRMClient) ExecuteNonInteractive(cmd *RemoteCommand) error {
	endpoint := winrm.NewEndpoint(
		w.Config.RemoteAddr,
		w.Config.Port,
		w.Config.HTTPS,
		w.Config.SkipVerify,
		nil,
		nil,
		nil,
		60,
	)

	transporter := &AdvancedTransporter{
		auth: w.Config,
	}

	fmt.Printf("WinRM ExecuteNonInteractive: cmd: \n%+v\n", cmd)
	fmt.Printf("WinRM ExecuteNonInteractive: this: \n%+v\n", *w)

	fmt.Printf("Set up WinRM client: start.\n")

	params := winrm.DefaultParameters
	params.Timeout = "PT12M"
	params.TransportDecorator = func() winrm.Transporter { return transporter }
	client, err := winrm.NewClientWithParameters(endpoint, w.Config.User, w.Config.Password, params)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Set up WinRM client: complete.\n")

	fmt.Printf("WinRM CreateShell: start\n")
	shell, err := client.CreateShell()
	if err != nil {
		log.Printf("[ERROR] error creating shell: %s", err)
		return err
	}
	fmt.Printf("WinRM CreateShell: complete\n")

	err = shell.Close()
	if err != nil {
		log.Printf("[ERROR] error closing shell: %s", err)
		return err
	}

	winfp, err := filepath.NewRenderer("windows")
	if err != nil {
		panic(err)
	}

	if winfp.Ext(cmd.Command) == `.ps1` && !strings.Contains(cmd.Command, " ") {
		cmdstrbuf := new(bytes.Buffer)
		err = elevatedCommandTemplate.Execute(cmdstrbuf, struct{ Path string }{
			Path: cmd.Command,
		})
		if err != nil {
			return err
		}

		escp := new(bytes.Buffer)
		err = xml.EscapeText(escp, cmdstrbuf.Bytes())
		if err != nil {
			return err
		}

		eo := elevatedOptions{
			User:              w.Config.User,
			Password:          w.Config.Password,
			TaskName:          winfp.Base(cmd.Command),
			LogFile:           fmt.Sprintf("%s.log", cmd.Command),
			TaskDescription:   "running laforge command",
			XMLEscapedCommand: escp.String(),
		}

		outbuf := new(bytes.Buffer)
		err = elevatedTemplate.Execute(outbuf, eo)
		if err != nil {
			return err
		}

		encoded := Powershell(outbuf.String())
		cmd.Command = fmt.Sprintf("powershell -NoProfile -ExecutionPolicy Bypass -EncodedCommand %s", encoded)
	}

	fmt.Printf("WinRM client.Run: start\n")
	fmt.Printf("WinRM client.Run: Command: %s\n", cmd.Command)
	status, err := client.Run(cmd.Command, cmd.Stdout, cmd.Stderr)
	fmt.Printf("WinRM client.Run: complete. Status: %d\n", status)
	if err != nil {
		fmt.Printf("WinRM client.Run: complete. Error: %s\n", err.Error())
	}

	cmd.SetExitStatus(status, err)

	// return nil

	// go io.Copy(wcmd.Stdin, stdin)
	// go io.Copy(cmd.Stdout, wcmd.Stdout)
	// go io.Copy(cmd.Stderr, wcmd.Stderr)

	// err = cmd.Wait()
	// if err != nil {
	// 	panic(err)
	// }

	// err = shell.Close()
	// if err != nil {
	// 	panic(err)
	// }

	// status, err := client.Run(cmd.Command, cmd.Stdout, cmd.Stderr)

	// wcmd, err = shell.Execute(cmd.Command)
	// if err != nil {
	// 	panic(err)
	// }

	// if cmd.Stdin != nil {
	// 	go io.Copy(wcmd.Stdin, cmd.Stdin)
	// }

	// go io.Copy(cmd.Stdout, wcmd.Stdout)
	// go io.Copy(cmd.Stderr, wcmd.Stderr)

	// wcmd.Wait()
	// cmderr := wcmd.Close()
	// exitStatus := wcmd.ExitCode()

	// cmd.SetExitStatus(status, err)
	// if err != nil {
	// 	panic(err)
	// }
	return nil
}

type elevatedOptions struct {
	User              string
	Password          string
	TaskName          string
	TaskDescription   string
	LogFile           string
	XMLEscapedCommand string
}

// Powershell wraps a PowerShell script
// and prepares it for execution by the winrm client
func Powershell(psCmd string) string {
	// 2 byte chars to make PowerShell happy
	wideCmd := ""
	for _, b := range []byte(psCmd) {
		wideCmd += string(b) + "\x00"
	}

	// Base64 encode the command
	input := []uint8(wideCmd)
	encodedCmd := base64.StdEncoding.EncodeToString(input)

	// Create the powershell.exe command line to execute the script
	return fmt.Sprintf("%s", encodedCmd)
}

var elevatedCommandTemplate = template.Must(template.New("ElevatedCommandRunner").Parse(`powershell -noprofile -executionpolicy bypass "& { if (Test-Path variable:global:ProgressPreference){set-variable -name variable:global:ProgressPreference -value 'SilentlyContinue'}; &'{{.Path}}'; exit $LastExitCode }"`))

var elevatedTemplate = template.Must(template.New("ElevatedCommand").Parse(`
$name = "{{.TaskName}}"
$log = [System.Environment]::ExpandEnvironmentVariables("{{.LogFile}}")
$s = New-Object -ComObject "Schedule.Service"
$s.Connect()
$t = $s.NewTask($null)
$t.XmlText = @'
<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.2" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
	<Description>{{.TaskDescription}}</Description>
  </RegistrationInfo>
  <Principals>
    <Principal id="Author">
      <UserId>{{.User}}</UserId>
      <LogonType>Password</LogonType>
      <RunLevel>HighestAvailable</RunLevel>
    </Principal>
  </Principals>
  <Settings>
    <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
    <DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
    <StopIfGoingOnBatteries>false</StopIfGoingOnBatteries>
    <AllowHardTerminate>true</AllowHardTerminate>
    <StartWhenAvailable>false</StartWhenAvailable>
    <RunOnlyIfNetworkAvailable>false</RunOnlyIfNetworkAvailable>
    <IdleSettings>
      <StopOnIdleEnd>false</StopOnIdleEnd>
      <RestartOnIdle>false</RestartOnIdle>
    </IdleSettings>
    <AllowStartOnDemand>true</AllowStartOnDemand>
    <Enabled>true</Enabled>
    <Hidden>false</Hidden>
    <RunOnlyIfIdle>false</RunOnlyIfIdle>
    <WakeToRun>false</WakeToRun>
    <ExecutionTimeLimit>PT24H</ExecutionTimeLimit>
    <Priority>4</Priority>
  </Settings>
  <Actions Context="Author">
    <Exec>
      <Command>cmd</Command>
      <Arguments>/c {{.XMLEscapedCommand}}</Arguments>
    </Exec>
  </Actions>
</Task>
'@
if (Test-Path variable:global:ProgressPreference){$ProgressPreference="SilentlyContinue"}
$f = $s.GetFolder("\")
$f.RegisterTaskDefinition($name, $t, 6, "{{.User}}", "{{.Password}}", 1, $null) | Out-Null
$t = $f.GetTask("\$name")
$t.Run($null) | Out-Null
$timeout = 10
$sec = 0
while ((!($t.state -eq 4)) -and ($sec -lt $timeout)) {
  Start-Sleep -s 1
  $sec++
}
$line = 0
do {
  Start-Sleep -m 100
  if (Test-Path $log) {
    Get-Content $log | select -skip $line | ForEach {
      $line += 1
      Write-Output "$_"
    }
  }
} while (!($t.state -eq 3))
$result = $t.LastTaskResult
if (Test-Path $log) {
    Remove-Item $log -Force -ErrorAction SilentlyContinue | Out-Null
}
[System.Runtime.Interopservices.Marshal]::ReleaseComObject($s) | Out-Null
exit $result`))

type AdvancedTransporter struct {
	auth      *WinRMAuthConfig
	transport http.RoundTripper
	endpoint  *winrm.Endpoint
	client    *http.Client
}

func (a *AdvancedTransporter) Transport(endpoint *winrm.Endpoint) error {
	t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: false,
		}).DialContext,
		MaxIdleConns:          1,
		IdleConnTimeout:       45 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: 60 * time.Second,
	}

	a.transport = t
	a.endpoint = endpoint

	client := &http.Client{
		Transport: a.transport,
	}

	a.client = client
	return nil
}

func (a *AdvancedTransporter) URL() string {
	var scheme string
	if a.endpoint.HTTPS {
		scheme = "https"
	} else {
		scheme = "http"
	}

	return fmt.Sprintf("%s://%s:%d/wsman", scheme, a.endpoint.Host, a.endpoint.Port)
}

func (a *AdvancedTransporter) Post(client *winrm.Client, request *soap.SoapMessage) (string, error) {
	req, err := http.NewRequest("POST", a.URL(), strings.NewReader(request.String()))
	if err != nil {
		return "", errors.Wrap(err, "impossible to create http request")
	}

	req.Header.Set("Content-Type", "application/soap+xml;charset=UTF-8")
	req.SetBasicAuth(a.auth.User, a.auth.Password)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req = req.WithContext(ctx)

	resp, err := a.client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "http request did not finish successfully")
	}

	defer resp.Body.Close()

	if !strings.Contains(resp.Header.Get("Content-Type"), "application/soap+xml") {
		return "", errors.Errorf("invalid content type returned: %s (expected application/soap+xml)", resp.Header.Get("Content-Type"))
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "error reading http response body")
	}

	if resp.StatusCode != 200 {
		return string(data), errors.Errorf("non-200 response from the WinRM server: %s (%d)", resp.Status, resp.StatusCode)
	}

	return string(data), nil
}
