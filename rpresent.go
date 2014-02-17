// Copyright 2014, Chandra Sekar S.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the README.md file.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	httpAddr  = flag.String("http", ":8080", "HTTP address to listen on")
	slidesDir = flag.String("d", "slides", "Directory to store slides in")
	baseURL   = flag.String("b", "http://localhost:8080", "Base URL for slides")
)

func main() {
	flag.Parse()

	if err := os.MkdirAll(*slidesDir, 0700); err != nil {
		log.Fatalln("Failed to create slides directory:", err)
	}

	if err := loadIndex(filepath.Join(*slidesDir, "index.json")); err != nil {
		log.Fatalln("Failed to load index:", err)
	}

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/static/", statics)

	fmt.Println("Listening at", *httpAddr)
	http.ListenAndServe(*httpAddr, nil)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		if r.URL.Path == "/" {
			uploadTmpl.Execute(w, nil)
		} else {
			presentSlide(w, r)
		}

	case "POST":
		processUpload(w, r)

	default:
		http.Error(w, "405 method not allowed", http.StatusMethodNotAllowed)
	}
}

func statics(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(r.URL.Path, "/", 3)

	switch parts[2] {
	case "slides.js":
		w.Header()["Content-Type"] = []string{"application/javascript"}
		w.Write([]byte(slidesJS))
	case "styles.css":
		w.Header()["Content-Type"] = []string{"text/css"}
		w.Write([]byte(stylesCSS))
	default:
		http.NotFound(w, r)
	}
}
