package locking

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/lni/dragonboat/v4/statemachine"
)

type LockCommand struct {
	Resource string `json:"resource"`
	Action   string `json:"action"` // "lock" or "unlock"
}

type LockStateMachine struct {
	lockMap map[string]bool
}

func NewLockStateMachine() statemachine.IStateMachine {
	return &LockStateMachine{
		lockMap: make(map[string]bool),
	}
}

// Apply a Raft proposal
func (s *LockStateMachine) Update(data []byte) (uint64, error) {
	var cmd LockCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		return 0, err
	}

	switch cmd.Action {
	case "lock":
		s.lockMap[cmd.Resource] = true
	case "unlock":
		s.lockMap[cmd.Resource] = false
	default:
		return 0, fmt.Errorf("unknown action: %s", cmd.Action)
	}
	return 1, nil
}

// Lookup (not replicated, local)
func (s *LockStateMachine) Lookup(query interface{}) (interface{}, error) {
	resource, ok := query.(string)
	if !ok {
		return nil, fmt.Errorf("expected string")
	}
	locked := s.lockMap[resource]
	return locked, nil
}

func (s *LockStateMachine) SaveSnapshot(w io.Writer, f statemachine.ISnapshotFileCollection, done <-chan struct{}) error {
	return nil // no snapshot support
}

func (s *LockStateMachine) RecoverFromSnapshot(r io.Reader, f []statemachine.SnapshotFile, done <-chan struct{}) error {
	return nil // no snapshot support
}

func (s *LockStateMachine) Close() error {
	return nil
}
