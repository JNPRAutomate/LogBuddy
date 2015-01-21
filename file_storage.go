package logbuddy

import (
	"bufio"
	"net/url"
	"os"
	"strings"
)

//FileStorage Allows data to be written to a file
type FileStorage struct {
	Location *url.URL        //location of the storage file
	MsgChan  chan LogMessage //Channel to recieve messages from
	LogFile  *os.File        //File the file being used
}

//NewFileStorage Create an initialized NewFileStorage
func NewFileStorage(loc *url.URL) *FileStorage {
	return &FileStorage{Location: loc}
}

//Write Write data to the destination file
func (s *FileStorage) Write(data ...LogMessage) error {
	w := bufio.NewWriter(s.LogFile)

	for msg := range data {
		_, err := w.WriteString(strings.Join([]string{data[msg].String(), string('\n')}, ""))
		if err != nil {
			return err
		}
	}

	err := w.Flush()

	if err != nil {
		return err
	}
	return nil
}

//Read Reads data from the
func (s *FileStorage) Read(data ...[]byte) error {
	return nil
}

//Open Open the file for writing
func (s *FileStorage) Open() error {
	var err error
	if _, err = os.Stat(s.Location.Path); os.IsNotExist(err) {
		s.LogFile, err = os.Create(s.Location.Path)
		if err != nil {
			return err
		}
	} else {
		s.LogFile, err = os.OpenFile(s.Location.Path, os.O_RDWR|os.O_RDWR, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

//Close Close the file for writing
func (s *FileStorage) Close() error {
	err := s.LogFile.Close()
	return err
}

//SetDest Set the destination URL for writing
func (s *FileStorage) SetDest(dest string) error {
	return nil
}
