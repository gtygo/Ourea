package cache

import (
	"math/rand"
	"strconv"
	"testing"
)
type kvi struct{
	k string
	v string
}


func BenchmarkCache_Set(b *testing.B) {


	c:=NewCache()

	table:=make([]kvi,b.N)


	for i:=0;i<b.N;i++{
		table[i]=kvi{
			k:strconv.Itoa(i)+"kkkkk",
			v:strconv.Itoa(i)+"vvvvv",
		}
	}

	b.ResetTimer()

	for i:=0;i<b.N;i++{
		c.Set(table[i].k,table[i].v)
	}

}


func BenchmarkCache_Get(b *testing.B) {


	c:=NewCache()


	for i:=0;i<b.N;i++{
		c.Set(strconv.Itoa(i)+"kkkkk",strconv.Itoa(i)+"vvvvv")
	}

	b.ResetTimer()
	for i:=0;i<b.N;i++{
		c.Get(strconv.Itoa(i)+"kkkkk")
	}
}

func BenchmarkCache_Set2(b *testing.B) {
	c:=NewCache()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m := strconv.Itoa(rand.Intn(100))
			n := strconv.Itoa(rand.Intn(200))
			c.Set(m,n)
		}
	})
}

func BenchmarkCache_Get2(b *testing.B) {
	c:=NewCache()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			m := strconv.Itoa(rand.Intn(100))
			c.Get(m)
		}
	})
}