package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

var path = "/tmp/golang/tes.txt"
var alphaNumSet = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz012345679")
var randomizer = rand.New(rand.NewSource(time.Now().UnixNano()))

// buat random string sesuai panjangnya
func randomString(length int) string {
	str := make([]byte, length)

	for i := range str {
		str[i] = alphaNumSet[randomizer.Intn(len(alphaNumSet))]
	}

	return string(str)
}

// function buat deteksi error
func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return err != nil
}

func CreateFile() {
	// deteksi apakah file sudah ada
	// os.Stat mengembalikan dua data, informasi tentang path yang dicari dan error jika ada
	_, err := os.Stat(path)

	// buat file jika belum ada
	if os.IsNotExist(err) {
		// os.Create mengembalikan 2 data yang pertama object bertipe *os.File yang baru dibuat dengan statusnya adalah otomatis open
		var file, err = os.Create(path)
		if isError(err) {
			return
		}

		// setelah operasi dengan file selesai file harus di close
		defer file.Close()
	}

	fmt.Println("file berhasil di buat", path)
}

func WriteFile() {
	// buka file dengan level akses write
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if isError(err) {
		return
	}
	defer file.Close()

	data := randomString(10 << 20)

	_, err = file.WriteString(data)
	if isError(err) {
		return
	}

	// simpan perubahan file
	err = file.Sync()
	if isError(err) {
		return
	}

	fmt.Println("file berhasil di isi")
}

// file yang mau dibaca harus dibuka terlebih dahulu dengan level akses minimal read
func ReadFile() {
	// membuka file
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if isError(err) {
		return
	}
	defer file.Close()

	// melihat info file buat mencari tahu file size
	fileinfo, err := file.Stat()
	if isError(err) {
		return
	}

	// baca file
	text := make([]byte, fileinfo.Size())
	for {
		n, err := file.Read(text)
		if err != io.EOF {
			if isError(err) {
				break
			}
		}
		if n == 0 {
			break
		}
	}

	// fmt.Printf("file berhasil di baca\n\n")
	// fmt.Println(string(text))
}

// membaca file dengan method channel
func ChanReadFile(path string, dataCh chan<- []byte, errCh chan<- error) {
	// membuka file
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if isError(err) {
		errCh <- err
		close(errCh)
		return
	}
	defer file.Close()

	// melihat info file buat mencari tahu file size
	fileinfo, err := file.Stat()
	if isError(err) {
		errCh <- err
		return
	}

	// membaca file
	for {
		buffer := make([]byte, fileinfo.Size())
		n, err := file.Read(buffer)
		if err != nil {
			if err != io.EOF {
				errCh <- err
			}
			break
		}

		if n > 0 {
			dataCh <- buffer[:n]
		}
	}
	close(errCh)
	close(dataCh)
}

// delete file
func DeleteFile() {
	err := os.Remove(path)
	if isError(err) {
		return
	}

	fmt.Println("file berhasil di hapus")
}

// jalankan function chanReadFile
func main() {
	dataChan := make(chan []byte)
	errCh := make(chan error)

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

	fmt.Println("\nread file complete")
}
