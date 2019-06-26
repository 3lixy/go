package common

import (
	"crypto/md5"
	"encoding/json"
	"github.com/coocood/freecache"
)

//内存缓存单例
var cacheNew *freecache.Cache

//NewFreeCache 获取cache实例
func NewFreeCache() *freecache.Cache {
	if cacheNew == nil {
		cacheSize := 1 * 1024 * 1024 //预分配内存1024KB
		cacheNew = freecache.NewCache(cacheSize)
	}
	return cacheNew
}

func getKey(key string) []byte {
	hash := md5.New()
	hash.Write([]byte(key))
	return hash.Sum(nil)
}

//Set 添加缓存
func SetCache(key string, val interface{}, expireSeconds int) (err error) {
	keyMd5 := getKey(key)
	valJson, err := json.Marshal(val)
	if err != nil {
		return
	}
	err = NewFreeCache().Set(keyMd5, valJson, expireSeconds)
	if err != nil {
		return
	}
	return nil
}

//Get 获取缓存
func GetCache(key string) ([]byte, error) {
	keyMd5 := getKey(key)
	val, err := NewFreeCache().Get(keyMd5)
	return val, err
}
