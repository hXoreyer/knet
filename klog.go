package knet

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

type log struct {
	TimeStamp string `json:"timestamp"`
	Level     string `json:"level"`
	Content   string `json:"content"`
}

type Klog struct {
	path   string
	uptime string
}

func NewKlog(path string) *Klog {
	os.Mkdir(path, 0644)
	os.Chmod(path, 0644)
	kl := &Klog{
		path:   path,
		uptime: "00:00:00",
	}
	go kl.dateFile("success")
	go kl.dateFile("error")
	go kl.dateFile("info")
	return kl
}

func (l *Klog) Success(format string, v ...interface{}) {
	l.writeFile("success", fmt.Sprintf(format, v...))
}

func (l *Klog) Info(format string, v ...interface{}) {
	l.writeFile("info", fmt.Sprintf(format, v...))
}

func (l *Klog) Error(format string, v ...interface{}) {
	l.writeFile("error", fmt.Sprintf(format, v...))
}

func (l *Klog) writeFile(level, content string) {
	logger := &log{
		Level:     level,
		TimeStamp: time.Now().String(),
		Content:   content,
	}
	buffer, _ := json.Marshal(&logger)
	fl, err := os.OpenFile(path.Join(l.path, level+".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	buffer = append(buffer, '\r', '\n')
	defer fl.Close()
	fl.Write(buffer)
}

func (l *Klog) dateFile(level string) {
	ts := strings.Split(l.uptime, ":")
	var h, m, s int
	h, _ = strconv.Atoi(ts[0])
	m, _ = strconv.Atoi(ts[1])
	s, _ = strconv.Atoi(ts[2])
	ticker := time.NewTicker(time.Second)
	for {
		ch := <-ticker.C
		if ch.Hour() == h && ch.Minute() == m && ch.Second() == s {
			fl, err := os.Open(path.Join(l.path, level+".log"))
			if err != nil {
				return
			}
			f, err := os.OpenFile(path.Join(l.path, fmt.Sprintf("%d-%d-%d %s", ch.Year(), ch.Month(), ch.Day(), level)), os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return
			}
			io.Copy(f, fl)
			f.Close()
			fl.Close()
			os.Remove(path.Join(l.path, level+".log"))
		}
	}
}

//设置每日更新文件时间(h:m:s) 默认0点整
func (l *Klog) SetUpdateTime(t string) {
	l.uptime = t
}
