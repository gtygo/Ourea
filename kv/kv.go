package kv

type Item interface {
	Get(k []byte) ([]byte, error)
	Set(k []byte, v []byte) error
	Delete(k []byte) error
	Close()
}