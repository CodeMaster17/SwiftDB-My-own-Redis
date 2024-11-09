package resp

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
input format
$5\r\nAhmed\r\n

*/

func main() {

	input := "$5\r\nAhmed\r\n"

	reader := bufio.NewReader(strings.NewReader(input))

	// Read the first line
	b, _ := reader.ReadByte()

	if b != '$' {
		fmt.Println("Invalid type, expecting bulk strings only")
		os.Exit(1)
	}

	// reading the number
	size, _ := reader.ReadByte()

	strSize, _ := strconv.ParseInt(string(size), 10, 64)

	// consume /r/n
	reader.ReadByte()
	reader.ReadByte()

	/*
		By doing this, we have read the byte that determines the data type, followed by the number that indicates the number of characters in the string. Then, we read an additional 2 bytes to get rid of the ‘\r\n’ that follows the number.
		Now, our reader object is positioned at the 5th byte, which is the letter ‘A’.

		$
		5
		\r
		\n
		A
		h
		m
		e
		d
		\r
		\n

		Since we know the number of characters or bytes we need to read, we just need to create a buffer array and read it.


	*/

	name := make([]byte, strSize)
	reader.Read(name)
	fmt.Println(string(name))
}
