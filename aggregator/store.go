package main

import "github.com/jpmoraess/toll-calculator/common"

type MemoryStore struct{}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{}
}

func (m *MemoryStore) Insert(distance common.Distance) error {
	return nil
}
