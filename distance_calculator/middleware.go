package main

import (
	"time"

	"github.com/jpmoraess/toll-calculator/common"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next CalculatorServicer
}

func NewLogMiddleware(next CalculatorServicer) CalculatorServicer {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) CalculateDistance(data common.OBUData) (distance float64, err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":     time.Since(start),
			"err":      err,
			"distance": distance,
		}).Info("calculate distance")
	}(time.Now())
	return m.next.CalculateDistance(data)
}
