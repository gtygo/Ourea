package boltkv

import (
	"github.com/boltdb/bolt"
	"reflect"
	"testing"
)

func TestBoltItem_Get(t *testing.T) {
	type fields struct {
		Db *bolt.DB
	}
	type args struct {
		k []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bki := &BoltItem{
				Db: tt.fields.Db,
			}
			got, err := bki.Get(tt.args.k)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBoltItem_Set(t *testing.T) {
	type fields struct {
		Db *bolt.DB
	}
	type args struct {
		k []byte
		v []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bki := &BoltItem{
				Db: tt.fields.Db,
			}
			if err := bki.Set(tt.args.k, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBoltItem_Delete(t *testing.T) {
	type fields struct {
		Db *bolt.DB
	}
	type args struct {
		k []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bki := &BoltItem{
				Db: tt.fields.Db,
			}
			if err := bki.Delete(tt.args.k); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}