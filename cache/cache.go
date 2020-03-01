package cache

import (
	"errors"
	cmap "github.com/orcaman/concurrent-map"

)

const maxChSize=100000

var ErrNotFound=errors.New("Not Found ")

type Cache struct{
	m cmap.ConcurrentMap

	keyCh chan string
}

func NewCache()*Cache{
	return &Cache{
		m:      cmap.New(),
		keyCh:  make(chan string,maxChSize),
	}
}

func (c *Cache)Set(k string,v string){
	c.m.Set(k,v)
}

func (c *Cache)Get(k string)(string,error){
	v,ok:=c.m.Get(k)
	if !ok{
		return "",ErrNotFound
	}
	return v.(string),nil
}

func (c *Cache)Del(k string){
	c.m.Remove(k)
}

