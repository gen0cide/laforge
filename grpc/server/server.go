package server

//go:generate fileb0x assets.toml
import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/agentstatus"
	"github.com/gen0cide/laforge/ent/provisionedhost"
	pb "github.com/gen0cide/laforge/grpc/proto"
	"github.com/gen0cide/laforge/grpc/server/static"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	Port     = ":50051"
	CertFile = "service.pem"
	KeyFile  = "service.key"
	webPort  = ":5000"
)

const (
	TaskFailed    = "Failed"
	TaskRunning   = "Running"
	TaskSucceeded = "Completed"
)

type Server struct {
	Client *ent.Client
	pb.UnimplementedLaforgeServer
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

//GetHeartBeat Recives Heartbeat from client and sends back a reply
func (s *Server) GetHeartBeat(ctx context.Context, in *pb.HeartbeatRequest) (*pb.HeartbeatReply, error) {
	message := fmt.Sprintf("Recived ID: %v | Hostname: %v | Uptime: %v | Boot Time: %v| Number of Running Processes: %v| OS Arch: %v| Host ID: %v| Load1: %v| Load5: %v| Load15: %v| Total Memory: %v| Avalible Memory: %v| Used Memory: %v", in.GetClientId(), in.GetHostname(), in.GetUptime(), in.GetBoottime(), in.GetNumprocs(), in.GetOs(), in.GetHostid(), in.GetLoad1(), in.GetLoad5(), in.GetLoad15(), ByteCountIEC(in.GetTotalmem()), ByteCountIEC(in.GetFreemem()), ByteCountIEC(in.GetUsedmem()))
	avalibleTasks := false
	uuid, err := strconv.Atoi(in.GetClientId())

	if err != nil {
		return nil, fmt.Errorf("failed casting UUID to int: %v", err)
	}

	statusExist, err := s.Client.AgentStatus.Query().Where(agentstatus.ClientIDEQ(in.GetClientId())).Exist(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed querying Agent Status: %v", err)
	}

	// TODO; Will Wanna change this relationship to be 1-M between Host and AgentStatus
	if statusExist {
		agentStatus, err := s.Client.AgentStatus.Query().Where(agentstatus.ClientIDEQ(in.GetClientId())).Only(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed querying Agent Status: %v", err)
		}
		agentStatus, err = agentStatus.Update().
			SetUpTime(int64(in.GetUptime())).
			SetBootTime(int64(in.GetBoottime())).
			SetNumProcs(int64(in.GetNumprocs())).
			SetLoad1(in.GetLoad1()).
			SetLoad5(in.GetLoad5()).
			SetLoad15(in.GetLoad15()).
			SetTotalMem(int64(in.GetTotalmem())).
			SetFreeMem(int64(in.GetFreemem())).
			SetUsedMem(int64(in.GetUsedmem())).
			SetTimestamp(in.GetTimestamp().Seconds).
			Save(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed updating Agent Status: %v", err)
		}
	} else {
		ph, err := s.Client.ProvisionedHost.Query().Where(provisionedhost.IDEQ(uuid)).Only(ctx)

		_, err = s.Client.AgentStatus.
			Create().
			SetClientID(in.GetClientId()).
			SetHostname(in.GetHostname()).
			SetUpTime(int64(in.GetUptime())).
			SetBootTime(int64(in.GetBoottime())).
			SetNumProcs(int64(in.GetNumprocs())).
			SetOs(in.GetOs()).
			SetHostID(in.GetHostid()).
			SetLoad1(in.GetLoad1()).
			SetLoad5(in.GetLoad5()).
			SetLoad15(in.GetLoad15()).
			SetTotalMem(int64(in.GetTotalmem())).
			SetFreeMem(int64(in.GetFreemem())).
			SetUsedMem(int64(in.GetUsedmem())).
			SetTimestamp(in.GetTimestamp().Seconds).
			AddAgentStatusToProvisionedHost(ph).
			Save(ctx)

		if err != nil {
			return nil, fmt.Errorf("failed Creating Agent Status: %v", err)
		}
	}

	// TODO: Implement This with ENT
	// tasks := make([]Task, 0)
	// db.Find(&tasks, map[string]interface{}{"client_id": in.GetClientId(), "completed": false})
	// if len(tasks) > 0 {
	// 	avalibleTasks = true
	// }
	return &pb.HeartbeatReply{Status: message, AvalibleTasks: avalibleTasks}, nil
}

//GetTask Gets a task that needs to be run on the client and sends it over
func (s *Server) GetTask(ctx context.Context, in *pb.TaskRequest) (*pb.TaskReply, error) {
	// TODO: Implement This with ENT
	// clientID := in.ClientId
	// db.Order("task_id asc").Find(&tasks, map[string]interface{}{"client_id": clientID, "completed": false})
	// tasks := make([]Task, 0)
	// if len(tasks) > 0 {
	// 	task := tasks[0]
	// 	return &pb.TaskReply{Id: task.TaskID, Command: pb.TaskReply_Command(task.CommandID), Args: task.Args}, nil
	// }

	return &pb.TaskReply{Id: 0, Command: pb.TaskReply_DEFAULT}, nil
}

// InformTaskStatus Updates the status of a Task on a client from the response of the client
func (s *Server) InformTaskStatus(ctx context.Context, in *pb.TaskStatusRequest) (*pb.TaskStatusReply, error) {
	// TODO: Implement This with ENT
	// clientID := in.ClientId
	// tasks := make([]Task, 0)
	// db.Order("task_id asc").Find(&tasks, map[string]interface{}{"client_id": clientID, "completed": false})
	// task := tasks[0]

	// switch in.Status {
	// 	case TaskRunning:
	// 		task.Status = TaskRunning
	// 	case TaskFailed:
	// 		task.Status = TaskFailed
	// 	case TaskSucceeded:
	// 		task.Status = TaskSucceeded
	// 		task.Completed = true
	// }

	// db.Save(&task)
	return &pb.TaskStatusReply{Status: TaskSucceeded}, nil
}

func main() {
	lis, err := net.Listen("tcp", Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	client, err := ent.Open("sqlite3", "file:test.sqlite?_loc=auto&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}

	ctx := context.Background()
	defer client.Close()

	// Run the auto migration tool.
	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

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
	certPem, certerr := static.ReadFile(CertFile)
	if certerr != nil {
		fmt.Println("File reading error", certerr)
		return
	}
	keyPem, keyerr := static.ReadFile(KeyFile)
	if keyerr != nil {
		fmt.Println("File reading error", keyerr)
		return
	}

	cert, tlserr := tls.X509KeyPair(certPem, keyPem)
	if tlserr != nil {
		fmt.Println("File reading error", tlserr)
		return
	}

	creds := credentials.NewServerTLSFromCert(&cert)
	s := grpc.NewServer(grpc.Creds(creds))

	log.Printf("Starting Laforge Server on port " + Port)

	pb.RegisterLaforgeServer(s, &Server{
		Client:                     client,
		UnimplementedLaforgeServer: pb.UnimplementedLaforgeServer{},
	})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
