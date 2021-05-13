package server

//go:generate fileb0x assets.toml
import (
	"context"
	"fmt"

	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/agenttask"
	"github.com/gen0cide/laforge/ent/provisionedhost"
	pb "github.com/gen0cide/laforge/grpc/proto"
	"github.com/google/uuid"
)

var (
	Port     = ":50051"
	CertFile = "service.pem"
	KeyFile  = "service.key"
	// webPort  = ":5000"
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
	uuid, err := uuid.Parse(in.GetClientId())

	if err != nil {
		return nil, fmt.Errorf("failed casting UUID to UUID: %v", err)
	}

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
		SetAgentStatusToProvisionedHost(ph).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed Creating Agent Status: %v", err)
	}

	avalibleTasks, err := ph.QueryProvisionedHostToAgentTask().Where(
		agenttask.StateEQ(agenttask.StateAWAITING),
	).Exist(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to Query Agent Tasks: %v", err)
	}

	return &pb.HeartbeatReply{Status: message, AvalibleTasks: avalibleTasks}, nil
}

//GetTask Gets a task that needs to be run on the client and sends it over
func (s *Server) GetTask(ctx context.Context, in *pb.TaskRequest) (*pb.TaskReply, error) {
	uuid, err := uuid.Parse(in.GetClientId())
	if err != nil {
		return nil, fmt.Errorf("failed casting UUID to UUID: %v", err)
	}
	ph, err := s.Client.ProvisionedHost.Query().Where(provisionedhost.IDEQ(uuid)).Only(ctx)
	entAgentTask, err := ph.QueryProvisionedHostToAgentTask().Order(ent.Asc(agenttask.FieldNumber)).Where(agenttask.StateEQ(agenttask.StateAWAITING)).First(ctx)
	if err != nil {
		if err != err.(*ent.NotFoundError) {
			return &pb.TaskReply{Id: "", Command: pb.TaskReply_DEFAULT}, nil
		}
	} else {
		return nil, err
	}

	return &pb.TaskReply{
		Id: entAgentTask.ID.String(),
		Command: pb.TaskReply_Command(
			pb.TaskReply_Command_value[string(entAgentTask.Command)],
		),
		Args: entAgentTask.Args,
	}, nil
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
