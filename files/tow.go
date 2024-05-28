package files

import (
	"fmt"
	"io"
	"log"
	"os"
	"raktabeeja/helper"
	"strings"
)

const bytesPerRead = 10000
const bytesReadForConversion = 12000

type FileBuffer struct {
	buf  []byte
	from int
}

func ExportFile(path string) {
	split_path := strings.Split(path, "/")
	file_name := helper.Xor(split_path[len(split_path)-1], 5)
	if err := os.Mkdir("./storage/"+file_name, os.ModePerm); err != nil {
		fmt.Println(err)
	}
	buf := make(chan []byte)
	file_save := 0
	go buffer_reader(buf, path, bytesPerRead)
	for byt := range buf {
		file_save++
		multi_hash, err := helper.SplitBytes(byt)
		if err != nil {
			panic(err.Error())
		}

		hex_file := helper.Base64Encode([]byte{byte(file_save)})
		if err := os.Mkdir("./storage/"+file_name+"/"+hex_file, os.ModePerm); err != nil {
			fmt.Println(err)
		}
		for _, v := range multi_hash {
			if file, err := os.Create("./storage/" + file_name + "/" + hex_file + "/" + helper.HashBytes(v)); err != nil {
				fmt.Println(err)
			} else {
				file.Write(v)
			}
		}
	}
}

func buffer_reader(buf chan []byte, path string, per_buf int) {
	f, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}

	defer f.Close()

	buffer := make([]byte, per_buf)
	for {
		n, err := f.Read(buffer)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		if err == io.EOF {
			close(buf)
			break
		}
		buf <- buffer[:n]
	}
}

func readFilesInDir(finish chan FileBuffer, hash string, file_chunk int) {
	base_file_chunk := helper.Base64Encode([]byte{byte(file_chunk)})
	files, err := os.ReadDir("./storage/" + hash + "/" + base_file_chunk)
	if err != nil {
		return
	}
	go readFilesInDir(finish, hash, file_chunk+1)

	buf := make(chan []byte)
	buf_1 := make(chan []byte)
	go buffer_reader(buf, "./storage/"+hash+"/"+base_file_chunk+"/"+files[0].Name(), bytesReadForConversion)
	go buffer_reader(buf_1, "./storage/"+hash+"/"+base_file_chunk+"/"+files[1].Name(), bytesReadForConversion)

	finish <- FileBuffer{buf: helper.JoinBytes(<-buf, <-buf_1), from: file_chunk}
}

func ReadByt(hash string) {
	files, err := os.ReadDir("./storage/" + hash)
	if err != nil {
		panic("No File")
	}

	actual_chunk := make(chan FileBuffer, len(files)-1)

	go readFilesInDir(actual_chunk, hash, 1)

	file, err := os.OpenFile(helper.XorBack(hash, 5), os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	for i := 0; i < len(files); i++ {
		ac_d := (<-actual_chunk)
		_, err = file.Seek((int64(ac_d.from)-1)*bytesPerRead, 0)
		if err != nil {
			fmt.Println("Error seeking in file:", err)
			return
		}
		_, err := file.Write(ac_d.buf)
		if err != nil {
			fmt.Println("Error writing array1 to file:", err)
			return
		}
	}
}
