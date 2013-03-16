package util

import (
	"bytes"
	"logger"
	"strconv"
)

// 内嵌bytes.Buffer，支持连写
type Buffer struct {
	*bytes.Buffer
}

func NewBuffer() *Buffer {
	return &Buffer{Buffer: new(bytes.Buffer)}
}

func (this *Buffer) Append(s string) *Buffer {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorln("*****内存不够了！******")
		}
	}()
	this.Buffer.WriteString(s)
	return this
}

func (this *Buffer) AppendInt(i int) *Buffer {
	return this.Append(strconv.Itoa(i))
}
