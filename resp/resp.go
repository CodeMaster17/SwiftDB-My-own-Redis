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
	typ   string  // used to deefine the data type carried by value
	str   string  // holds the value of string recevied from simple strings
	num   int     // holds value of integers recevied from value
	bulk  string  // used to hold the strings received from bulk strings
	array []Value // holds all the values received from the arrays
}

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

// read array
func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	// reading length of array
	length, _, err := r.readInteger()

	if err != nil {
		return v, err
	}

	// foreach line, parse and read the value
	v.array = make([]Value, length)
	for i := 0; i < length; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		// add parsed value to array
		v.array[i] = val
	}
	return v, nil
}

// read bulk string
func (r *Resp) readBulk() (Value, error) {
	v := Value{}

	v.typ = "bulk"
	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)
	v.bulk = string(bulk)

	// Read the trailing CRLF
	/*
		Note that we call r.readLine() after reading the string to read the ‘\r\n’ that follows each bulk string. If we don’t do this, the pointer will be left at ‘\r’ and the Read method won’t be able to read the next bulk string correctly.
	*/
	r.readLine()

	return v, nil
}
