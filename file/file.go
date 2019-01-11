package file

import (
    "errors"
    "fmt"
    "os"
    "path/filepath"
    "time"
)


// 初始化日志文件, 不存在则新建
func initFile(f string) error {
	if !logFileExists(f) {
		f, err := openFile(f)
		if err != nil {
			return errors.New(fmt.Sprintf("log: cannot create log: %v\n", err.Error()))
		}
		f.Close()
	}
	return nil
}

// 获取日志文件大小
func getLogSize(file string) uint64 {
	f, _ := os.Stat(file)
	return uint64(f.Size())
}

func renameFile(now time.Time, file string) {
	newName := fmt.Sprintf("%04d%02d%02d_%02d%02d%02d.log",
		now.Year(),
		now.Month(),
		now.Day(),
		now.Hour(),
		now.Minute(),
		now.Second(),
	)
	os.Rename(file, filepath.Join(filepath.Dir(file), newName))
}

// 判断文件是否存在
func logFileExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}

// 创建日志文件
func openFile(file string) (*os.File, error) {
	ps := filepath.Dir(file)
	if e := os.MkdirAll(ps,0666);e == nil {
	}
	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return nil, fmt.Errorf("log: cannot create log: %v", err)
	}
	return f, nil
}
