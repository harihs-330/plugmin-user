package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var errUnimplemented = errors.New("not yet implemented")

type Record struct {
	Key    string
	Value  any
	Expiry time.Duration
}

type Config struct {
	Address  string
	DB       string
	DBID     int
	Password string
	URL      string
}

type Storage interface {
	Load(cfg *Config) Storage
	Get(context.Context, string) (Record, error)
	Set(context.Context, Record) (bool, error)
}

func (e Record) Int() int {
	v, _ := e.Value.(int)
	return v
}
func (e Record) String() string {
	v, _ := e.Value.(string)
	return v
}
func (e Record) StringMap() map[string]string {
	v, _ := e.Value.(map[string]string)
	return v
}
func (e Record) Slice() []interface{} {
	v, _ := e.Value.([]interface{})
	return v
}
func (e Record) Unmarshal(obj interface{}) error {
	var (
		data []byte
		err  error
	)

	// Check the type of e.Value
	switch v := e.Value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		data, err = json.Marshal(e.Value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
	}

	// Unmarshal the data into the provided obj
	if err := json.Unmarshal(data, obj); err != nil {
		return fmt.Errorf("failed to unmarshal value: %w", err)
	}

	return nil
}

/*
Struct implements all the methods that are defined on the parent struct
It is useful to embed this struct in the struct that implements the parent struct
to act as a default method in case you want to implement only some of the methods and not all
*/
type UnImplemented struct{}

var _ Storage = (*UnImplemented)(nil)

func (unImpl UnImplemented) Load(*Config) Storage {
	return &unImpl
}

func (unImpl UnImplemented) Get(context.Context, string) (Record, error) {
	return Record{}, errUnimplemented
}

func (unImpl UnImplemented) Set(context.Context, Record) (bool, error) {
	return false, errUnimplemented
}

// *************************** //
