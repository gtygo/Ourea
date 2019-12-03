package core

import (
	"github.com/gtygo/Ourea/tools"
	"testing"
)

type res struct {
	ans string
}

var (
	keyRange = 10
	caseNum  = 100000
	dbPath   = "../tools/test/"
)

func TestFsm_GetAndSet(t *testing.T) {
	fsm, _ := NewFsm(dbPath)

	table, keyList := generateTestTableAndKeyList(caseNum, keyRange)

	for _, x := range keyList {
		if x == "" {
			continue
		}
		err := fsm.Set(x, table[x].ans)
		if err != nil {
			t.Fatal("set error:", err)
		}
	}

	for _, x := range keyList {
		if x == "" {
			continue
		}
		val, err := fsm.Get(x)
		if err != nil {
			t.Log("get error:", err)
		}
		if val != table[x].ans {
			t.Fatalf("got error ,except output: %s ,obtain output: %s ", table[x].ans, val)
		}
	}
}

func generateTestTableAndKeyList(thor, keyRange int) (map[string]res, []string) {
	table := make(map[string]res, thor)
	keyList := make([]string, thor)

	for i := 0; i < thor; i++ {
		k := tools.RandGenerateStr(keyRange)
		if k == "" || k == " " {
			k = tools.RandGenerateStr(keyRange)
		}
		if _, ok := table[k]; !ok {
			val := tools.RandGenerateStr(keyRange * 10)
			table[k] = res{
				ans: val,
			}
			keyList = append(keyList, k)
		}
	}
	return table, keyList
}
