package client

import (
	"context"

	"github.com/jpmoraess/toll-calculator/common"
)

type AggregatorClient interface {
	AggregateInvoice(context.Context, common.Distance) error
}
