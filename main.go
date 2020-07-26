package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
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

func validateString(input string) bool {
	validationString := `^[a-zA-Z0-9-_.]+$`

	inputValidation := regexp.MustCompile(validationString)
	return inputValidation.MatchString(input)
}

func endpoint(writer http.ResponseWriter, request *http.Request) {
	operatingSystem := "aeacus-" + request.URL.Query().Get("os") + ".zip"
	password := request.URL.Query().Get("pass")

	pass, err := ioutil.ReadFile("pass.txt")
	if err != nil {
		fmt.Println("Could not read pass.txt")
	}

	if password != string(pass) {
		http.Error(writer, "Incorrect password!", 400)
		return
	}

	if operatingSystem == "" {
		http.Error(writer, "Query 'os' not specified in url.", 400)
		return
	}

	if !validateString(operatingSystem) {
		http.Error(writer, "Please specify a clean file path", 400)
		return
	} else if operatingSystem != "win32" {
		http.Error(writer, "OS not supported", 400)
		return
	} else if operatingSystem != "linux" {
		http.Error(writer, "OS not supported", 400)
		return
	}

	fmt.Println("Client requests: " + operatingSystem)

	openFile, err := os.Open(operatingSystem)
	defer openFile.Close()
	if err != nil {
		http.Error(writer, "File "+operatingSystem+" not found", 404)
		return
	}

	fileHeader := make([]byte, 512)
	openFile.Read(fileHeader)
	contentType := http.DetectContentType(fileHeader)

	fileStat, _ := openFile.Stat()
	fileSize := strconv.FormatInt(fileStat.Size(), 10)

	writer.Header().Set("Content-Disposition", "attachment; filename="+operatingSystem)
	writer.Header().Set("Content-Type", contentType)
	writer.Header().Set("Content-Length", fileSize)

	openFile.Seek(0, 0)
	io.Copy(writer, openFile)
	return
}
