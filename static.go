package main

import (
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func sendFile(filePath string, w http.ResponseWriter) ([]byte, error) {

	mimeTypes := map[string]string{
		"jsx": "text/jsx",
	}

	ext := filepath.Ext(filePath)
	fileExt := ext[1:len(ext)]
	bytes, err := ioutil.ReadFile("public/" + filePath)
	if err != nil {
		return nil, errors.New("Error while reading static file.")
	}
	w.Header().Add("Content-Type", mimeTypes[fileExt])
	w.Write(bytes)
	return bytes, err
}
