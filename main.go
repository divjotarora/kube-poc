package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"log"
	"net"
)

const (
	networkAddr = "localhost:8080"
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
	return nil

	// var conn net.Conn
	// var err error
	// for i := 0; i < 10; i++ {
	// 	conn, err = net.Dial("tcp", networkAddr)
	// 	if err != nil {
	// 		log.Printf("connection establishment error: %v\n", err)
	// 		time.Sleep(1 * time.Second)
	// 		continue
	// 	}
	// 	break
	// }
	// if err != nil {
	// 	return err
	// }

	// defer func() {
	// 	err = conn.Close()
	// 	if err != nil {
	// 		log.Printf("connection close error: %v\n", err)
	// 	}
	// }()

	// log.Printf("writing")
	// msg := []byte("hello\n")
	// n, err := conn.Write(msg)
	// if err != nil {
	// 	return err
	// }
	// if n != len(msg) {
	// 	return fmt.Errorf("expected to write %d bytes but wrote %d\n", len(msg), n)
	// }
	// return nil
}

func reader() error {
	listener, err := net.Listen("tcp", networkAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		defer func() {
			err = conn.Close()
			if err != nil {
				log.Printf("connection close error: %v\n", err)
			}
		}()

		log.Printf("accepted conn, reading")
		msg, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}

		log.Printf("read msg from conn: %v\n", msg)
		if len(msg) > 0 {
			break
		}
	}

	return nil
}
