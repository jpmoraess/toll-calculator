package client

import (
	"context"
	"log"

	"github.com/jpmoraess/toll-calculator/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AggregatorGRPCClient struct {
	Client common.AggregatorClient
}

func NewAggregatorGRPCClient(endpoint string) *AggregatorGRPCClient {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c := common.NewAggregatorClient(conn)

	return &AggregatorGRPCClient{
		Client: c,
	}
}

func (c *AggregatorGRPCClient) AggregateInvoice(ctx context.Context, distance common.Distance) error {
	req := &common.AggregateRequest{
		ObuID: int32(distance.OBUID),
		Value: distance.Value,
		Unix:  distance.Unix,
	}
	_, err := c.Client.Aggregate(ctx, req)
	if err != nil {
		return err
	}
	return nil
}
