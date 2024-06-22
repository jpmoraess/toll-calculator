package main

import (
	"fmt"

	"github.com/jpmoraess/toll-calculator/common"
)

type Aggregator interface {
	AggregateDistance(common.Distance) error
}

type Storer interface {
	Insert(common.Distance) error
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
	fmt.Println("processing and inserting distance in the storage:", distance)
	return i.store.Insert(distance)
}
