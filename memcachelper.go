package memcachelper

import (
	"log"
	"strconv"

	"github.com/bradfitz/gomemcache/memcache"
)

type CacheMeta struct {
	Client   *memcache.Client
	Key      string
	DataType string
	Data     interface{}
	Refresh  bool
}

func (cm *CacheMeta) convertFn() []byte {
	// TODO add more types
	var b []byte
	switch cm.DataType {
	case "int":
		d := strconv.Itoa(cm.Data.(int))
		b = []byte(string(d))
		break
	case "string":
		b = []byte(cm.Data.(string))
		break
	}
	return b
}

func (cm *CacheMeta) revertFn(b []byte) interface{} {
	// TODO add more types
	switch cm.DataType {
	case "int":
		d, err := strconv.Atoi(string(b))
		if err != nil {
			log.Println(err)
		}
		return d
	case "string":
		d := string(b)
		return d
	}
	return nil
}

func Get(meta CacheMeta) interface{} {
	exists := retrieve(meta)
	if exists == nil || meta.Refresh {
		store(meta)
	} else {
		d := meta.revertFn(exists)
		meta.Data = d
	}
	return meta.Data
}

func remove(meta CacheMeta) bool {
	err := meta.Client.Delete(meta.Key)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func store(meta CacheMeta) bool {
	err := meta.Client.Set(&memcache.Item{Key: meta.Key, Value: meta.convertFn()})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func retrieve(meta CacheMeta) []byte {
	it, err := meta.Client.Get(meta.Key)
	if it == nil && err != nil {
		return nil
	} else {
		return it.Value
	}
}
