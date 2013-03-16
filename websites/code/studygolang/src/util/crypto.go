package util

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"time"
)

func Md5(text string) string {
	hashMd5 := md5.New()
	io.WriteString(hashMd5, text)
	return fmt.Sprintf("%x", hashMd5.Sum(nil))
}

// 产生唯一的id
func GenUUID() string {
	buf := make([]byte, 16)
	io.ReadFull(rand.Reader, buf)
	return fmt.Sprintf("%x%x", buf, time.Now().UnixNano())
}
