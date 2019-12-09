package core

import (
	"math/rand"
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

func _TestFsm_GetAndSet(t *testing.T) {
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
		k := RandGenerateStr(keyRange)
		if k == "" || k == " " {
			k = RandGenerateStr(keyRange)
		}
		if _, ok := table[k]; !ok {
			val := RandGenerateStr(keyRange * 10)
			table[k] = res{
				ans: val,
			}
			keyList = append(keyList, k)
		}
	}
	return table, keyList
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ{}:~!@#$%^&*()_+|:?></")

func RandGenerateStr(rang int) string {
	n := rand.Intn(rang)
	n++
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
