package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/singnet/snet-daemon/config"
)

const defaultFormatterConfigJSON = `
	{
		"type": "json",
		"timestamp_format": "UTC"
	}`

const defaultOutputConfigJSON = `
	{
		"type": "file",
		"file_pattern": "/tmp/snet-daemon.%Y%m%d.log",
		"current_link": "/tmp/snet-daemon.log",
		"max_size_in_mb": 86400,
		"max_age_in_days": 604800,
		"rotation_count": 0
	}`

const defaultLogConfigJSON = `
	{
		"level": "info",
		"timezone": "UTC",
		"formatter": {
			"type": "json",
			"timestamp_format": "UTC"
		},
		"output": {
			"type": "file",
			"file_pattern": "/tmp/snet-daemon.%Y%m%d.log",
			"current_link": "/tmp/snet-daemon.log",
			"max_size_in_mb": 86400,
			"max_age_in_days": 604800,
			"rotation_count": 0
		}
	}`

func TestMain(m *testing.M) {
	result := m.Run()

	removeLogFiles("/tmp/snet-daemon*.log")
	removeLogFiles("/tmp/file-rotatelogs-test.*.log")

	os.Exit(result)
}

func removeLogFiles(pattern string) {
	var err error
	var files []string

	files, err = filepath.Glob(pattern)
	if err != nil {
		panic(fmt.Sprintf("Cannot find files using pattern: %v", err))
	}

	for _, file := range files {
		err = os.Remove(file)
		if err != nil {
			panic(fmt.Sprintf("Cannot remove file: %v, error: %v", file, err))
		}
	}
}

func TestCreateEncoderConfig(t *testing.T) {
	var formatterText = `
	{
	    "type": "text",
	    "timestamp_format": "UTC"
	}`

	config.LoadConfig(formatterText)
	_, err := createEncoderConfig()

	assert.Nil(t, err)
}

func TestGetLocationTimezone(t *testing.T) {
	logConfigText := `
	{
		"log": {
			"level": "info",
			"timezone": "UTC",
			"formatter": {
				"type": "json",
				"timestamp_format": "UTC"
			},
			"output": {
				"type": "file",
				"file_pattern": "/tmp/snet-daemon.%Y%m%d.log",
				"current_link": "/tmp/snet-daemon.log",
				"max_size_in_mb": 86400,
				"max_age_in_days": 604800,
				"rotation_count": 0
			}
		}
	}`

	config.LoadConfig(logConfigText)
	v := config.Vip()
	assert.NotNil(t, config.Vip())

	timezoneViper := v.GetString(LogTimezoneKey)
	assert.NotEmpty(t, timezoneViper)

	timezone, err := getLocationTimezone()
	assert.NoError(t, err)
	assert.NotNil(t, timezone)

	currentTime := time.Now()
	assert.Equal(t, currentTime.Format(timezoneViper), currentTime.Format(timezone.String()))
}

func TestGetLoggerEncoder(t *testing.T) {
	logConfigText := `{
		"log": {
			"level": "info",
			"timezone": "UTC",
			"formatter": {
				"type": "json",
				"timestamp_format": "UTC"
			},
			"output": {
				"type": "file",
				"file_pattern": "/tmp/snet-daemon.%Y%m%d.log",
				"current_link": "/tmp/snet-daemon.log",
				"max_size_in_mb": 86400,
				"max_age_in_days": 604800,
				"rotation_count": 0
			}
		}
	}`

	config.LoadConfig(logConfigText)
	assert.NotNil(t, config.Vip())

	encoderConf, err := createEncoderConfig()
	assert.NotNil(t, encoderConf)
	assert.Nil(t, err)

	encoder, err := createLoggerEncoder(encoderConf)
	assert.NotNil(t, encoder)
}

func TestGetLoggerLevel(t *testing.T) {
	level := "info"
	logLevel, err := getLoggerLevel(level)
	assert.Nil(t, err)
	assert.Equal(t, zap.InfoLevel, logLevel)

	level = "debug"
	logLevel, err = getLoggerLevel(level)
	assert.Nil(t, err)
	assert.Equal(t, zap.DebugLevel, logLevel)

	level = "warn"
	logLevel, err = getLoggerLevel(level)
	assert.Nil(t, err)
	assert.Equal(t, zap.WarnLevel, logLevel)

	level = "error"
	logLevel, err = getLoggerLevel(level)
	assert.Nil(t, err)
	assert.Equal(t, zap.ErrorLevel, logLevel)

	level = "panic"
	logLevel, err = getLoggerLevel(level)
	assert.Nil(t, err)
	assert.Equal(t, zap.PanicLevel, logLevel)

	level = "wrong_level"
	logLevel, err = getLoggerLevel(level)
	assert.NotNil(t, err)
	assert.Equal(t, "wrong string for level: wrong_level. Availible options: debug, info, warn, error, panic", err.Error())
}

func TestFormatFileName(t *testing.T) {
	mockTime := time.Date(2024, 7, 4, 12, 34, 56, 789000000, time.UTC)

	filePatternName := "./snet-daemon.%Y-----%m-----%d--%M.log"
	expectedFileName := "./snet-daemon.2024-----07-----04--34.log"
	fileName, err := formatFileName(filePatternName, mockTime)
	assert.Nil(t, err)

	assert.Equal(t, expectedFileName, fileName)

	filePatternNameError := "./snet-daemon.%x%l%f.log"
	_, err = formatFileName(filePatternNameError, mockTime)
	assert.NotNil(t, err)
}
