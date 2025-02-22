package cloud

type Config struct {
	Region    string
	SecretKey string
	AppKey    string
}

type Storage interface {
	Get(string) ([]byte, error)
	Put(string, []byte) (bool, error)
}

type Provider interface {
	Load(cfg *Config) Provider
	Storage() Storage
}

/*
Struct implements all the methods that are defined on the parent struct
It is useful to embed this struct in the struct that implements the parent struct
to act as a default method in case you want to implement only some of the methods and not all
*/
type UnImplemented struct{}

var _ Provider = (*UnImplemented)(nil)

func (unImpl *UnImplemented) Load(_ *Config) Provider {
	return unImpl
}

func (unImpl UnImplemented) Storage() Storage {
	return StorageUnimplemented{}
}

type StorageUnimplemented struct{}

func (str StorageUnimplemented) Get(string) ([]byte, error) {
	return []byte("not yet implemented"), nil
}

func (str StorageUnimplemented) Put(string, []byte) (bool, error) {
	return false, nil
}

// *************************************************************** //
