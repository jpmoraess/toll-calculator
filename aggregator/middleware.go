package main

import (
	"time"

	"github.com/jpmoraess/toll-calculator/common"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{next: next}
}

func (m *LogMiddleware) AggregateDistance(data common.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("AggregateDistance")
	}(time.Now())
	err = m.next.AggregateDistance(data)
	return
}

func (m *LogMiddleware) GetInvoice(obuID int) (invoice *common.Invoice, err error) {
	defer func(start time.Time) {
		var (
			amount   float64
			distance float64
		)
		if invoice != nil {
			amount = invoice.TotalAmount
			distance = invoice.TotalDistance
		}

		logrus.WithFields(logrus.Fields{
			"took":     time.Since(start),
			"err":      err,
			"obuID":    obuID,
			"amount":   amount,
			"distance": distance,
		}).Info("GetInvoice")
	}(time.Now())
	return m.next.GetInvoice(obuID)
}
