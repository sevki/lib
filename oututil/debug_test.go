package oututil

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

func init() {
	debugging = true
}

func ExamplePrefixLine() {
	const body = "Go is a general-purpose language designed with systems programming in mind."
	req, err := http.NewRequest("PUT", "http://www.example.org", strings.NewReader(body))
	if err != nil {
		log.Fatal(err)
	}

	dump, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		panic(err)
	}
	b := bytes.NewBuffer(dump)
	PrefixLine(b, "> ", 1)
}
