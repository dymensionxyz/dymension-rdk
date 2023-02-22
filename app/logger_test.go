package app_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dymensionxyz/rollapp/app"
	"github.com/dymensionxyz/rollapp/utils"
)

func TestLogLevel(t *testing.T) {
	var buf bytes.Buffer

	logger := app.NewLogger("", "error", nil)
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

	logger := app.NewLogger("", "info", nil)
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

	logger := app.NewLogger("", "info", nil)
	logger.SetOutput(&buf)

	logger1 := logger.With("module", "module1")
	logger2 := logger1.With("arg", "custom_arg") //on top logger1

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

	moduleOverrides := utils.ConvertStringToStringMap("module2:error", ",", ":")
	logger := app.NewLogger("", "info", moduleOverrides)
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
