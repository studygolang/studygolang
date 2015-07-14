package service

import (
	"encoding/json"
	gio "io"

	"config"
	"github.com/qiniu/api.v6/conf"
	"github.com/qiniu/api.v6/io"
	"github.com/qiniu/api.v6/rs"
	"logger"
)

var uptoken string

func InitQiniu() {
	conf.ACCESS_KEY = config.Config["qiniu_access_key"]
	conf.SECRET_KEY = config.Config["qiniu_secret_key"]

	putPolicy := rs.PutPolicy{
		Scope: config.Config["qiniu_bucket_name"],
		// CallbackUrl:  callbackUrl,
		// CallbackBody: callbackBody,
		// ReturnUrl:    returnUrl,
		// ReturnBody:   returnBody,
		// AsyncOps:     asyncOps,
		// EndUser:      endUser,
		// Expires:      expires,
	}

	uptoken = putPolicy.Token(nil)
}

func uploadLocalFile(localFile, key string) (err error) {
	InitQiniu()
	var ret io.PutRet
	var extra = &io.PutExtra{
	// Params:   params,
	// MimeType: mieType,
	// Crc32:    crc32,
	// CheckCrc: CheckCrc,
	}

	// ret       变量用于存取返回的信息，详情见 io.PutRet
	// uptoken   为业务服务器生成的上传口令
	// key       为文件存储的标识(文件名)
	// localFile 为本地文件名
	// extra     为上传文件的额外信息，详情见 io.PutExtra，可选
	err = io.PutFile(nil, &ret, uptoken, key, localFile, extra)

	if err != nil {
		//上传产生错误
		logger.Errorln("io.PutFile failed:", err)
		return
	}

	//上传成功，处理返回值
	logger.Debugln(ret.Hash, ret.Key)

	return
}

func UploadMemoryFile(r gio.Reader, key string) (err error) {
	InitQiniu()
	var ret io.PutRet
	var extra = &io.PutExtra{
	// Params:   params,
	// MimeType: mieType,
	// Crc32:    crc32,
	// CheckCrc: CheckCrc,
	}

	// ret       变量用于存取返回的信息，详情见 io.PutRet
	// uptoken   为业务服务器端生成的上传口令
	// key       为文件存储的标识
	// r         为io.Reader类型，用于从其读取数据
	// extra     为上传文件的额外信息,可为空， 详情见 io.PutExtra, 可选
	err = io.Put(nil, &ret, uptoken, key, r, extra)

	// 上传产生错误
	if err != nil {
		logger.Errorln("io.Put failed:", err)

		errInfo := make(map[string]interface{})
		err = json.Unmarshal([]byte(err.Error()), &errInfo)
		if err != nil {
			logger.Errorln("io.Put Unmarshal failed:", err)
			return
		}

		code, ok := errInfo["code"]
		if ok && code == 614 {
			err = nil
		}

		return
	}

	// 上传成功，处理返回值
	logger.Debugln(ret.Hash, ret.Key)

	return
}
