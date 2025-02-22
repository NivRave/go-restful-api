package cache

import "net/http"

type Writer struct {
	writer   http.ResponseWriter
	response response
	resource string
}

// Header implements http.ResponseWriter.
func (w *Writer) Header() http.Header {
	return w.response.header
}

// Write implements http.ResponseWriter.
func (w *Writer) Write(b []byte) (int, error) {
	w.response.body = make([]byte, len(b))
	for k, v := range b {
		w.response.body[k] = v
	}
	copyHeader(w.Header(), w.writer.Header())
	set(w.resource, &w.response)
	return w.writer.Write(b)
}

// WriteHeader implements http.ResponseWriter.
func (w *Writer) WriteHeader(statusCode int) {
	copyHeader(w.response.header, w.writer.Header())
	w.response.code = statusCode
	w.writer.WriteHeader(statusCode)
}

func NewWriter(w http.ResponseWriter, r *http.Request) *Writer {
	return &Writer{
		writer:   w,
		resource: MakeResource(r),
		response: response{
			header: http.Header{},
		},
	}
}
