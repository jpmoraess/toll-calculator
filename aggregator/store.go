package main

import "github.com/jpmoraess/toll-calculator/common"

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (m *MemoryStore) Insert(distance common.Distance) error {
	m.data[distance.OBUID] += distance.Value
	return nil
}
