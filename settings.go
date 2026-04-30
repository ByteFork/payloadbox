package main

import (
	"reflect"
	"strconv"

	"github.com/ByteFork/payloadbox/internal/env"
)

type ServerSettings struct {
	Address           string `default:":8080" env:"LISTEN_ADDRESS"`
	MaxBodySizeBytes  int64  `default:"5120" env:"MAX_BODY_SIZE_BYTES"`
	MaxRecordsToStore int    `default:"200" env:"MAX_RECORDS_TO_STORE"`
	LogRequests       bool   `default:"true" env:"LOG_HTTP_REQUESTS"`
	LogLevel          string `default:"info" env:"LOG_LEVEL"`
}

func NewSettings() *ServerSettings {
	settings := &ServerSettings{}
	val := reflect.ValueOf(settings).Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		if !fieldValue.CanSet() {
			continue
		}

		defaultTag := field.Tag.Get("default")

		if defaultTag != "" {
			switch field.Type.Kind() {
			case reflect.String:
				fieldValue.SetString(defaultTag)
			case reflect.Bool:
				if boolVal, err := strconv.ParseBool(defaultTag); err == nil {
					fieldValue.SetBool(boolVal)
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if intVal, err := strconv.ParseInt(defaultTag, 10, 64); err == nil {
					fieldValue.SetInt(intVal)
				}
			default:
				// Other kinds (chan, map, slice, etc.) are not supported as settings fields.
			}
		}

		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}

		switch field.Type.Kind() {
		case reflect.String:
			fieldValue.SetString(env.String(envTag, fieldValue.String()))
		case reflect.Bool:
			fieldValue.SetBool(env.Bool(envTag, fieldValue.Bool()))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
			fieldValue.SetInt(int64(env.Int(envTag, int(fieldValue.Int()))))
		case reflect.Int64:
			fieldValue.SetInt(env.Int64(envTag, fieldValue.Int()))
		default:
			// Other kinds (chan, map, slice, etc.) are not supported as settings fields.
		}
	}

	return settings
}
