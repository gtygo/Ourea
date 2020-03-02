package cache

import (
	"errors"
	"github.com/gtygo/Ourea/boltkv"
	"github.com/gtygo/Ourea/kv"
	cmap "github.com/orcaman/concurrent-map"
	"sync"
)

const maxChSize=100000

var ErrNotFound=errors.New("Not Found ")

var mu sync.Mutex

type Cache struct{
	obj sync.Pool

	m cmap.ConcurrentMap
	hm cmap.ConcurrentMap

	store kv.Item

}

func NewCache()*Cache{
	return &Cache{
		m:      cmap.New(),
		hm:     cmap.New(),
		obj:   sync.Pool{
			New: func() interface{} {
				return cmap.New()
			},
		},
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

func (c *Cache)GetAllKey()[]string{
	return c.m.Keys()
}

func (c *Cache)Hset(hashName string,key string,v string){
	mu.Lock()
	defer mu.Unlock()
	sigMap:=c.obj.Get().(cmap.ConcurrentMap)
	defer c.obj.Put(sigMap)

	sigMap.Set(key,v)
	c.hm.Set(hashName,sigMap)
	return
}

func (c *Cache)Hget(hashName string,key string)(string,error){
	mu.Lock()
	defer mu.Unlock()
	sigMap:=c.obj.Get().(cmap.ConcurrentMap)
	defer c.obj.Put(sigMap)

	hmap,ok:=c.hm.Get(hashName)
	if !ok{
		return "",ErrNotFound
	}
	sigMap,ok=hmap.(cmap.ConcurrentMap)
	if !ok{
		panic("type assert concurrentMap failed!")
	}

	v,ok:=sigMap.Get(key)
	if !ok{
		return "",ErrNotFound
	}
	return v.(string),nil
}

func (c *Cache)Dump()error{
	db,err:=boltkv.Open("dump.rdb")
	if err!=nil{
		return err
	}
	c.store=db

	items:=c.m.Items()
	if err:=c.store.Set(items);err!=nil{
		return err
	}
	c.store.Close()
	return nil
}

func (c *Cache)Cover()error{
	db,err:=boltkv.Open("dump.rdb")
	if err!=nil{
		return err
	}
	c.store=db

	coverM,err:=c.store.Get()
	if err!=nil{
		return err
	}
	c.m.MSet(coverM)
	db.Close()
	return nil
}


