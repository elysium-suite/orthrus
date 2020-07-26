package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Starting endpoint...")
	http.HandleFunc("/", endpoint)
	err := http.ListenAndServe(":6969", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func endpoint(writer http.ResponseWriter, request *http.Request) {
	fileName := request.URL.Query().Get("filename")
	password := request.URL.Query().Get("password")

	pass, err := ioutil.ReadFile("pass.txt")
	if err != nil {
		fmt.Println("Could not read pass.txt")
	}

	if password != string(pass) {
		http.Error(writer, "Incorrect password!", 400)
		return
	}

	if fileName == "" {
		http.Error(writer, "Get 'filename' not specified in url.", 400)
		return
	}

	fmt.Println("Client requests: " + fileName)

	openFile, err := os.Open(fileName)
	defer openFile.Close()
	if err != nil {
		http.Error(writer, "File"+fileName+"not found", 404)
		return
	}

	fileHeader := make([]byte, 512)
	openFile.Read(fileHeader)
	contentType := http.DetectContentType(fileHeader)

	FileStat, _ := openFile.Stat()
	FileSize := strconv.FormatInt(FileStat.Size(), 10)

	writer.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	writer.Header().Set("Content-Type", contentType)
	writer.Header().Set("Content-Length", FileSize)

	openFile.Seek(0, 0)
	io.Copy(writer, openFile)
	return
}
