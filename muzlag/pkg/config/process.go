package config

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/kelseyhightower/envconfig"
)

var ErrProcessForNotPointer = errors.New("process not possible for non-pointers")

func Process(appName string, cfg any) error {
	if reflect.ValueOf(cfg).Kind() != reflect.Ptr {
		return ErrProcessForNotPointer
	}

	err := envconfig.Process(appName, cfg)
	if err != nil {
		return fmt.Errorf("process config: %w", err)
	}

	return nil
}
