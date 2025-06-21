package limiter

import "hash/fnv"

type ShardSelector struct {
	nodes    []string
	strategy string
}

func NewShardSelector(nodes []string, strategy string) *ShardSelector {
	return &ShardSelector{
		nodes:    nodes,
		strategy: strategy,
	}
}

func (s *ShardSelector) GetRedisURL(key string) string {
	if len(s.nodes) == 0 {
		return "redis://localhost:6379/0" // fallback
	}

	switch s.strategy {
	case "hash_mod":
		return s.hashModSelect(key)
	case "consistent_hash":
		return s.consistentHashSelect(key)
	default:
		return s.nodes[0] // fallback
	}
}

func (s *ShardSelector) hashModSelect(key string) string {
	hash := s.hash32(key)
	nodeIndex := hash % uint32(len(s.nodes))
	return s.nodes[nodeIndex]
}

func (s *ShardSelector) consistentHashSelect(key string) string {
	keyHash := s.hash32(key)

	var selectedNode string
	minDistance := uint32(^uint32(0)) // max uint32

	for _, node := range s.nodes {
		nodeHash := s.hash32(node)
		var distance uint32
		if nodeHash >= keyHash {
			distance = nodeHash - keyHash
		} else {
			distance = (^uint32(0) - keyHash) + nodeHash
		}

		if distance < minDistance {
			minDistance = distance
			selectedNode = node
		}
	}

	if selectedNode == "" {
		return s.nodes[0]
	}

	return selectedNode
}

func (s *ShardSelector) hash32(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}
