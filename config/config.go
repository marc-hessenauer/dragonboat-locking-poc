package config

import (
	"github.com/lni/dragonboat/v3/config"
)

func GetNodeHostConfig(nodeID uint64, address string) config.NodeHostConfig {
	return config.NodeHostConfig{
		NodeHostDir:    "", // memory only
		RaftAddress:    address,
		RTTMillisecond: 200,
		EnableMetrics:  false,
		Expert: config.ExpertConfig{
			DisableLogDB: true, // in-memory only
		},
	}
}

func GetRaftConfig(clusterID, nodeID uint64) config.Config {
	return config.Config{
		NodeID:             nodeID,
		ClusterID:          clusterID,
		ElectionRTT:        10,
		HeartbeatRTT:       1,
		CheckQuorum:        true,
		SnapshotEntries:    0, // disable snapshots
		CompactionOverhead: 0, // no compaction needed
	}
}
