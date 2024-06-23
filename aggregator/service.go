package main

import (
	"github.com/jpmoraess/toll-calculator/common"
)

type Aggregator interface {
	AggregateDistance(common.Distance) error
	GetInvoice(int) (*common.Invoice, error)
}

type Storer interface {
	Insert(common.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) *InvoiceAggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance common.Distance) error {
	return i.store.Insert(distance)
}

func (i *InvoiceAggregator) GetInvoice(obuID int) (*common.Invoice, error) {
	distance, err := i.store.Get(obuID)
	if err != nil {
		return nil, err
	}
	return &common.Invoice{
		OBUID:         obuID,
		TotalDistance: distance,
		TotalAmount:   distance * 1.5,
	}, nil
}
