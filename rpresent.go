// Copyright 2014, Chandra Sekar S.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the README.md file.

package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"html/template"
	"time"

	"code.google.com/p/go.net/websocket"
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

	if err := index.load(filepath.Join(*slidesDir, "index.json")); err != nil {
		log.Fatalln("Failed to load index:", err)
	}

	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/help", help)
	http.HandleFunc("/static/", statics)
	http.Handle("/p", websocket.Handler(handlePresenter))
	http.Handle("/v", websocket.Handler(handleViewer))

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

func help(w http.ResponseWriter, r *http.Request) {
	helpTmpl.Execute(w, nil)
}

func statics(w http.ResponseWriter, r *http.Request) {
	parts := strings.SplitN(r.URL.Path, "/", 3)

	switch parts[2] {
	case "slides.js":
		w.Header()["Content-Type"] = []string{"application/javascript"}
		w.Write([]byte(slidesJS))
	case "remote.js":
		w.Header()["Content-Type"] = []string{"application/javascript"}
		w.Write([]byte(remoteJS))
	case "styles.css":
		w.Header()["Content-Type"] = []string{"text/css"}
		w.Write([]byte(stylesCSS))
	default:
		http.NotFound(w, r)
	}
}

var helpTmpl = template.Must(template.New("share").Parse(`
<!doctype html>
<html>
<head>
	<title>Help - RPresent</title>
	<style>
	body {
		color: rgb(51, 51, 51);
	}
	</style>
</head>
<body>
<h1>RPresent</h1>
	RPresent is a remote presentation tool. The presenter will be able to remotely control the slide displayed ini your screen.
	<p>
	<h2>Keyboard Shortcuts</h2>
	<ul>
		<li>Arrows - Navigate slides</li>
		<li>Pause - Pause/Resume presenter's remote control. Can be useful to review slides during presentation without being dragged back by the presenter.</li>
	</ul>
</body>
</html>`))
