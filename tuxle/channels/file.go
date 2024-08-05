package channels

// A file that message list will read from.
//
// An *os.File will always satisfy this interface.
type File interface {
	Read([]byte) (int, error)
	ReadAt([]byte, int64) (int, error)
	WriteAt([]byte, int64) (int, error)
}
