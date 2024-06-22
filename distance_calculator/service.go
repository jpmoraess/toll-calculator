package main

import (
	"math"

	"github.com/jpmoraess/toll-calculator/common"
)

// We like the end our interface with (er) :shrug:
type CalculatorServicer interface {
	CalculateDistance(common.OBUData) (float64, error)
}

type CalculatorService struct {
	previousPoint []float64
}

func NewCalculatorService() *CalculatorService {
	return &CalculatorService{}
}

func (s *CalculatorService) CalculateDistance(data common.OBUData) (float64, error) {
	distance := 0.0
	if len(s.previousPoint) > 0 {
		distance = calculateDistance(s.previousPoint[0], s.previousPoint[1], data.Lat, data.Long)
	}
	s.previousPoint = []float64{data.Lat, data.Long}
	return distance, nil
}

func calculateDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2-x1, 2) + math.Pow(y2-y1, 2))
}
