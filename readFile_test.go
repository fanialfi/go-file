package main

import (
	"fmt"
	"testing"
)

func BenchmarkChanReadFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dataChan := make(chan []byte, 2)
		errCh := make(chan error, 2)

		go ChanReadFile(path, dataChan, errCh)

		for {
			select {
			case _, ok := <-dataChan:
				if !ok {
					dataChan = nil
				} else {
					// fmt.Printf("read %d bytes : %s\n\n", len(data), string(data))
				}
			case err, ok := <-errCh:
				if ok && err != nil {
					fmt.Println("Error :", err.Error())
				} else {
					errCh = nil
				}
			}

			// hentikan for lop jika kedua channel sudah ditutup
			if dataChan == nil && errCh == nil {
				break
			}
		}

		// fmt.Println("\nread file complete")
	}
}

func BenchmarkReadFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ReadFile()
	}
}
