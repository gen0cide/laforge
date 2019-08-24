package core

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/gen0cide/laforge/core/cli"
	"github.com/juju/utils/filepath"
	"github.com/masterzen/winrm"
	"github.com/masterzen/winrm/soap"
	"github.com/pkg/errors"
)

var DefaultWinRMTimeout = 60

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

// TestConnection makes a basic connection to the WinRM server to validate it is working
func (w *WinRMClient) TestConnection() bool {
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

	client, err := winrm.NewClient(endpoint, w.Config.User, w.Config.Password)
	if err != nil {
		return false
	}

	shell, err := client.CreateShell()
	if err != nil {
		return false
	}

	shell.Close()
	return true
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
		return errors.WithMessage(err, "could not create WinRM shell connection")
	}
	var cmd *winrm.Command
	cmd, err = shell.Execute("powershell -NoProfile -ExecutionPolicy Bypass")
	if err != nil {
		return errors.WithMessage(err, "could not execute PowerShell for interactive session")
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
	timeout := DefaultWinRMTimeout
	if cmd.Timeout > 0 {
		timeout = cmd.Timeout
	}
	endpoint := winrm.NewEndpoint(
		w.Config.RemoteAddr,
		w.Config.Port,
		w.Config.HTTPS,
		w.Config.SkipVerify,
		nil,
		nil,
		nil,
		(time.Duration(timeout) * time.Second),
	)

	transporter := &AdvancedTransporter{
		auth:    w.Config,
		Timeout: timeout,
	}

	params := winrm.DefaultParameters
	params.Timeout = fmt.Sprintf("PT%dM", (timeout / 60))
	params.TransportDecorator = func() winrm.Transporter { return transporter }
	client, err := winrm.NewClientWithParameters(endpoint, w.Config.User, w.Config.Password, params)
	if err != nil {
		panic(err)
	}

	shell, err := client.CreateShell()
	if err != nil {
		cli.Logger.Errorf("Failed to create a WinRM shell successfully: %v", err)
		return err
	}

	err = shell.Close()
	if err != nil {
		cli.Logger.Errorf("Failed to close the WinRM shell successfully: %v", err)
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

	cli.Logger.Debug("Executing WinRM command...")
	status, err := client.Run(cmd.Command, cmd.Stdout, cmd.Stderr)
	cli.Logger.Debugf("Completed WinRM execution with exit code %d (errored=%v)", status, (err != nil))
	cmd.SetExitStatus(status, err)

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
	Timeout   int
}

func (a *AdvancedTransporter) Transport(endpoint *winrm.Endpoint) error {
	timeout := a.Timeout
	if a.Timeout == 0 {
		timeout = 60
	}

	t := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(timeout+5) * time.Second,
			KeepAlive: time.Duration(timeout+5) * time.Second,
			DualStack: false,
		}).DialContext,
		MaxIdleConns:          1,
		IdleConnTimeout:       time.Duration(timeout+15) * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		ResponseHeaderTimeout: time.Duration(timeout+10) * time.Second,
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
	timeout := a.Timeout
	if timeout == 0 {
		timeout = 60
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
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
