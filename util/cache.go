package util

import (
	"time"

	"github.com/muesli/cache2go"
)

type Cache struct {
	*cache2go.CacheTable
}

type CacheData struct {
	Count int
	Data  []IpMap
}

type IpMap struct {
	IP   string
	Info InfoDetail
}

type InfoDetail struct {
	Port int
	User string
	Pwd  string
}

var cache *cache2go.CacheTable
var valkey = "installList"

func init() {
	cache = cache2go.Cache("installCache")
}

//"设置缓存"
func (c *Cache) Set(val CacheData) error {
	cache.Add(valkey, 5*time.Second, val)
	return nil
}

//"获取缓存"
func (c *Cache) Get() (CacheData, error) {
	var cachedata CacheData
	val, err := cache.Value(valkey)
	if err != nil {
		//fmt.Println("valkey错误")
		return cachedata, err
	}
	v := val.Data().(CacheData)
	return v, nil
}

//"Del删除缓存"
func (c *Cache) Del() error {
	_, err := cache.Delete(valkey)
	if err != nil {
		return err
	}
	cache.Flush()
	return nil
}
