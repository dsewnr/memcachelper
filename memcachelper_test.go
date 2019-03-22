package memcachelper

import (
	"fmt"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
)

func Test_GetFromCache(t *testing.T) {
	meta := CacheMeta{
		Client:   memcache.New("127.0.0.1:11211"),
		Key:      fmt.Sprintf("%s", "testFn"),
		DataType: "float64",
		Data:     3.1415927410125732,
		Refresh:  false,
	}
	ret := Get(meta)
	t.Log("cache hit", ret)
	meta.Refresh = true
	ret = Get(meta)
	t.Log("cache miss", ret)
}
