package traefik_plugin_template_headers

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"text/template"
)

// TemplateHeader holds one advanced-headers configuration section
type TemplateHeader struct {
	Template string `json:"template,omitempty"`
	Header   string `json:"header,omitempty"`
}

type internalTemplateHeader struct {
	Header   string
	Template *template.Template
}

// Config holds the plugin configuration
type Config struct {
	TemplateHeaders []TemplateHeader `json:"template-headers,omitempty"`
}

func CreateConfig() *Config {
	return &Config{}
}

type templateHeaders struct {
	name            string
	next            http.Handler
	templateHeaders []internalTemplateHeader
}

func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	fmt.Printf("Starting with config: %v\n", config.TemplateHeaders)
	templates := make([]internalTemplateHeader, len(config.TemplateHeaders))

	for i, tmpl := range config.TemplateHeaders {
		tmpTmpl, err := template.New(fmt.Sprintf("template-%d", i)).Parse(tmpl.Template)
		if err != nil {
			return nil, fmt.Errorf("error parsing template %s: %w", tmpl.Template, err)
		}
		templates[i] = internalTemplateHeader{
			Header:   tmpl.Header,
			Template: tmpTmpl,
		}
	}

	return &templateHeaders{
		name:            name,
		next:            next,
		templateHeaders: templates,
	}, nil
}

type templateData struct {
	Path                string
	Scheme              string
	Host                string
	Method              string
	Proto               string
	Query               string
	RequestURI          string
	HttpXForwardedProto string
	HttpXForwardedHost  string
	HttpHost            string
}

func (r *templateHeaders) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	tmplData := templateData{
		Path:                req.URL.EscapedPath(),
		Scheme:              req.URL.Scheme,
		Host:                req.URL.Host,
		Method:              req.Method,
		Proto:               req.Proto,
		Query:               req.URL.RawQuery,
		RequestURI:          req.URL.RequestURI(),
		HttpXForwardedProto: req.Header.Get("X-Forwarded-Proto"),
		HttpXForwardedHost:  req.Header.Get("X-Forwarded-Host"),
		HttpHost:            req.Header.Get("Host"),
	}
	for i, tmpl := range r.templateHeaders {
		if len(tmpl.Header) > 0 {
			// TODO: check for whether header already exists?
			// overwrite? add?

			var val bytes.Buffer
			err := tmpl.Template.Execute(&val, tmplData)
			if err != nil {
				fmt.Printf("error executing template for header %d '%s': %s", i, tmpl.Header, err)
				continue
			}

			// fmt.Printf("running with value: %v\n", val.String())
			// fmt.Printf("running with template: %v\n", tmpl.Template.Name())
			fmt.Printf("running with template data: %v\n", tmplData.RequestURI)
			// fmt.Printf("running with template data: %v\n", tmplData.Scheme)
			// fmt.Printf("running with template data: %v\n", tmplData.Path)
			// fmt.Printf("running with template data: %v\n", tmplData.Host)
			req.Header.Add(tmpl.Header, val.String())
		}
	}

	r.next.ServeHTTP(rw, req)
}
