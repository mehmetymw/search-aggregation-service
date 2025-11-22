package logger

import "github.com/mehmetymw/search-aggregation-service/backend/domain/ports"

type logField struct {
	key   string
	value any
}

func (f logField) Key() string {
	return f.key
}

func (f logField) Value() any {
	return f.value
}

func String(key, value string) ports.Field {
	return logField{key: key, value: value}
}

func Int(key string, value int) ports.Field {
	return logField{key: key, value: value}
}

func Int64(key string, value int64) ports.Field {
	return logField{key: key, value: value}
}

func Float64(key string, value float64) ports.Field {
	return logField{key: key, value: value}
}

func Bool(key string, value bool) ports.Field {
	return logField{key: key, value: value}
}

func Error(err error) ports.Field {
	return logField{key: "error", value: err}
}

func Any(key string, value any) ports.Field {
	return logField{key: key, value: value}
}
