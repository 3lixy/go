package common

import (
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"golang.org/x/net/context"
	"time"
)

//Upload 上传七牛
//c nova对象
//localFile需要上传的文件路径
//key重命名的文件名称
func Upload(localFile string, key string) (string, string, error) {
	//获取配置
	accessKey := GetConfig().String("qiniu::access_key")
	secretKey := GetConfig().String("qiniu::secret_key")
	bucket := GetConfig().String("qiniu::bucket")
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	// 可选配置
	putExtra := storage.PutExtra{}
	err := formUploader.PutFile(context.Background(), &ret, upToken, key, localFile, &putExtra)
	if err != nil {
		panic(err)
	}
	return ret.Key, ret.Hash, err
}

//Download 七牛下载
//c nova对象
//key 需要下载的文件名
func Download(key string) string {
	//获取配置
	accessKey := GetConfig().String("qiniu::access_key")
	secretKey := GetConfig().String("qiniu::secret_key")
	domain := GetConfig().String("qiniu::domain")
	expireTime, _ := GetConfig().Int64("qiniu::expire_time")
	deadline := time.Now().Add(time.Second * time.Duration(expireTime)).Unix()
	mac := qbox.NewMac(accessKey, secretKey)
	privateAccessURL := storage.MakePrivateURL(mac, domain, key, deadline)
	return privateAccessURL
}
