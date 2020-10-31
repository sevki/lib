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
	const body = "Any seemingly pointless activity which is actually necessary to solve a problem which solves a problem which, several levels of recursion later, solves the real problem you're working on."
	req, err := http.NewRequest("PUT", "https://bovineshave.club", strings.NewReader(body))
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
	// >  Host: bovineshave.club
	// >  User-Agent: Go-http-client/1.1
	// >  Content-Length: 187
	// >  Accept-Encoding: gzip
	// > 
	// >  Any seemingly pointless activity which is actually necessary to solve a problem which solves a problem which, several levels of recursion later, solves the real problem you're working on.
}
