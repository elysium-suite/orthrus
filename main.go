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

func sendFile(platform string, writer http.ResponseWriter) {
	file := "aeacus-" + platform + ".zip"

	openFile, err := os.Open(file)
	defer openFile.Close()
	if err != nil {
		http.Error(writer, "File "+file+" not found", 404)
		return
	}

	fileHeader := make([]byte, 512)
	openFile.Read(fileHeader)
	contentType := http.DetectContentType(fileHeader)

	fileStat, _ := openFile.Stat()
	fileSize := strconv.FormatInt(fileStat.Size(), 10)

	writer.Header().Set("Content-Disposition", "attachment; filename="+file)
	writer.Header().Set("Content-Type", contentType)
	writer.Header().Set("Content-Length", fileSize)

	openFile.Seek(0, 0)
	io.Copy(writer, openFile)
}

func endpoint(writer http.ResponseWriter, request *http.Request) {
	platform := request.URL.Query().Get("os")
	password := request.URL.Query().Get("pass")

	pass, err := ioutil.ReadFile("pass.txt")
	if err != nil {
		http.Error(writer, "Error: server could not read pass.txt.", 400)
		return
	}

	if password != string(pass) {
		http.Error(writer, "Error: Incorrect password.", 400)
		return
	}

	if !validateString(platform) {
		http.Error(writer, "Error: DTA attempted. This incident will be reported.", 400)
		return
	}

	switch platform {
	case "win32":
		sendFile("win32", writer)
	case "linux":
		sendFile("linux", writer)
	default:
		http.Error(writer, "Error: invalid OS", 400)
	}

	return
}
