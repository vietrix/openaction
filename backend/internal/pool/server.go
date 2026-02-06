package pool

import (
	"context"
	"time"

	"github.com/google/uuid"

	"openaction/pkg/poolpb"
)

type Server struct {
	poolpb.UnimplementedPoolServiceServer
}

func (s *Server) Register(ctx context.Context, req *poolpb.RegisterRequest) (*poolpb.RegisterResponse, error) {
	assigned := uuid.NewString()
	if req.Info != nil && req.Info.Id != "" {
		assigned = req.Info.Id
	}
	return &poolpb.RegisterResponse{AssignedId: assigned}, nil
}

func (s *Server) Heartbeat(ctx context.Context, req *poolpb.HeartbeatRequest) (*poolpb.HeartbeatResponse, error) {
	return &poolpb.HeartbeatResponse{Ok: true}, nil
}

func (s *Server) FetchJob(ctx context.Context, req *poolpb.JobRequest) (*poolpb.JobResponse, error) {
	return &poolpb.JobResponse{JobId: "", Payload: ""}, nil
}

func (s *Server) ReportStep(ctx context.Context, req *poolpb.StepReport) (*poolpb.StepReportResponse, error) {
	_ = time.Now()
	return &poolpb.StepReportResponse{Ok: true}, nil
}
