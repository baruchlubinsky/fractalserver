package main

import (
	"net/http"
	"os"
	"io/ioutil"
)

func index(response http.ResponseWriter, request *http.Request) {
	returnFile(response, "index.html")
}

func init() {
	http.HandleFunc("/", index)
}

func returnFile(response http.ResponseWriter, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		response.WriteHeader(400)
		return
	}
	data, _ := ioutil.ReadAll(file)
	response.Write(data) 
}