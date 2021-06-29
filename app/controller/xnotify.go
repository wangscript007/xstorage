package controller

import (
	"fmt"
	"xstorage/internal"

	"github.com/wlgd/xproto"
)

type XNotify struct {
}

// NewAccessData 实例化对象
func NewXNotify() *XNotify {
	return new(XNotify)
}

// AccessHandler 设备接入
func (o *XNotify) AccessHandler(data string, register *xproto.LinkAccess) error {
	fmt.Printf("%s\n", data)
	if !register.OnLine {
		internal.CacheClose(register.DeviceId, register.Session)
	}
	return nil
}

// AVFrameHandler 音视频数据
func (o *XNotify) AVFrameHandler(deviceId, ss string, channel, ctype uint16, timestamp uint64, data []byte) {
	fmt.Printf("%s >> | channel %d type: %d timestamp: %v length %d\n", ss, channel, ctype, timestamp, len(data))
	internal.CacheGet(deviceId, ss).SyncWrite(channel, ctype, timestamp, data)
}

// RawFrameHandler 普通二进制数据
func (o *XNotify) RawFrameHandler(deviceId, ss string, channel, ctype uint16, timestamp uint64, data []byte) {

}
