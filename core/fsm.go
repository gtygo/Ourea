package core

import (
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	hraft "github.com/hashicorp/raft"
)

type fsm struct {
	db     DB
	logger *log.Logger
}

type snapShot struct {
	db     DB
	logger *log.Logger
}

//NewFsm return  a fsm instance
func NewFsm(path string) (*fsm, error) {
	db, err := NewDB(path, path)
	if err != nil {
		return nil, err
	}
	return &fsm{
		db:     db,
		logger: log.New(os.Stderr, "[fsm] ", log.LstdFlags),
	}, nil
}

//Get return value
func (f *fsm) Get(key string) (string, error) {
	v, err := f.db.get([]byte(key))
	if err != nil {
		f.logger.Printf("get key %s error: %s", key, err)
		return "", err
	}
	return string(v), nil
}

func (f *fsm) Set(key, value string) error {
	err := f.db.set([]byte(key), []byte(value))
	if err != nil {
		f.logger.Fatalf("set key: %s value: %s error: %s ", key, value, err)
		return err
	}
	return nil
}

func (f *fsm) Delete(key string) error {
	isSuccess, err := f.db.delete([]byte(key))
	if err != nil || !isSuccess {
		f.logger.Fatalf("delete key: %s error: %s ", key, err)
		return err
	}

	return nil
}

func (f *fsm) Apply(log *hraft.Log) interface{} {
	var c command
	if err := json.Unmarshal(log.Data, &c); err != nil {
		panic("failed to unmarshal raft log")
	}
	switch strings.ToLower(c.Op) {
	case "set":
		return f.Set(c.Key, c.Value)
	case "delete":
		return f.Delete(c.Key)
	default:
		panic("command type not support")
	}
}

func (f *fsm) Snapshot() (hraft.FSMSnapshot, error) {
	f.logger.Printf("Generate FSMSnapshot")
	return &snapShot{
		db:     f.db,
		logger: log.New(os.Stderr, "[fsmSnapshot] ", log.LstdFlags),
	}, nil
}

func (f *fsm) Restore(readClose io.ReadCloser) error {
	f.logger.Printf("Restore snapshot from FSMSnapshot")
	defer readClose.Close()

	var (
		readBuf  []byte
		protoBuf *proto.Buffer
		err      error
		keyCount int = 0
	)

	f.logger.Printf("Read all data")

	if readBuf, err = ioutil.ReadAll(readClose); err != nil {
		f.logger.Printf("new protoBuf length %d bytes", len(protoBuf.Bytes()))
		return err
	}
	protoBuf = proto.NewBuffer(readBuf)

	f.logger.Printf("new protoBuf length %d bytes", len(protoBuf.Bytes()))

	// decode messages from 1M block file
	// the last message could decode failed with io.ErrUnexpectedEOF
	for {
		item := &ProtoKVItem{}
		if err = protoBuf.DecodeMessage(item); err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			f.logger.Printf("DecodeMessage failed %v", err)
			return err
		}
		// apply item to store
		f.logger.Printf("Set key %v to %v count: %d", item.Key, item.Value, keyCount)
		err = f.db.set(item.Key, item.Value)
		if err != nil {
			f.logger.Printf("Snapshot load failed %v", err)
			return err
		}
		keyCount = keyCount + 1
	}

	f.logger.Printf("Restore total %d keys", keyCount)

	return nil

}

func (f snapShot) Persist(sink hraft.SnapshotSink) error {
	f.logger.Printf("Persist action in fsmSnapshot")
	defer sink.Close()

	ch := f.db.snapShotItems()

	keyCount := 0

	// read kv item from channel
	for {
		buff := proto.NewBuffer([]byte{})

		dataItem := <-ch
		item := dataItem.(*KVItem)

		if item.IsFinished() {
			break
		}

		// create new protobuf item
		protoKVItem := &ProtoKVItem{
			Key:   item.key,
			Value: item.val,
		}

		keyCount = keyCount + 1

		// encode message
		buff.EncodeMessage(protoKVItem)

		if _, err := sink.Write(buff.Bytes()); err != nil {
			return err
		}
	}
	f.logger.Printf("Persist total %d keys", keyCount)

	return nil
}
func (f *snapShot) Release() {
	f.logger.Printf("Release action in fsmSnapshot")
}
