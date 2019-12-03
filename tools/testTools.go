package tools

import "math/rand"

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
