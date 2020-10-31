package oututil

import (
	"bytes"
	"fmt"
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
	fmt.Println(b.String())
	// Output:
	// >  PUT / HTTP/1.1
	// >  Host: www.example.org
	// >  User-Agent: Go-http-client/1.1
	// >  Content-Length: 75
	// >  Accept-Encoding: gzip
	// >
	// >  Go is a general-purpose language designed with systems programming in mind.
}
