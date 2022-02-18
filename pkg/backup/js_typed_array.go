package backup

import (
	"io"
	"strconv"
	"strings"
)

const jsArraySep = ","

type JsTypedArray struct {
	wr io.Writer
}

func NewJsTypedArray(wr io.Writer) *JsTypedArray {
	return &JsTypedArray{wr: wr}
}

func (j *JsTypedArray) Write(data []byte) (int, error) {
	// last char will be a separator after join
	out := make([]string, len(data)+1)

	for ix := range data {
		out[ix] = strconv.Itoa(int(data[ix]))
	}

	serialized := strings.Join(out, jsArraySep)
	return j.wr.Write([]byte(serialized))
}
