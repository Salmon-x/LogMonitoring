package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan string)
}

type LogProcess struct {
	rc chan []byte		// 读取管道
	wc chan string		// 写入管道
	read Reader
	write Writer
}

type ReadFromFile struct {
	path string		// 读取文件的路径
}

type WriteToInfluxDB struct {
	influxDbDsn string		// influxdb data source
}

// 读取文件
func (r *ReadFromFile)Read(rc chan []byte) {
	// 读取文件
	f, err := os.Open(r.path)
	if err != nil {
		panic(fmt.Sprintf("open file error:%s",err.Error()))
	}
	// 从文件末尾开始逐行读取文件内容
	f.Seek(0,2)
	rd := bufio.NewReader(f)
	for  {
		line,err := rd.ReadBytes('\n')
		if err == io.EOF {
			time.Sleep(500*time.Millisecond)
			continue
		}else if err != nil {
			panic(fmt.Sprintf("ReadBytes error:%s",err))
		}
		/*此处[:len(line)-1]是为了去掉每行末尾的换行符号*/
		rc <- line[:len(line)-1]
	}
}


// 写入管道
func (l *WriteToInfluxDB)Write(wc chan string)  {
	// 写入模块
	for v := range wc {
		fmt.Println(v)
	}
}


func (l *LogProcess)Process()  {
	// 解析模块
	for v := range l.rc{
		l.wc <- strings.ToUpper(string(v))
	}
}


func main()  {
	r :=  &ReadFromFile{
		path: "./access.log",
	}

	w := &WriteToInfluxDB{
		influxDbDsn: "asd",
	}

	lp := &LogProcess{
		rc: make(chan []byte),
		wc: make(chan string),
		read: r,
		write: w,
	}

	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)
	time.Sleep(60*time.Second)
}
