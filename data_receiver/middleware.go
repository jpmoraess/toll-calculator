package main

import (
	"time"

	"github.com/jpmoraess/toll-calculator/common"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next DataProducer
}

func NewLogMiddleware(next DataProducer) DataProducer {
	return &LogMiddleware{
		next: next,
	}
}

func (l *LogMiddleware) ProduceData(data common.OBUData) error {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuID": data.OBUID,
			"lat":   data.Lat,
			"long":  data.Long,
			"took":  time.Since(start),
		}).Info("producing to kafka")
	}(time.Now())
	return l.next.ProduceData(data)
}
