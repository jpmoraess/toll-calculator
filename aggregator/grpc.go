package main

import (
	"context"

	"github.com/jpmoraess/toll-calculator/common"
)

type GRPCServer struct {
	common.UnimplementedAggregatorServer
	service Aggregator
}

func NewGRPCServer(service Aggregator) *GRPCServer {
	return &GRPCServer{
		service: service,
	}
}

// transport layer
// business layer

func (s *GRPCServer) Aggregate(ctx context.Context, request *common.AggregateRequest) (*common.None, error) {
	distance := common.Distance{
		OBUID: int(request.ObuID),
		Value: request.Value,
		Unix:  request.Unix,
	}
	if err := s.service.AggregateDistance(distance); err != nil {
		return nil, err
	}
	return &common.None{}, nil
}
