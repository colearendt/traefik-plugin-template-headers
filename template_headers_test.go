package traefik_plugin_template_headers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	t.Run("check that header can be set", func(t *testing.T) {
		config := &Config{
			TemplateHeaders: []TemplateHeader{
				genConfig("X-RSC-Request", "{{ .Scheme }}://{{ .Host }}/rsc{{ .RequestURI }}"),
				genConfig("X-Exact-Req", "{{ .Scheme }}://{{ .Host }}{{ .RequestURI }}"),
			},
		}

		req1 := "http://localhost:80/"
		reqHeader1 := testHttpRequest(t, config, "X-RSC-Request", req1)
		if reqHeader1 != "http://localhost:80/rsc/" {
			t.Errorf("Request Header did not match expected (%v): %v\n", "http", reqHeader1)
		}

		req2 := "http://localhost:80/something?hello=1"
		exp2 := "http://localhost:80/rsc/something?hello=1"
		reqHeader2 := testHttpRequest(t, config, "X-RSC-Request", req2)
		if reqHeader2 != exp2 {
			t.Errorf("Request Header did not match expected (%v): %v\n", exp2, reqHeader2)
		}

		req3 := "https://some.example.com/some-path/level?query=true"
		reqHeader3 := testHttpRequest(t, config, "X-Exact-Req", req3)
		if reqHeader3 != req3 {
			t.Errorf("Request Header did not match expected (%v): %v\n", req3, reqHeader3)
		}
	})
}

func genConfig(name string, template string) TemplateHeader {
	return TemplateHeader{
		Header:   name,
		Template: template,
	}
}

func testHttpRequest(t *testing.T, config *Config, header string, request string) string {
	next := func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Add(header, req.Header.Get(header))
		rw.WriteHeader(http.StatusOK)
	}

	templateHeaders, err := New(context.Background(), http.HandlerFunc(next), config, "templateHeaders")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, request, nil)

	templateHeaders.ServeHTTP(recorder, req)

	return recorder.Result().Header.Get(header)
}
