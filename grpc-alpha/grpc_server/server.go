package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/frybin/laforge/grpc-alpha/laforge_proto_agent"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"gorm.io/gorm"
)

const (
	port     = ":50051"
	certFile = "server.crt"
	keyFile  = "server.key"
	webPort  = ":5000"

	TaskFailed = "Failed"
	TaskRunning = "Running"
	TaskSucceeded = "Completed"
)

var (
	db *gorm.DB
)

type server struct {
	pb.UnimplementedLaforgeServer
}

/* TEST MESSAGES */

//Ping Info
func (s *server) GetPing(ctx context.Context, in *pb.PingRequest) (*pb.PingReply, error) {
	log.Printf("Received: %v | ID: %v", in.GetName(), in.GetId())
	return &pb.PingReply{Name: "Hello " + in.GetName(), Id: in.GetId()}, nil
}

//HostTest Info
func (s *server) GetHostTest(ctx context.Context, in *pb.HostTestRequest) (*pb.HostTestReply, error) {
	log.Printf("Got Host: %v | ID: %v | IP: %s | OS: %s", in.GetName(), in.GetId(), in.GetIp(), in.GetOs())
	return &pb.HostTestReply{Name: in.GetName(), Id: in.GetId(), Ip: in.GetIp(), Os: in.GetOs()}, nil
}

//ByteCountIEC Converts Bytes to Higher Order
func ByteCountIEC(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

//GetHeartBeat Info
func (s *server) GetHeartBeat(ctx context.Context, in *pb.HeartbeatRequest) (*pb.HeartbeatReply, error) {
	message := fmt.Sprintf("Recived ID: %v | Hostname: %v | Uptime: %v | Boot Time: %v| Number of Running Processes: %v| OS Arch: %v| Host ID: %v| Load1: %v| Load5: %v| Load15: %v| Total Memory: %v| Avalible Memory: %v| Used Memory: %v", in.GetClientId(), in.GetHostname(), in.GetUptime(), in.GetBoottime(), in.GetNumprocs(), in.GetOs(), in.GetHostid(), in.GetLoad1(), in.GetLoad5(), in.GetLoad15(), ByteCountIEC(in.GetTotalmem()), ByteCountIEC(in.GetFreemem()), ByteCountIEC(in.GetUsedmem()))
	log.Printf(message)
	avalibleTasks := false
	tasks := make([]Task, 0)
	db.Find(&tasks, map[string]interface{}{"client_id": in.GetClientId(), "completed": false})
	if len(tasks) > 0 {
		avalibleTasks = true
	}
	return &pb.HeartbeatReply{Status: message, AvalibleTasks: avalibleTasks}, nil
}

//GetTask Info
func (s *server) GetTask(ctx context.Context, in *pb.TaskRequest) (*pb.TaskReply, error) {
	clientID := in.ClientId
	tasks := make([]Task, 0)
	db.Order("task_id asc").Find(&tasks, map[string]interface{}{"client_id": clientID, "completed": false})
	task := tasks[0]

	if len(tasks) > 0 {
		task := tasks[0]
		return &pb.TaskReply{Id: task.TaskID, Command: pb.TaskReply_Command(task.CommandID), Args: task.Args}, nil
	}
	return &pb.TaskReply{Id: 0, Command: pb.TaskReply_DEFAULT}, nil
}

// InformTaskStatus
func (s *server) InformTaskStatus(ctx context.Context, in *pb.TaskStatusRequest) (*pb.TaskStatusReply, error) {
	clientID := in.ClientId
	tasks := make([]Task, 0)
	db.Order("task_id asc").Find(&tasks, map[string]interface{}{"client_id": clientID, "completed": false})
	task := tasks[0]

	switch status := in.status; status {
		case TaskRunning:
			task.Status = TaskRunning
		case TaskFailed:
			task.Status = TaskFailed
		case TaskSucceeded:
			task.Status = TaskSucceeded
			task.Completed = true
	}

	db.Save(&task)
	return &pb.TaskStatusReply{Status: task.Status}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	db = OpenDB()

	log.Printf("Starting API Server on port " + webPort)
	web := gin.Default()
	web.GET("/download/:file_id", TempURLHandler)

	apiGroup := web.Group("/api")
	{
		add := apiGroup.Group("add")
		add.POST("task", TaskAdder)
		add.POST("file", FileAdder)
	}

	go web.Run(webPort)

	// secure server
	creds, _ := credentials.NewServerTLSFromFile(certFile, keyFile)
	s := grpc.NewServer(grpc.Creds(creds))

	//insecure server
	// s := grpc.NewServer()

	log.Printf("Starting Laforge Server on port " + port)

	pb.RegisterLaforgeServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
