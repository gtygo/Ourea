package kv

type Item interface {
	Get() (map[string]interface{}, error)
	Set(map[string]interface{}) error
	Delete(k []byte) error
	Close()
}
