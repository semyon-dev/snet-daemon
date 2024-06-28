package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	LogLevelKey    = "log.level"
	LogTimezoneKey = "log.timezone"

	LogFormatterTypeKey   = "log.formatter.type"
	LogTimestampFormatKey = "log.formatter.timestamp_format"

	LogOutputTypeKey        = "log.output.type"
	LogOutputFilePatternKey = "log.output.file_pattern"
	LogOutputCurrentLinkKey = "log.output.current_link"
	LogRotationTimeKey      = "log.output.rotation_time_in_sec"
	LogMaxAgeKey            = "log.output.max_age_in_sec"
	LogRotationCountKey     = "log.output.rotation_count"
)

// InitLogger initializes logger using configuration provided by viper
// instance.
//
// Function designed to configure few different loggers with different
// formatter and output settings. To achieve this viper configuration
// contains separate sections for each logger, each output and
// each formatter.

func Initialize(config *viper.Viper) {

	levelString := config.GetString(LogLevelKey)
	level, err := getLoggerLevel(levelString)
	if err != nil {
		panic(fmt.Errorf("failed to get logger level: %v", err))
	}

	encoderConfig, err := createEncoderConfig(config)
	if err != nil {
		panic(fmt.Errorf("failed to create encoder config, error: %v", err))
	}

	encoder, err := createLoggerEncoder(config, encoderConfig)
	if err != nil {
		panic(fmt.Errorf("failed to get encoder, error: %v", err))
	}

	writerSyncer, err := createLoggerWriterSyncer(config)
	if err != nil {
		panic(fmt.Errorf("failed to get logger writer, error: %v", err))
	}

	core := zapcore.NewCore(encoder, writerSyncer, level)
	logger := zap.New(core)
	zap.ReplaceGlobals(logger)

	// for _, hookConfigName := range config.GetStringSlice(LogHooksKey) {
	// 	err = addHookByConfig(logger, config.Sub(hookConfigName))
	// 	if err != nil {
	// 		return fmt.Errorf("Unable to add log hook \"%v\", error: %v", hookConfigName, err)
	// 	}
	// }

	logger.Info("Logger initialized")
}

func getLoggerLevel(levelString string) (zapcore.Level, error) {
	switch levelString {
	case "debug":
		return zap.DebugLevel, nil
	case "info":
		return zap.InfoLevel, nil
	case "warn":
		return zap.WarnLevel, nil
	case "error":
		return zap.ErrorLevel, nil
	case "panic":
		return zap.PanicLevel, nil
	default:
		return zapcore.Level(0), fmt.Errorf("wrong string for level: %v. Availible options: debug, info, warn, error, panic", levelString)
	}
}

func createEncoderConfig(config *viper.Viper) (*zapcore.EncoderConfig, error) {
	location, err := getLocationTimezone(config)
	if err != nil {
		return nil, err
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	if config.GetString(LogTimestampFormatKey) != "" {
		customTimeFormat := config.GetString(LogTimestampFormatKey)
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.In(location).Format(customTimeFormat))
		}
	} else {
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.In(location).Format(time.RFC3339))
		}
	}
	return &encoderConfig, nil
}

func getLocationTimezone(config *viper.Viper) (location *time.Location, err error) {
	timezone := config.GetString(LogTimezoneKey)
	location, err = time.LoadLocation(timezone)
	return
}

func createLoggerEncoder(config *viper.Viper, encoderConfig *zapcore.EncoderConfig) (zapcore.Encoder, error) {
	var encoder zapcore.Encoder
	switch config.GetString(LogFormatterTypeKey) {
	case "json":
		encoder = zapcore.NewJSONEncoder(*encoderConfig)
	case "text":
		encoder = zapcore.NewConsoleEncoder(*encoderConfig)
	default:
		return nil, fmt.Errorf("unsupported log formatter type: %v", config.GetString(LogFormatterTypeKey))
	}
	return encoder, nil
}

func createLoggerWriterSyncer(config *viper.Viper) (zapcore.WriteSyncer, error) {

	var ws zapcore.WriteSyncer
	switch config.GetString(LogOutputTypeKey) {
	case "stdout":
		ws = zapcore.AddSync(os.Stdout)
	case "stderr":
		ws = zapcore.AddSync(os.Stderr)
	case "file":
		ws = zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.GetString(LogOutputCurrentLinkKey),
			MaxSize:    config.GetInt(LogRotationTimeKey) / 1024 / 1024, // Convert from seconds to megabytes
			MaxAge:     config.GetInt(LogMaxAgeKey) / 86400,             // Convert from seconds to days
			MaxBackups: config.GetInt(LogRotationCountKey),
			Compress:   true,
		})

		// if _, err := os.Lstat(); err != nil {
		// 	if config.GetString(LogOutputCurrentLinkKey) != "" {
		// 		if err := os.Symlink(config.GetString(LogOutputFilePatternKey), config.GetString(LogOutputCurrentLinkKey)); err != nil {
		// 			return ws, fmt.Errorf("failed to create symlink: %v", err)
		// 		}
		// 	}
		// }
	default:
		return ws, fmt.Errorf("unsupported log output type: %v", config.GetString(LogOutputTypeKey))
	}
	return ws, nil
}
