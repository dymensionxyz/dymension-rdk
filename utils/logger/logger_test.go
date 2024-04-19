package logger_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/dymensionxyz/dymension-rdk/utils/logger"
	log "github.com/dymensionxyz/dymension-rdk/utils/logger"
)

func TestLogLevel(t *testing.T) {
	var buf bytes.Buffer

	logger := log.NewLogger("", 0, "error", nil)
	logger.SetOutput(&buf)

	logger.Debug("debug msg")
	logger.Info("info msg")
	logger.Error("error msg")
	msg := strings.TrimSpace(buf.String())

	assert.Contains(t, msg, "error msg", "error msg wasn't logged")
	assert.NotContains(t, msg, "info")
	assert.NotContains(t, msg, "debug")
}

func TestMultipleLoggerWithMethod(t *testing.T) {
	var buf bytes.Buffer

	logger := log.NewLogger("", 0, "info", nil)
	logger.SetOutput(&buf)

	logger1 := logger.With("module", "logger1")
	logger2 := logger.With("module", "logger2")

	logger1.Info("testing")
	msg := strings.TrimSpace(buf.String())
	assert.Contains(t, msg, "logger1", "logger1 didn't have the module name as expected")
	assert.NotContains(t, msg, "logger2", "logger1 shouldn't have logger2 in it's log")

	buf.Reset()
	logger2.Info("testing")
	msg = strings.TrimSpace(buf.String())
	assert.Contains(t, msg, "logger2", "logger2 didn't have the module name as expected")
	assert.NotContains(t, msg, "logger1", "logger2 shouldn't have logger1 in it's log")
}

func TestMultipleWithCalls(t *testing.T) {
	// t.Skip("nested With calls not supported")
	var buf bytes.Buffer

	logger := log.NewLogger("", 0, "info", nil)
	logger.SetOutput(&buf)

	logger1 := logger.With("module", "module1")
	logger2 := logger1.With("arg", "custom_arg") // on top logger1

	logger1.Info("testing")
	msg := strings.TrimSpace(buf.String())
	assert.Contains(t, msg, "module1", "logger didn't have the module name as expected")
	assert.NotContains(t, msg, "custom_arg", "logger shouldn't have logger2 in it's log")

	buf.Reset()
	logger2.Info("testing")
	msg = strings.TrimSpace(buf.String())
	assert.Contains(t, msg, "module1", "logger2 didn't have the module name as expected")
	assert.Contains(t, msg, "custom_arg", "logger2 didn't have the module name as expected")
}

func TestModuleOverrideLevel(t *testing.T) {
	var buf bytes.Buffer

	moduleOverrides := logger.ConvertStringToStringMap("module2:error", ",", ":")
	logger := log.NewLogger("", 0, "info", moduleOverrides)
	logger.SetOutput(&buf)

	logger1 := logger.With("module", "module1")
	logger2 := logger.With("module", "module2")

	/* -------------------------- check module 1 logger ------------------------- */
	logger1.Debug("testing debug")
	logger1.Info("testing info")
	logger1.Error("testing error")
	msg := strings.TrimSpace(buf.String())
	assert.Contains(t, msg, "testing error", "logger expected to log error")
	assert.Contains(t, msg, "testing info", "logger expected to log info")
	assert.NotContains(t, msg, "testing debug", "logger expected to NOT log debug")

	buf.Reset()

	/* -------------------------- check module 2 logger ------------------------- */

	logger2.Debug("testing debug")
	logger2.Info("testing info")
	logger2.Error("testing error")
	msg = strings.TrimSpace(buf.String())
	assert.Contains(t, msg, "level=error", "logger expected to log error")
	assert.NotContains(t, msg, "level=info", "logger expected to NOT log info")
	assert.NotContains(t, msg, "level=debug", "logger expected to NOT log debug")
}

func TestMaxFileSize(t *testing.T) {
	seed := "test-" + time.Now().Format("20060102-150405")
	logFileName := fmt.Sprintf("test-log-%s.log", seed)

	logDir := filepath.Join("/tmp", seed)
	err := os.Mkdir(logDir, 0o700)
	assert.NoError(t, err)

	logPath := filepath.Join(logDir, logFileName)
	_, err = os.Create(logPath)
	assert.NoError(t, err)

	defer func() {
		os.RemoveAll(logDir)
	}()

	logger := log.NewLogger(logPath, 1, "info", nil)
	fillUpLog(&logger)

	files, err := os.ReadDir(logDir)
	assert.NoError(t, err)

	assert.Greater(t, len(files), 1)
}

func fillUpLog(logger *log.Logger) {
	for i := 0; i < 100000; i++ {
		logger.Infof("Log message %d", i)
	}
}
