package internal

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"xstorage/pkg/utils"
)

// 20byte cache flag 'cache2mp4.0.0.3'
// [
//	8byte package flag,
//  4byte package size,
//  12byte data header,
//  n byte data
// ]

const (
	cacheflag = "cache2mp4.0.0.3"
)

var lock sync.RWMutex

type Cache2Mp4 struct {
	file          *os.File
	deviceNo      string
	lastTimestamp uint64
}

var defCacheMng = make(map[string]*Cache2Mp4)

func CacheGet(deviceId, ss string) *Cache2Mp4 {
	lock.Lock()
	defer lock.Unlock()
	if v, ok := defCacheMng[ss]; ok {
		return v
	}
	cache := &Cache2Mp4{}
	cache.deviceNo = deviceId
	defCacheMng[ss] = cache
	return cache
}

func CacheClose(deviceId, ss string) {
	lock.Lock()
	defer lock.Unlock()
	v, ok := defCacheMng[ss]
	if !ok || v.file == nil {
		return
	}
	tstamp := int64(v.lastTimestamp / 1000 / 1000)
	dtStr := utils.Unix2Str(tstamp, "20060102 150405", 0)
	timeStr := dtStr[9:]
	oldname := v.file.Name()
	v.file.Close()
	delete(defCacheMng, ss)
	newname := fmt.Sprintf("%s_%s.cache", v.file.Name(), timeStr)
	os.Rename(oldname, newname)
}

func packageCache(channel, ctype uint16, timestamp uint64, length int) []byte {
	var data [24]byte
	copy(data[:8], "cache...")
	binary.LittleEndian.PutUint32(data[8:], uint32(length))
	binary.LittleEndian.PutUint16(data[12:], channel)
	binary.LittleEndian.PutUint16(data[14:], ctype)
	binary.LittleEndian.PutUint64(data[16:], timestamp)
	return data[0:]
}

func (c *Cache2Mp4) syncOpenFile(channel uint16, timestamp uint64) error {
	tstamp := int64(timestamp / 1000 / 1000)
	dtStr := utils.Unix2Str(tstamp, "20060102 150405", 0)
	dateStr := dtStr[:8]
	timeStr := dtStr[9:]
	fpName := fmt.Sprintf("%s/%s/CH%02d_%s_%s", dateStr, c.deviceNo, channel, dateStr, timeStr)
	dir := filepath.Dir(fpName)
	os.MkdirAll(dir, os.ModePerm)
	fp, err := os.OpenFile(fpName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	c.file = fp
	var flag [20]byte
	copy(flag[0:], cacheflag)
	c.file.Write(flag[:20])
	return nil
}

func (c *Cache2Mp4) SyncWrite(channel, ctype uint16, timestamp uint64, data []byte) error {
	if c.file == nil {
		if err := c.syncOpenFile(channel, timestamp); err != nil {
			return err
		}
	}
	res := packageCache(channel, ctype, timestamp, len(data))
	c.file.Write(res)
	c.file.Write(data)
	c.lastTimestamp = timestamp
	return nil
}
