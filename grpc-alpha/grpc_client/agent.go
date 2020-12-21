package main

//go:generate fileb0x assets.toml
import (
	"context"
	"crypto/md5"
	"crypto/x509"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gen0cide/laforge/grpc-alpha/grpc_client/static"
	pb "github.com/gen0cide/laforge/grpc-alpha/laforge_proto_agent"
	"github.com/kardianos/service"
	"github.com/mholt/archiver"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	TaskFailed = "Failed"
	TaskRunning = "Running"
	TaskSucceeded = "Completed"
)

var (
	logger service.Logger
	address          = "localhost:50051"
	defaultName      = "Laforge Agent"
	certFile         = "service.pem"
	heartbeatSeconds = 1
	clientID         = "1"
)

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

// Start What is Run when the executable is started up.
func (p *program) Start(s service.Service) error {
	p.exit = make(chan struct{})

	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

// ExecuteCommand Runs the Command that is inputted and either returns the error or output
func ExecuteCommand(command string, args ...string) string {
	output, err := exec.Command(command, args...).Output()
	if err != nil {
		return err.Error()
	}
	return string(output)
}

// DeleteObject Deletes the Object that is inputted and either returns the error or nothing
func DeleteObject(file string) error {
	err := os.RemoveAll(file)
	if err != nil {
		return err
	}
	return nil
}

// Reboot Reboots Host Operating System
func Reboot() {
	RebootSystem()
}

// ExtractArchive will extract archive to foler path.
func ExtractArchive(filepath string, folderpath string) error {
	err := archiver.Unarchive(filepath, folderpath)
	return err
}

// CreateUser will create a new user.
func CreateUser(username string, password string) error {
	return CreateSystemUser(username, password)
}

// ChangeUserPassword will change the users password
func ChangeUserPassword(username string, password string) error {
	return ChangeSystemUserPassword(username, password)
}

// AddUserGroup will extract archive to foler path.
func AddUserGroup(groupname string, username string) error {
	return AddSystemUserGroup(groupname, username)
}

// DownloadFile will download a url to a local file.
func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}


// ValidateMD5Hash Validates the MD5 Hash of a file with the provided MD5 Hash
func ValidateMD5Hash(filepath string, md5hash string) error {
	var calculatedMD5Hash string

	// Open the file 
	file, err := os.Open(filepath)

	// Can't open the file, assuming false
	if err != nil {
		return err
	}

	// Close the file when we're done
	defer file.Close()

	// Open a new hash interface
	hash := md5.New()

	// Hash the file
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	byteHash := hash.Sum(nil)[:16]

	// Convert bytes to string
	calculatedMD5Hash = hex.EncodeToString(byteHash)

	if calculatedMD5Hash == md5hash {
		return errors.New("MD5 hashes do not match")
	} else {
		return nil
	}
}

// RequestTask Function Requests task from the GRPC server to be run on the client
func RequestTask(c pb.LaforgeClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	request := &pb.TaskRequest{ClientId: clientID}
	r, err := c.GetTask(ctx, request)

	taskRequest := &pb.TaskStatusRequest{ClientId: clientID, Status: TaskRunning}
	c.InformTaskStatus(ctx, taskRequest)

	if err != nil {
		logger.Errorf("Error: %v", err)
	} else {
		switch r.Command {
		case pb.TaskReply_EXECUTE:
			taskArgs := strings.Split(r.Args, " ")
			command := taskArgs[0]
			args := taskArgs[1:]
			output := ExecuteCommand(command, args...)
			logger.Infof("Command Output: %s", output)
			RequestTaskStatusRequest(nil, clientID, c)
		case pb.TaskReply_DOWNLOAD:
			taskArgs := strings.Split(r.Args, ",")
			filepath := taskArgs[0]
			url := taskArgs[1]
			taskerr := DownloadFile(filepath, url)
			RequestTaskStatusRequest(taskerr, clientID, c)
		case pb.TaskReply_EXTRACT:
			taskArgs := strings.Split(r.Args, ",")
			filepath := taskArgs[0]
			folder := taskArgs[1]
			taskerr := ExtractArchive(filepath, folder)
			RequestTaskStatusRequest(taskerr, clientID, c)
		case pb.TaskReply_DELETE:
			taskerr := DeleteObject(r.Args)
			RequestTaskStatusRequest(taskerr, clientID, c)
		case pb.TaskReply_REBOOT:
			Reboot()
			taskRequest := &pb.TaskStatusRequest{ClientId: clientID, Status: TaskSucceeded}
			c.InformTaskStatus(ctx, taskRequest)
		case pb.TaskReply_CREATEUSER:
			taskArgs := strings.Split(r.Args, ",")
			username := taskArgs[0]
			password := taskArgs[1]
			taskerr := CreateUser(username, password)
			RequestTaskStatusRequest(taskerr, clientID, c)
		case pb.TaskReply_ADDTOGROUP:
			taskArgs := strings.Split(r.Args, ",")
			group := taskArgs[0]
			username := taskArgs[1]
			taskerr := AddUserGroup(group, username)
			RequestTaskStatusRequest(taskerr, clientID, c)
		case pb.TaskReply_CREATEUSERPASS:
			taskArgs := strings.Split(r.Args, ",")
			username := taskArgs[0]
			password := taskArgs[1]
			taskerr := ChangeUserPassword(username, password)
			RequestTaskStatusRequest(taskerr, clientID, c)
		case pb.TaskReply_VALIDATE:
			taskArgs := strings.Split(r.Args, ",")
			filepath := taskArgs[0]
			md5hash := taskArgs[1]
			taskerr := ValidateMD5Hash(filepath, md5hash)
			RequestTaskStatusRequest(taskerr, clientID, c)
		default:
			logger.Infof("Response Message: %v", r)
		    RequestTaskStatusRequest(nil, clientID, c)
		}
	}
}


