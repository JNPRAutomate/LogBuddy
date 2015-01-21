package logbuddy

import (
	"net"
	"net/url"
	"os"
	"strings"
	"testing"
)

func TestFileStorage(t *testing.T) {
	var err error
	var fileURL *url.URL
	//Create test directory
	err = os.Mkdir("_test", 0777)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	wd, _ := os.Getwd()
	fileURL, err = url.Parse(strings.Join([]string{"file://", wd, "/_test/test.txt"}, ""))
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
	var i int
	for i = 0; i < 100; i++ {
		err = fs.Write(LogMessage{Network: "udp4", SrcIP: net.ParseIP("1.1.1.1"), DstIP: net.ParseIP("2.2.2.2"), DstPort: 5000, SrcPort: 5000, Message: []byte("Eat the test")})
		if err != nil {
			t.Fatalf("%s", err.Error())
		}
	}
	t.Log("Write to file")
	fs.Close()
	//Remove test files
	err = os.RemoveAll("./_test")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
}
