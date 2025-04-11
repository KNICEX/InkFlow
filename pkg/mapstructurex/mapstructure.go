package mapstructurex

import (
	"github.com/mitchellh/mapstructure"
	"reflect"
	"time"
)

func timeHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if t != reflect.TypeOf(time.Time{}) {
		return data, nil
	}

	if t != reflect.TypeOf(time.Time{}) {
		return data, nil
	}

	switch f.Kind() {
	case reflect.String:
		return time.Parse(time.RFC3339, data.(string))
	case reflect.Int64:
		return time.Unix(data.(int64), 0), nil
	default:
		return data, nil
	}
}

func Decode(input any, output any) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook: mapstructure.ComposeDecodeHookFunc(timeHook),
		Result:     output,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(input)
}
