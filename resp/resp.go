package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// it is file that contains all the code related to serializing and deserializing on the buffer

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// struct used in serialization and deserialization process

type Value struct {
	Typ   string  // used to deefine the data type carried by value
	Str   string  // holds the value of string recevied from simple strings
	Num   int     // holds value of integers recevied from value
	Bulk  string  // used to hold the strings received from Bulk strings
	Array []Value // holds all the values received from the arrays
}

// func NewValue(Typ, Str string, num int, Bulk string, Array []Value) Value {
// 	return Value{
// 		yp:   Typ,
// 		Str:   Str,
// 		num:   num,
// 		Bulk:  Bulk,
// 		Array: Array,
// 	}
// }

// The Reader
type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// readLine reads the line from the buffer.
// readInteger reads the integer from the buffer.

/*In this function, we read one byte at a time until we reach ‘\r’, which indicates the end of the line. Then, we return the line without the last 2 bytes, which are ‘\r\n’, and the number of bytes in the line.*/
func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}

		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil

}

// method that will read from the buffer recursively
func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()

	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}

// read Array
func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.Typ = "array"

	// reading length of Array
	length, _, err := r.readInteger()

	if err != nil {
		return v, err
	}

	// foreach line, parse and read the value
	v.Array = make([]Value, length)
	for i := 0; i < length; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		// add parsed value to Array
		v.Array[i] = val
	}
	return v, nil
}

// read Bulk string
func (r *Resp) readBulk() (Value, error) {
	v := Value{}

	v.Typ = "Bulk"
	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	Bulk := make([]byte, len)

	r.reader.Read(Bulk)
	v.Bulk = string(Bulk)

	// Read the trailing CRLF
	/*
		Note that we call r.readLine() after reading the string to read the ‘\r\n’ that follows each Bulk string. If we don’t do this, the pointer will be left at ‘\r’ and the Read method won’t be able to read the next Bulk string correctly.
	*/
	r.readLine()

	return v, nil
}

// -- Writing resp

func (v Value) Marshal() []byte {
	switch v.Typ {
	case "Array":
		return v.marshalArray()
	case "Bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshallNull()
	case "error":
		return v.marshallError()
	default:
		return []byte{}
	}
}

// for simple strings
func (v Value) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// for Bulk string
func (v Value) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.Bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.Bulk...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

// for Array
func (v Value) marshalArray() []byte {
	len := len(v.Array)
	var bytes []byte
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < len; i++ {
		bytes = append(bytes, v.Array[i].Marshal()...)
	}

	return bytes
}

func (v Value) marshallError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.Str...)
	bytes = append(bytes, '\r', '\n')

	return bytes
}

func (v Value) marshallNull() []byte {
	return []byte("$-1\r\n")
}

// writer

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w}
}

// method that takes Value and writes the bytes it gets from the Marshal method to the Writer.

func (w *Writer) Write(v Value) error {
	var bytes = v.Marshal()

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}
	return nil
}
