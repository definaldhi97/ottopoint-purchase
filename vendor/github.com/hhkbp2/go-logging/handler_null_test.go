package logging

import (
	"github.com/hhkbp2/testify/require"
	"testing"
)

func TestNullHandler(t *testing.T) {
	handler := NewNullHandler()
	logger := GetLogger("null")
	logger.AddHandler(handler)
	require.Equal(t, 1, len(logger.GetHandlers()))
	message := "test"
	logger.Debugf(message)
	logger.Fatalf(message)
}
