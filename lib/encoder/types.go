package encoder

import "bufio"

// IBinary interface for binary marshaller
type IBinary interface {
	Marshal(io *bufio.Writer) error
	Unmarshall(io *bufio.Reader) error
}
