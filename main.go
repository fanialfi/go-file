package main

import (
	"fmt"
	"os"
)

var path = "/tmp/golang/tes.txt"

// function buat deteksi error
func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return err != nil
}

func createFile() {
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

func main() {
	createFile()
}
