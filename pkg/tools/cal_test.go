package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	assert.Equal(t, 2, Add(1, 1))
}

func BenchmarkAdd(b *testing.B) {
	for n := 0; n < b.N; n++ {
		Add(1, 1)
	}
}
