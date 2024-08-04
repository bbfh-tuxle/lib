package field

// The size of the MessageList header size in bytes
const LIST_HEADER_SIZE int64 = 8

// A file that message list will read from.
//
// An *os.File will always satisfy this interface.
type File interface {
	Read([]byte) (int, error)
	ReadAt([]byte, int64) (int, error)
}
