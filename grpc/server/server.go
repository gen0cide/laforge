package server

//go:generate fileb0x assets.toml
import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/agenttask"
	"github.com/gen0cide/laforge/ent/provisionedhost"
	pb "github.com/gen0cide/laforge/grpc/proto"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	Port     = ":50051"
	CertFile = "service.pem"
	KeyFile  = "service.key"

	// webPort  = ":5000"
)

type Server struct {
	Client *ent.Client
	pb.UnimplementedLaforgeServer
	RDB *redis.Client
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
	// message := fmt.Sprintf("Recived ID: %v | Hostname: %v | Uptime: %v | Boot Time: %v| Number of Running Processes: %v| OS Arch: %v| Host ID: %v| Load1: %v| Load5: %v| Load15: %v| Total Memory: %v| Avalible Memory: %v| Used Memory: %v", in.GetClientId(), in.GetHostname(), in.GetUptime(), in.GetBoottime(), in.GetNumprocs(), in.GetOs(), in.GetHostid(), in.GetLoad1(), in.GetLoad5(), in.GetLoad15(), ByteCountIEC(in.GetTotalmem()), ByteCountIEC(in.GetFreemem()), ByteCountIEC(in.GetUsedmem()))
	message := fmt.Sprintf("Heartbeat Recived: %v", time.Now().Unix())
	uuid, err := uuid.Parse(in.GetClientId())

	if err != nil {
		logrus.Errorf("GRPC SERVER ERROR: failed casting UUID to UUID: %v", err)
		return &pb.HeartbeatReply{Status: message, AvalibleTasks: false}, nil
	}
	ph, err := s.Client.ProvisionedHost.Query().Where(provisionedhost.IDEQ(uuid)).Only(ctx)
	if err != nil {
		logrus.Errorf("GRPC SERVER ERROR: Cannot find client %v. Error: %v", in.GetClientId(), err)
		return &pb.HeartbeatReply{Status: message, AvalibleTasks: false}, nil
	}
	pn, err := ph.QueryProvisionedHostToProvisionedNetwork().Only(ctx)
	if err != nil {
		logrus.Errorf("GRPC SERVER ERROR: Cannot find client %v network. Error: %v", in.GetClientId(), err)
		return &pb.HeartbeatReply{Status: message, AvalibleTasks: false}, nil
	}
	b, err := pn.QueryProvisionedNetworkToBuild().Only(ctx)
	if err != nil {
		logrus.Errorf("GRPC SERVER ERROR: Cannot find client %v build. Error: %v", in.GetClientId(), err)
		return &pb.HeartbeatReply{Status: message, AvalibleTasks: false}, nil
	}
	existingEntAgentStatus, err := ph.QueryProvisionedHostToAgentStatus().First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			createdEntAgentStatus, err := s.Client.AgentStatus.
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
				SetTimestamp(time.Now().Unix()).
				SetAgentStatusToProvisionedHost(ph).
				SetAgentStatusToProvisionedNetwork(pn).
				SetAgentStatusToBuild(b).
				Save(ctx)

			if err != nil {
				logrus.Errorf("GRPC SERVER ERROR: failed Creating Agent Status: %v", err)
				return &pb.HeartbeatReply{Status: message, AvalibleTasks: false}, nil
			}

			s.RDB.Publish(ctx, "newAgentStatus", createdEntAgentStatus.ID.String())
		} else {
			logrus.Errorf("GRPC SERVER ERROR: failed Query Agent Status: %v", err)
			return &pb.HeartbeatReply{Status: message, AvalibleTasks: false}, nil
		}
	} else {
		existingEntAgentStatus, err = existingEntAgentStatus.
			Update().
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
			SetTimestamp(time.Now().Unix()).
			SetAgentStatusToProvisionedHost(ph).
			SetAgentStatusToProvisionedNetwork(pn).
			SetAgentStatusToBuild(b).
			Save(ctx)

		if err != nil {
			logrus.Errorf("GRPC SERVER ERROR: failed Update Agent Status: %v", err)
			return &pb.HeartbeatReply{Status: message, AvalibleTasks: false}, nil
		}

		s.RDB.Publish(ctx, "newAgentStatus", existingEntAgentStatus.ID.String())
	}
	// logrus.Debugf("GRPC SERVER DEBUG: Agent for client %v has sent heartbeat", in.GetClientId())

	avalibleTasks, err := ph.QueryProvisionedHostToAgentTask().Where(
		agenttask.Or(
			agenttask.StateEQ(agenttask.StateAWAITING),
			agenttask.StateEQ(agenttask.StateINPROGRESS),
		),
	).Exist(ctx)
	if err != nil {
		logrus.Errorf("GRPC SERVER ERROR: failed to Query Agent Tasks: %v", err)
		return &pb.HeartbeatReply{Status: message, AvalibleTasks: false}, nil
	}

	return &pb.HeartbeatReply{Status: message, AvalibleTasks: avalibleTasks}, nil
}

