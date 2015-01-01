package logbuddy

import (
	"net/url"
)

//FileStorage Allows data to be written to a file
type FileStorage struct {
	Location *url.URL     //location of the storage file
	MsgChan  chan Message //Channel to recieve messages from
}

//NewFileStorage Create an initialized NewFileStorage
func NewFileStorage(loc *url.URL) *FileStorage {
	return &FileStorage{Location: loc}
}

//Write Write data to the destination file
func (s *FileStorage) Write(data ...[]byte) error {
	return nil
}

//Read Reads data from the
func (s *FileStorage) Read(data ...[]byte) error {
	return nil
}

//Open Open the file for writing
func (s *FileStorage) Open() error {
	return nil
}

//Close Close the file for writing
func (s *FileStorage) Close() error {
	return nil
}

//SetDest Set the destination URL for writing
func (s *FileStorage) SetDest(dest string) error {
	return nil
}
