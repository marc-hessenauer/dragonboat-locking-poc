package config

import (
	"github.com/lni/dragonboat/v4/config"
)

func GetNodeHostConfig(nodeID uint64, address string) config.NodeHostConfig {
	return config.NodeHostConfig{
		NodeHostDir:    "", // memory only
		RaftAddress:    address,
		RTTMillisecond: 200,
		EnableMetrics:  false,
		Expert: config.ExpertConfig{},
	}
}

func GetRaftConfig(clusterID, nodeID uint64) config.Config {
	return config.Config{
		ReplicaID:          nodeID,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    0, // disable snapshots
		CompactionOverhead: 0, // no compaction needed
	}
}
