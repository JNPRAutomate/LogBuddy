package logbuddy

import (
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestFileStorage(t *testing.T) {
	var err error
	var fileURL *url.URL
	wd, _ := os.Getwd()
	fileURL, err = url.Parse(strings.Join([]string{"file://", wd, "/test/test.txt"}, ""))
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	fs := &FileStorage{Location: fileURL}
	err = fs.Open()
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	t.Log("File Opened")
	//Write to file
	err = fs.Write([]byte("Log Msg to Write to File"))
	if err != nil {
		t.Fatalf("%s", err)
	}
	fs.Close()
}
