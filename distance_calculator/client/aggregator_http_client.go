package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jpmoraess/toll-calculator/common"
)

type AggregatorHttpClient struct {
	Endpoint string
}

func NewAggregatorHttpClient(endpoint string) *AggregatorHttpClient {
	return &AggregatorHttpClient{
		Endpoint: endpoint,
	}
}

func (c *AggregatorHttpClient) AggregateInvoice(ctx context.Context, distance common.Distance) error {
	b, err := json.Marshal(distance)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the service responded with non 200 status code: %d", resp.StatusCode)
	}
	return nil
}
