package store

type DataItem interface{}

type DB interface {
	Get(key []byte) ([]byte, error)
	Set(key, value []byte) error
	Delete(key []byte) (bool, error)
	SnapShotItems() <-chan DataItem
}
