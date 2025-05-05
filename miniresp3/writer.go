package miniresp3

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Writer struct {
	bw *bufio.Writer
	sb strings.Builder
}

func NewWriter(w io.Writer) *Writer {
	bw := bufio.NewWriter(w)
	bw.Reset(w)
	return &Writer{bw: bw, sb: strings.Builder{}}
}

func (w *Writer) Reset() {
	w.sb.Reset()
}

func (w *Writer) Write() error {
	if _, err := w.bw.Write([]byte(w.sb.String())); err != nil {
		return err
	}

	w.Reset()
	return w.bw.Flush()
}

func (w *Writer) AppendNil() {
	w.sb.WriteByte(byte(RESPNull))
	w.sb.WriteByte('\r')
	w.sb.WriteByte('\n')
}

func (w *Writer) AppendInt(num int) {
	w.sb.WriteByte(byte(RESPNumber))
	w.sb.WriteString(strconv.Itoa(num))
	w.sb.WriteByte('\r')
	w.sb.WriteByte('\n')
}

func (w *Writer) AppendFloat64(num float64) {
	w.sb.WriteByte(byte(RESPDoubles))
	w.sb.WriteString(fmt.Sprintf("%.2f", num))
	w.sb.WriteByte('\r')
	w.sb.WriteByte('\n')
}

func (w *Writer) AppendSimpleStr(str string) {
	w.sb.WriteByte(byte(RESPSimpleString))
	w.sb.WriteString(str)
	w.sb.WriteByte('\r')
	w.sb.WriteByte('\n')
}

func (w *Writer) AppendBulkStr(str string) {
	w.sb.WriteByte(byte(RESPBulkString))
	w.sb.WriteString(strconv.Itoa(len(str)))
	w.sb.WriteByte('\r')
	w.sb.WriteByte('\n')

	w.sb.WriteString(str)
	w.sb.WriteByte('\r')
	w.sb.WriteByte('\n')
}

func (w *Writer) AppendArrInt(arr []int) {
	w.sb.WriteByte(byte(RESPArray))
	w.sb.WriteString(strconv.Itoa(len(arr)))
	w.sb.WriteByte('\r')
	w.sb.WriteByte('\n')

	for _, val := range arr {
		w.AppendInt(val)
	}
}

func (w *Writer) AppendArrStr(arr []string) {
	w.sb.WriteByte(byte(RESPArray))
	w.sb.WriteString(strconv.Itoa(len(arr)))
	w.sb.WriteByte('\r')
	w.sb.WriteByte('\n')

	for _, val := range arr {
		w.AppendBulkStr(val)
	}
}

func (w *Writer) AppendArrAny(arr []interface{}) {
	w.sb.WriteByte(byte(RESPArray))
	w.sb.WriteString(strconv.Itoa(len(arr)))
	w.sb.WriteByte('\r')
	w.sb.WriteByte('\n')

	for _, val := range arr {
		switch val := val.(type) {
		case int:
			w.AppendInt(val)
		case string:
			w.AppendBulkStr(val)
		case []int:
			w.AppendArrInt(val)
		case []string:
			w.AppendArrStr(val)
		}
	}
}

func (w *Writer) AppendArrHeader(len int) {
	w.sb.WriteByte(byte(RESPArray))
	w.sb.WriteString(strconv.Itoa(len))
	w.sb.WriteByte('\r')
	w.sb.WriteByte('\n')
}

func (w *Writer) AppendSimpleError(errMsg string) {
	w.sb.WriteByte(byte(RESPSimpleError))
	w.sb.WriteString(errMsg)
	w.sb.WriteByte('\r')
	w.sb.WriteByte('\n')
}

func (w *Writer) AppendBulkErr(errMsg string) {
	w.sb.WriteByte(byte(RESPBulkError))
	w.sb.WriteString(strconv.Itoa(len(errMsg)))
	w.sb.WriteByte('\r')
	w.sb.WriteByte('\n')

	w.sb.WriteString(errMsg)
	w.sb.WriteByte('\r')
	w.sb.WriteByte('\n')
}

func (w *Writer) AppendMap(m map[string]interface{}) {
	w.sb.WriteByte(byte(RESPMap))
	w.sb.WriteString(strconv.Itoa(len(m)))
	w.sb.WriteByte('\r')
	w.sb.WriteByte('\n')

	for key, val := range m {
		w.AppendSimpleStr(key)
		switch val := val.(type) {
		case int:
			w.AppendInt(val)
		case string:
			w.AppendBulkStr(val)
		case []int:
			w.AppendArrInt(val)
		case []string:
			w.AppendArrStr(val)
		}
	}
}
