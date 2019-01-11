package logger

import (
	"testing"
)

func TestLoggerPrint(t *testing.T){
	err := Init(true, "./ok.log" ,5)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		Infof("test msg is ok:%d",i)
	}
}
