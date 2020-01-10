package main

import (
	"github.com/gtygo/Ourea/db"
)

func main() {
	newDB := db.NewDB()
	defer newDB.Item.Close()
	err := newDB.Item.Set([]byte("kkkkkk"), []byte("vvvvvv"))
	if err != nil {
		println(err)
	}

}
