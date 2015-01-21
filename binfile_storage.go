package logbuddy

import (
	"bufio"
	"encoding/gob"
	"net/url"
	"os"
)

//BinFileStorage Store logs in a binary format
type BinFileStorage struct {
	Location *url.URL        //location of the storage file
	MsgChan  chan LogMessage //Channel to recieve messages from
	LogFile  *os.File        //File the file being used
}

//NewBinFileStorage Create an initialized NewBinFileStorage
func NewBinFileStorage(loc *url.URL) *BinFileStorage {
	return &BinFileStorage{Location: loc}
}

//Write Write data to the destination file
func (s *BinFileStorage) Write(data ...LogMessage) error {
	w := bufio.NewWriter(s.LogFile)
	e := gob.NewEncoder(w)
	for msg := range data {
		err := e.Encode(data[msg])
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
func (s *BinFileStorage) Read(data ...[]byte) error {
	return nil
}

//Open Open the file for writing
func (s *BinFileStorage) Open() error {
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
func (s *BinFileStorage) Close() error {
	err := s.LogFile.Close()
	return err
}

//SetDest Set the destination URL for writing
func (s *BinFileStorage) SetDest(dest string) error {
	return nil
}
