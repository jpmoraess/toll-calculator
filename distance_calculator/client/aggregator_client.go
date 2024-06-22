package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jpmoraess/toll-calculator/common"
)

type AggregatorClient struct {
	Endpoint string
}

func NewAggregatorClient(endpoint string) *AggregatorClient {
	return &AggregatorClient{
		Endpoint: endpoint,
	}
}

func (c *AggregatorClient) AggregateInvoice(distance common.Distance) error {
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
