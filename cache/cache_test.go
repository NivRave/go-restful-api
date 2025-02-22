package cache

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestSetAndGet(t *testing.T) {
	resp := &response{header: http.Header{"Content-Type": {"text/plain"}}, code: 200, body: []byte("Hello")}
	set("/test", resp)
	if got := get("/test"); got == nil || got.code != 200 || string(got.body) != "Hello" {
		t.Errorf("Expected response not found")
	}
}

func TestOverwriteSet(t *testing.T) {
	resp1 := &response{code: 200, body: []byte("First")}
	resp2 := &response{code: 201, body: []byte("Second")}
	set("/test", resp1)
	set("/test", resp2)
	if got := get("/test"); got == nil || got.code != 201 || string(got.body) != "Second" {
		t.Errorf("Expected updated response not found")
	}
}

func TestRemoveSet(t *testing.T) {
	resp := &response{code: 200, body: []byte("Hello")}
	set("/test", resp)
	set("/test", nil)
	if get("/test") != nil {
		t.Errorf("Expected nil response after removal")
	}
}

func TestGetNonExistent(t *testing.T) {
	if get("/notfound") != nil {
		t.Errorf("Expected nil response for missing resource")
	}
}

func TestConcurrentGet(t *testing.T) {
	resp := &response{code: 200, body: []byte("Concurrent")}
	set("/concurrent", resp)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			get("/concurrent")
		}()
	}
	wg.Wait()
}

func TestMakeResource(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://example.com/test/", nil)
	if got := MakeResource(r); got != "/test" {
		t.Errorf("Unexpected resource: %s", got)
	}
}

func TestMakeResourceNil(t *testing.T) {
	if MakeResource(nil) != "" {
		t.Errorf("Expected empty string for nil request")
	}
}

func TestCopyHeader(t *testing.T) {
	src := http.Header{"Content-Type": {"text/html"}}
	dst := http.Header{}
	copyHeader(src, dst)
	if dst.Get("Content-Type") != "text/html" {
		t.Errorf("Header not copied correctly")
	}
}

func TestClean(t *testing.T) {
	set("/one", &response{code: 200})
	set("/two", &response{code: 404})
	Clean()
	if get("/one") != nil || get("/two") != nil {
		t.Errorf("Expected empty cache after Clean")
	}
}

func TestDrop(t *testing.T) {
	set("/test", &response{code: 200})
	Drop("/test")
	if get("/test") != nil {
		t.Errorf("Expected nil response after Drop")
	}
}

func TestServeNilInputs(t *testing.T) {
	if Serve(nil, nil) {
		t.Errorf("Expected false for nil inputs")
	}
}

func TestServeNoCacheControl(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://example.com/test", nil)
	r.Header.Set("Cache-Control", "no-cache")
	w := httptest.NewRecorder()
	if Serve(w, r) {
		t.Errorf("Expected false when Cache-Control is no-cache")
	}
}

func TestServeNotCached(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://example.com/missing", nil)
	w := httptest.NewRecorder()
	if Serve(w, r) {
		t.Errorf("Expected false for missing resource")
	}
}

func TestServeCachedResponse(t *testing.T) {
	resp := &response{header: http.Header{"Content-Type": {"text/plain"}}, code: 200, body: []byte("Hello")}
	set("/test", resp)
	r, _ := http.NewRequest("GET", "http://example.com/test", nil)
	w := httptest.NewRecorder()
	if !Serve(w, r) {
		t.Errorf("Expected true for cached resource")
	}
	if w.Code != 200 || w.Body.String() != "Hello" {
		t.Errorf("Unexpected response: %d %s", w.Code, w.Body.String())
	}
}

func TestServeHeadRequest(t *testing.T) {
	resp := &response{header: http.Header{"Content-Type": {"text/plain"}}, code: 200, body: []byte("Hello")}
	set("/test", resp)
	r, _ := http.NewRequest("HEAD", "http://example.com/test", nil)
	w := httptest.NewRecorder()
	if !Serve(w, r) {
		t.Errorf("Expected true for HEAD request")
	}
	if w.Code != 200 || w.Body.Len() != 0 {
		t.Errorf("Unexpected body length for HEAD request: %d", w.Body.Len())
	}
}
