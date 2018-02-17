// Copyright 2018 Sevki <s@sevki.org>. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http // import "sevki.org/x/http"

import (
	"mime"
	"net/http"
	"path"
	"strings"
)

// ContentTypeHandler attaches the correct Content-Type
// to http requests based on the file extension.
func ContentTypeHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		mimetype := "application/octet-stream"
		ext := path.Ext(r.URL.Path)
		if strings.HasSuffix(r.URL.Path, "/") {
			ext = ".html"
		}

		mt := mime.TypeByExtension(ext)
		if mt != "" {
			mimetype = mt
		}
		if ext == "" {
			mimetype = "text/html"
		}
		w.Header().Set("Content-Type", mimetype)
		h.ServeHTTP(w, r)
	})
}
