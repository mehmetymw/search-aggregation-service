package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/mehmetymw/search-aggregation-service/backend/domain/ports"
)

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger() (*ZapLogger, error) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(config),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &ZapLogger{logger: logger}, nil
}

func (l *ZapLogger) Info(msg string, fields ...ports.Field) {
	l.logger.Info(msg, l.convertFields(fields)...)
}

func (l *ZapLogger) Error(msg string, fields ...ports.Field) {
	l.logger.Error(msg, l.convertFields(fields)...)
}

func (l *ZapLogger) Debug(msg string, fields ...ports.Field) {
	l.logger.Debug(msg, l.convertFields(fields)...)
}

func (l *ZapLogger) Warn(msg string, fields ...ports.Field) {
	l.logger.Warn(msg, l.convertFields(fields)...)
}

func (l *ZapLogger) With(fields ...ports.Field) ports.Logger {
	return &ZapLogger{logger: l.logger.With(l.convertFields(fields)...)}
}

func (l *ZapLogger) Sync() error {
	return l.logger.Sync()
}

func (l *ZapLogger) convertFields(fields []ports.Field) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key(), f.Value()))
	}
	return zapFields
}
