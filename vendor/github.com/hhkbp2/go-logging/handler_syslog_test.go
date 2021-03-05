// +build !windows

package logging

import (
	"log/syslog"
	"testing"

	"github.com/hhkbp2/testify/require"
)

func TestSyslogHandler(t *testing.T) {
	defer Shutdown()
	handler, err := NewSyslogHandler(
		syslog.LOG_USER|syslog.LOG_DEBUG,
		"atag")
	require.Nil(t, err)
	logger := GetLogger("syslog")
	logger.SetLevel(LevelDebug)
	logger.AddHandler(handler)
	prefix := "test syslog handler "
	logger.Debugf(prefix + "Debug() a message")
	logger.Errorf(prefix + "Error() a message")
}