//GetTask Gets a task that needs to be run on the client and sends it over
func (s *Server) GetTask(ctx context.Context, in *pb.TaskRequest) (*pb.TaskReply, error) {
	uuid, err := uuid.Parse(in.GetClientId())
	if err != nil {
		logrus.Errorf("GRPC SERVER ERROR: failed casting UUID to UUID: %v", err)
		return &pb.TaskReply{Id: "", Command: pb.TaskReply_DEFAULT}, nil
	}
	ph, err := s.Client.ProvisionedHost.Query().Where(provisionedhost.IDEQ(uuid)).Only(ctx)
	if err != nil {
		logrus.Errorf("GRPC SERVER ERROR: Cannot find client %v. Error: %v", in.GetClientId(), err)
		return &pb.TaskReply{Id: "", Command: pb.TaskReply_DEFAULT}, nil
	}
	entAgentTask, err := ph.QueryProvisionedHostToAgentTask().Order(ent.Asc(agenttask.FieldNumber)).Where(
		agenttask.Or(
			agenttask.StateEQ(agenttask.StateAWAITING),
			agenttask.StateEQ(agenttask.StateINPROGRESS),
		)).First(ctx)
	if err != nil {
		if !ent.IsNotFound(err) {
			logrus.Errorf("GRPC SERVER ERROR: Cannot find agent task for client %v. Error: %v", in.GetClientId(), err)
		}
		return &pb.TaskReply{Id: "", Command: pb.TaskReply_DEFAULT}, nil
	}
	// logrus.Debugf("GRPC SERVER DEBUG: Agent for client %v has requested task", in.GetClientId())
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
	if in.GetTaskId() == "" {
		return &pb.TaskStatusReply{Status: "ERROR"}, nil
	}
	uuid, err := uuid.Parse(in.GetTaskId())
	if err != nil {
		logrus.Errorf("GRPC SERVER ERROR: failed casting UUID to UUID: %v", err)
		return &pb.TaskStatusReply{Status: "ERROR"}, nil
	}
	entAgentTask, err := s.Client.AgentTask.Query().Where(agenttask.IDEQ(uuid)).First(ctx)
	if err != nil {
		logrus.Errorf("GRPC SERVER ERROR: failed Querying Agent Task %v: %v", uuid, err)
		return &pb.TaskStatusReply{Status: "ERROR"}, nil
	}
	output := strings.ReplaceAll(in.GetOutput(), "🔥", "\n")
	err = entAgentTask.Update().
		SetState(agenttask.State(in.GetStatus())).
		SetErrorMessage(in.GetErrorMessage()).
		SetOutput(output).
		Exec(ctx)
	if err != nil {
		logrus.Errorf("GRPC SERVER ERROR: failed Updating Agent Task %v: %v", uuid, err)
		return &pb.TaskStatusReply{Status: "ERROR"}, nil
	}
	s.RDB.Publish(ctx, "updatedAgentTask", entAgentTask.ID.String())
	return &pb.TaskStatusReply{Status: in.GetStatus()}, nil
}
