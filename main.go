package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lni/dragonboat/v4"
	"github.com/lni/dragonboat/v4/statemachine"
	"github.com/marc-hessenauer/dragonboat-locking-poc/config"
	"github.com/marc-hessenauer/dragonboat-locking-poc/locking"
)

func main() {
	// Setup
	nodeID := uint64(1)
	clusterID := uint64(100)
	address := "localhost:25000"

	nhc := config.GetNodeHostConfig(nodeID, address)
	raftConfig := config.GetRaftConfig(clusterID, nodeID)

	nh, err := dragonboat.NewNodeHost(nhc)
	if err != nil {
		panic(err)
	}
	defer nh.Close()

	peers := map[uint64]string{1: address}
	if err := nh.StartReplica(peers, false, func(clusterID uint64, nodeID uint64) statemachine.IStateMachine {
		return locking.NewLockStateMachine()
	}, raftConfig); err != nil {
		panic(err)
	}

	time.Sleep(2 * time.Second) // wait for Raft to stabilize

	// Simulate client proposal
	cmd := locking.LockCommand{
		Resource: "printer-1",
		Action:   "lock",
	}
	data, _ := json.Marshal(cmd)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	session := nh.GetNoOPSession(clusterID)
	_, err = nh.SyncPropose(ctx, session, data)
	if err != nil {
		panic(err)
	}

	// Lookup
	result, err := nh.SyncRead(ctx, clusterID, "printer-1")
	if err != nil {
		panic(err)
	}

	locked, _ := result.(bool)
	fmt.Println("Is 'printer-1' locked? ->", locked)

	// Keep alive
	select {}
}
