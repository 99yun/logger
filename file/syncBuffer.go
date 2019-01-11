package file

import (
	"bufio"
	"bytes"
	"fmt"
    "io"
    "os"
    "path/filepath"
    "time"
)

const bufferSize = 256 * 1024
const dateFormatter = "2006/01/02 15:04:05"

type syncBuffer struct {
	*bufio.Writer
	file     *os.File
	fileName string
	logDir   string
	nbytes uint64 // 当前文件已写入的字节数
	lastDate time.Time
}

func (sb *syncBuffer) Sync() error {
	return sb.file.Sync()
}

func (sb *syncBuffer) Write(p []byte) (n int, err error) {
	if sb.nbytes+uint64(len(p)) >= MAX_LOGFILE_SIZE || sb.lastDate.Day() != time.Now().Day() {
		if err = sb.rotateFile(sb.lastDate); err != nil {
		    return
		}
		sb.lastDate = time.Now()
	}
	n, err = sb.Writer.Write(p)
	sb.nbytes += uint64(n)
	sb.Writer.Flush()
	sb.Sync()
	return
}

// 关闭旧文件并创建一个新文件
// 并将旧文件重新命名,然后放入log的old目录下
func (sb *syncBuffer) rotateFile(now time.Time) error {
	if sb.file != nil {
		sb.Flush()
		sb.file.Close()
		renameFile(now,sb.file.Name())
	}

	var err error
	sb.file, err = openFile(sb.file.Name())
	sb.nbytes = 0
	if err != nil {
		return err
	}

	// 重新分配一个缓冲
	sb.Writer = bufio.NewWriterSize(sb.file, bufferSize)
	// Write header.
	return sb.writeHeader(now)
}

// 打印日志文件头部
func (sb *syncBuffer) writeHeader(now time.Time) error {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Log file created at: %s\n", now.Format(dateFormatter))
	//fmt.Fprint(&buf, "Log line format: [IWEF]yyyy-mm-dd hh:mm:ss msg\n")
	n, err := sb.file.Write(buf.Bytes())
	sb.nbytes += uint64(n)
	return err
}


func (sb *syncBuffer) initFile(file string, now time.Time) error {
    var err error
	sb.fileName = file
	sb.logDir = filepath.Dir(file)
	err = initFile(file)
	if err != nil {
	    return err
    }

    // 1. 获取旧日志的创建时间
    // 2. 获取旧日志的大小
	sb.nbytes = getLogSize(sb.fileName)
	sb.file, err = openFile(sb.fileName)
	if err != nil {
		return err
	}
	sb.Writer = bufio.NewWriterSize(sb.file, bufferSize)

	if sb.nbytes == 0 { // 空文件, 则写入文件头
		return sb.writeHeader(now)
	}
	return nil
}


func NewFileWriter(file string) (io.Writer, error) {
    sbuf := &syncBuffer{lastDate: time.Now()}
    err := sbuf.initFile(file,time.Now())
    return sbuf,err
}


