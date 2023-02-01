package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

const (
	fileName = "/tmp/foo.json"
)

func main() {
	var mode string
	flag.StringVar(&mode, "mode", "", "must be 'reader' or 'writer'")
	flag.Parse()

	var err error
	switch mode {
	case "writer":
		err = writer()
	case "reader":
		err = reader()
	default:
		log.Fatalf("invalid execution mode %q", mode)
	}
	if err != nil {
		log.Fatalf("execution error: %v", err)
	}
}

func writer() error {
	log.Printf("creating file")
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("error closing file in writer: %v", err)
		}
	}()

	msg := []byte("hello")
	n, err := file.Write(msg)
	if err != nil {
		return err
	}
	if n != len(msg) {
		return fmt.Errorf("expected to write %d bytes but only wrote %d", len(msg), n)
	}
	return nil
}

func reader() error {
	log.Printf("opening file with backoff")
	var file *os.File
	var err error
	for i := 0; i < 10; i++ {
		file, err = os.Open(fileName)
		if err == nil {
			break
		}

		log.Printf("error opening file: %v", err)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("error opening file in reader: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("error closing file in reader: %v", err)
		}
	}()

	res, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	log.Printf("read from file: %s", res)

	return nil
}