// RequestTaskStatusRequest Tell the server the status of a completed task
func RequestTaskStatusRequest(taskerr error, clientID string, c pb.LaforgeClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if taskerr != nil {
		logger.Errorf("Error: %v", taskerr)
		taskRequest := &pb.TaskStatusRequest{ClientId: clientID, Status: TaskFailed}
		c.InformTaskStatus(ctx, taskRequest)
	} else {
		taskRequest := &pb.TaskStatusRequest{ClientId: clientID, Status: TaskSucceeded}
		c.InformTaskStatus(ctx, taskRequest)
	}
}

// SendHeartBeat Send the GRPC server a Heartbeat with specified parameters
func SendHeartBeat(c pb.LaforgeClient, taskChannel chan pb.HeartbeatReply) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	request := &pb.HeartbeatRequest{ClientId: clientID}
	hostInfo, hostErr := host.Info()
	if hostErr == nil {
		(*request).Hostname = hostInfo.Hostname
		(*request).Uptime = hostInfo.Uptime
		(*request).Boottime = hostInfo.BootTime
		(*request).Numprocs = hostInfo.Procs
		(*request).Os = hostInfo.OS
		(*request).Hostid = hostInfo.HostID
	}
	mem, memErr := mem.VirtualMemory()
	if memErr == nil {
		(*request).Totalmem = mem.Total
		(*request).Freemem = mem.Free
		(*request).Usedmem = mem.Used
	}
	load, loadErr := load.Avg()
	if loadErr == nil {
		(*request).Load1 = load.Load1
		(*request).Load5 = load.Load5
		(*request).Load15 = load.Load15
	}
	r, err := c.GetHeartBeat(ctx, request)
	if err != nil {
		logger.Errorf("Error: %v", err)
	}else{
		taskChannel <- *r
	}
	
}

// StartTaskRunner Gets a Heartbeat reply from the task channel, and if there are avalible tasks it will request them
func StartTaskRunner(c pb.LaforgeClient, taskChannel chan pb.HeartbeatReply) {
	r := <-taskChannel
	logger.Infof("Response Message: %s", r.GetStatus())
	logger.Infof("Avalible Tasks: %s", r.GetAvalibleTasks())

	if r.GetAvalibleTasks() {
		RequestTask(c)
	}
}

// genSendHeartBeat A goroutine that is called, which periodically send a heartbeat to the GRPC Server
func genSendHeartBeat(p *program, c pb.LaforgeClient, taskChannel chan pb.HeartbeatReply, wg sync.WaitGroup) {
	ticker := time.NewTicker(time.Duration(heartbeatSeconds) * time.Second)
	defer wg.Done()
	for {
		select {
		case <-ticker.C:
			SendHeartBeat(c, taskChannel)
		case <-p.exit:
			ticker.Stop()
		}
	}
}

// genStartTaskRunner A goroutine that is called, which checks responses from GRPC server for avalible tasks
func genStartTaskRunner(p *program, c pb.LaforgeClient, taskChannel chan pb.HeartbeatReply, wg sync.WaitGroup) {
	ticker := time.NewTicker(time.Duration(heartbeatSeconds) * time.Second)
	defer wg.Done()
	for {
		select {
		case <-ticker.C:
			StartTaskRunner(c, taskChannel)
		case <-p.exit:
			ticker.Stop()
		}
	}
}

// run Function that is called when the program starts and run all the Go Routines
func (p *program) run() error {
	logger.Infof("I'm running %v.", service.Platform())
	var wg sync.WaitGroup

	// TLS Cert for verifying GRPC Server
	certPem,certerr := static.ReadFile(certFile)
	if certerr != nil {
        fmt.Println("File reading error", certerr)
        return nil
	}
	
	// Starts GRPC Connection with cert included in the binary
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(certPem)
	creds := credentials.NewClientTLSFromCert(certPool, "")
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(creds))

	if err != nil {
		logger.Errorf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewLaforgeClient(conn)
	
	// START VARS
	taskChannel := make(chan pb.HeartbeatReply)
	wg.Add(2)
	go genSendHeartBeat(p, c, taskChannel, wg)
	go genStartTaskRunner(p, c, taskChannel, wg)

	wg.Wait()
	return nil
}

// Stop Called when the Agent is closed
func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Info("I'm Stopping!")
	close(p.exit)
	return nil
}

// Service setup.
//   Define service config.
//   Create the service.
//   Setup the logger.
//   Handle service controls (optional).
//   Run the service.
func main() {
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	options := make(service.KeyValue)
	options["Restart"] = "on-success"
	options["SuccessExitStatus"] = "1 2 8 SIGKILL"
	svcConfig := &service.Config{
		Name:        "laforge-agent",
		DisplayName: "Laforge Client Agent",
		Description: "Laforge Client Agent",
		Dependencies: []string{
			"Requires=network.target",
			"After=network-online.target"},
		Option: options,
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		logger.Error(err)
	}
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		logger.Error(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				logger.Info(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			logger.Infof("Valid actions: %q\n", service.ControlAction)
			logger.Error(err)
		}
		return
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
