package hash

import (
	"testing"
)

func NewHashCirclePrepare(replicas int) *HashCircle {
	hashCircle := NewHashCircle(HashCircleOpts{
		replicas: replicas,
	})
	return hashCircle
}

func TestHashCircle(t *testing.T) {
	keys := []string{"127.0.0.1:8080", "127.0.0.1:8081", "127.0.0.1:8082`"}
	replicas := 10
	hashCircle := NewHashCirclePrepare(replicas)
	hashCircle.Add(keys...)

	// test keys count
	if len(hashCircle.keys) != len(keys)*replicas {
		t.Errorf("expected %d keys, got %d", len(keys)*replicas, len(hashCircle.keys))
	}

	// test keys are sorted
	i := 0
	for i < len(keys)-1 {
		if hashCircle.keys[i] > hashCircle.keys[i+1] {
			t.Errorf("keys are not sorted, expected key[%d] = %d < key[%d] = %d", i, hashCircle.keys[i], i+1, hashCircle.keys[i+1])
		}
		i++
	}

}
