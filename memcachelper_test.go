package memcachelper

import (
	"fmt"
	"testing"

	"github.com/bradfitz/gomemcache/memcache"
)

func testFn(i int) string {
	return fmt.Sprintf("%d", i)
}

func Test_GetFromCache(t *testing.T) {
	meta := CacheMeta{
		Client:   memcache.New("127.0.0.1:11211"),
		Key:      fmt.Sprintf("%s", "testFn"),
		DataType: "string",
		Data:     testFn(0),
		Refresh:  false,
	}
	ret := GetFromCache(meta)
	t.Log("cache hit", ret)
	meta.Refresh = true
	ret = GetFromCache(meta)
	t.Log("cache miss", ret)
}
