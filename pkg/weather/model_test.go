package weather

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJSONTime(t *testing.T) {
	jsonTime := NewJSONTime(time.Unix(1485789600, 0))
	assert.Equal(t, int64(1485789600), jsonTime.Unix())
}
