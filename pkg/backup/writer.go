package backup

import (
	"io"

	"github.com/wault-pw/alice/lib/encoder"
	"google.golang.org/protobuf/proto"
)

type Writer struct {
	wr io.Writer
}

func NewWriter(wr io.Writer) *Writer {
	return &Writer{wr: wr}
}

func (w *Writer) write(mark byte, message proto.Message) (int, error) {
	data, err := proto.Marshal(message)
	if err != nil {
		return 0, err
	}

	size := encoder.BinaryUint32(uint32(len(data)))
	out := make([]byte, 0, 1+len(size)+len(data))
	out = append(out, mark)
	out = append(out, size...)
	out = append(out, data...)

	return w.wr.Write(out)
}
