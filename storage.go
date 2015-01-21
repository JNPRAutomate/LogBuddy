package logbuddy

//Storage an interface for all storage types
type Storage interface {
	Write(...LogMessage) error
	Read() error //Read data from the
	SetDest(string) error
	Open() error  //Opens the storage location for writing
	Close() error //Closes the storage location
}
