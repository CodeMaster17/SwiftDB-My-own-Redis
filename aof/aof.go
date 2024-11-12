package aof

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"swiftdb/resp"
	"sync"
	"time"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}

	// Start a goroutine to sync AOF to disk every 1 second
	go func() {
		for {
			aof.mu.Lock()

			aof.file.Sync()

			aof.mu.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil
	/*
		- What happens here is that we first create the file if it doesn’t exist or open it if it does.
		- Then, we create the bufio.Reader to read from the file.
		- We start a goroutine to sync the AOF file to disk every 1 second while the server is running.

	*/

	/*
		- The idea of syncing every second ensures that the changes we made are always present on disk. Without the sync, it would be up to the OS to decide when to flush the file to disk. With this approach, we ensure that the data is always available even in case of a crash. If we lose any data, it would only be within the second of the crash, which is an acceptable rate.
		- If you want 100% durability, we won’t need the goroutine. Instead, we would sync the file every time a command is executed. However, this would result in poor performance for write operations because IO operations are expensive.
	*/

}

// The next method is Close, which ensures that the file is properly closed when the server shuts down.
func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}

// After that, we create the Write method, which will be used to write the command to the AOF file whenever we receive a request from the client.

func (aof *Aof) Write(value resp.Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	data := value.Marshal()
	fmt.Println("Writing data to AOF:", string(data)) // Log the data being written

	_, err := aof.file.Write(data)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}

	return nil
}

// reading the file
func (aof *Aof) Read(fn func(value resp.Value)) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	aof.file.Seek(0, io.SeekStart)

	reader := resp.NewResp(aof.file)

	for {
		value, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		fn(value)
	}

	return nil
}
