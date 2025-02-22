package handlers

import "net/http"

func RootHandler(writer http.ResponseWriter, reader *http.Request) {
	if reader.URL.Path != "/" {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("Asset not found\n"))
		return
	} else {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("Running API V1\n"))
	}
}
