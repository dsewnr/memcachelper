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

func (cm *CacheMeta) ConvertFn() []byte {
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

func (cm *CacheMeta) RevertFn(b []byte) interface{} {
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

func GetFromCache(meta CacheMeta) interface{} {
	exists := get(meta)
	if exists == nil || meta.Refresh {
		set(meta)
	} else {
		d := meta.RevertFn(exists)
		meta.Data = d
	}
	return meta.Data
}

func del(meta CacheMeta) bool {
	err := meta.Client.Delete(meta.Key)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func set(meta CacheMeta) bool {
	err := meta.Client.Set(&memcache.Item{Key: meta.Key, Value: meta.ConvertFn()})
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func get(meta CacheMeta) []byte {
	it, err := meta.Client.Get(meta.Key)
	if it == nil && err != nil {
		return nil
	} else {
		return it.Value
	}
}
