package main

import (
	"errors"
	"fmt"
)

//Bloqueio que é descrito no trabalho
type Lock struct {
	ItemID   string
	TrID     int
	Scope    string
	Duration string
	OpType   string
}

func (l *Lock) SetScope(value string) error {
	if value != "O" && value != "P" {
		return errors.New("O valor do atributo Scope deve ser O ou P")
	}
	l.Scope = value
	return nil
}

func (l *Lock) SetDuration(value string) error {
	if value != "L" && value != "C" {
		return errors.New("O valor do atributo Duration deve ser L ou C")
	}

	l.Duration = value
	return nil
}

func (l *Lock) SetOpType(value string) error {
	if value != "R" && value != "W" {
		return errors.New("O valor do atributo OpType deve ser R ou W")
	}

	l.OpType = value
	return nil
}

func NewLock(itemID string, trID int, scope, duration, opType string) (*Lock, error) {
	lock := &Lock{
		ItemID: itemID,
		TrID:   trID,
	}

	if err := lock.SetScope(scope); err != nil {
		return nil, err
	}

	if err := lock.SetDuration(duration); err != nil {
		return nil, err
	}

	if err := lock.SetOpType(opType); err != nil {
		return nil, err
	}

	return lock, nil
}

type LockTable struct {
	Locks          []Lock
	IsolationLevel int // 0 a 3
}

func (lt *LockTable) PrintLockTable() {
	fmt.Println("Lock Table:")
	for _, lock := range lt.Locks {
		fmt.Printf("ItemID: %s, TrID: %d, Scope: %s, Duration: %s, OpType: %s\n", lock.ItemID, lock.TrID, lock.Scope, lock.Duration, lock.OpType)
	}
}
