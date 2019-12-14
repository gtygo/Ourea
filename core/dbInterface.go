package core

type DataItem interface{}

type DB interface {
	get(key []byte) ([]byte, error)
	set(key, value []byte) error
	delete(key []byte) (bool, error)
	snapShotItems() <-chan DataItem
	close()
}
