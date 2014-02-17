// Copyright 2014, Chandra Sekar S.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the README.md file.

package main

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"io"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

const keyChars = "abcdefghijklmnopqrstuvwxyz0123456789"

var (
	errNoSlide       = errors.New("Archive must contain at least one .slide file")
	errTooManySlides = errors.New("Archive must contain only one .slide file")
)

func processUpload(w http.ResponseWriter, r *http.Request) {
	slideId, viewId := getIdPair(r.FormValue("existingId"))
	if slideId == "" || viewId == "" {
		slideId, viewId = generateKey(), generateKey()
	}

	file, _, err := r.FormFile("slideArchive")
	if err == http.ErrMissingFile {
		w.WriteHeader(http.StatusBadRequest)
		msgTmpl.Execute(w, map[string]interface{}{
			"title": "Missing File",
			"msg":   "No slide archive file was chosen",
			"error": true,
		})
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer file.Close()

	if err := extractArchive(filepath.Join(*slidesDir, slideId), file); err != nil {

		badFile := err == gzip.ErrChecksum || err == gzip.ErrHeader || err == tar.ErrHeader
		badContent := err == errNoSlide || err == errTooManySlides

		if badFile || badContent {
			w.WriteHeader(http.StatusBadRequest)
			msgTmpl.Execute(w, map[string]interface{}{
				"title": "Bad Upload",
				"msg":   err.Error(),
				"error": true,
			})
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	addSlide(slideId, viewId)
	saveIndex(filepath.Join(*slidesDir, "index.json"))

	shareTmpl.Execute(w, map[string]string{
		"slideId": slideId,
		"viewId":  viewId,
		"baseURL": *baseURL,
	})

	return
}

func generateKey() string {
	buf := make([]byte, 8)
	for i := range buf {
		buf[i] = keyChars[rand.Intn(len(keyChars))]
	}

	return string(buf)
}

func extractArchive(slideBase string, file multipart.File) (err error) {
	if err := os.RemoveAll(slideBase); err != nil {
		return err
	}

	if err := os.MkdirAll(slideBase, 0700); err != nil {
		return err
	}

	defer cleanupOnFailure(slideBase, &err)

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzipReader.Close()

	var slideCount int
	reader := tar.NewReader(gzipReader)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		fileName := filepath.Join(slideBase, header.Name)
		if header.FileInfo().IsDir() {
			if err := os.MkdirAll(fileName, 0700); err != nil {
				return err
			}

			continue
		}

		if strings.HasSuffix(fileName, ".slide") {
			fileName = filepath.Join(filepath.Dir(fileName), "main.slide")
			slideCount++
		}

		f, err := os.Create(fileName)
		if err != nil {
			return err
		}

		defer f.Close()

		if _, err := io.Copy(f, reader); err != nil {
			return err
		}
	}

	if slideCount == 0 {
		return errNoSlide
	}

	if slideCount > 1 {
		return errTooManySlides
	}

	return nil
}

func cleanupOnFailure(slideBaseDir string, failure *error) {
	if *failure == nil {
		return
	}

	if err := os.RemoveAll(slideBaseDir); err != nil {
		log.Printf("Failed to cleanup directory: %s : %s\n", slideBaseDir, err)
	}
}

var uploadTmpl = template.Must(template.New("upload").Parse(`
<!doctype html>
<html>
<head>
	<title>Remote Presentations</title>
	<style>
	body {
		color: rgb(51, 51, 51);
	}
	</style>
</head>
<body>
<h1>Add/Update Presentation</h1>
<form action="/" method="POST" enctype="multipart/form-data">
	<label for="slideArchive">Slide Archive (.tar.gz or .tgz):</label>
	<input type="file" id="slideArchive" name="slideArchive">
	<p>
	<label for="existingId">Existing Slide ID:</label>
	<input type="text" id="existingId" name="existingId">
	<p>
	<input type="submit" value="Upload">
	<input type="reset" value="Reset">
</form>
</body>
</html>`))

var msgTmpl = template.Must(template.New("message").Parse(`
<!doctype html>
<html>
<head>
	<title>{{.title}}</title>
	<style>
	body {
		color: rgb(51, 51, 51);
	}
	</style>
</head>
<body>
<h1{{if .error}} style="color: red"{{end}}>{{.msg}}</h1>
</body>
</html>`))

var shareTmpl = template.Must(template.New("share").Parse(`
<!doctype html>
<html>
<head>
	<title>Slide Uploaded</title>
	<style>
	body {
		color: rgb(51, 51, 51);
	}
	</style>
</head>
<body>
<h1>Successfully Uploaded Slide</h1>
	This information is required for updating and presenting the slide.
	<p>
	<label for="slideId">Slide ID:</label>
	<input type="text" id="slideId" readonly="readonly" value="{{.slideId}}">
	<p>
	<label for="presenterURL">Presenter URL:</label>
	<input type="text" id="presenterURL" readonly="readonly" value="{{.baseURL}}/{{.slideId}}">
	<p>
	<label for="presenterURL">View URL:</label>
	<input type="text" id="presenterURL" readonly="readonly" value="{{.baseURL}}/{{.viewId}}">
	<script>
		document.getElementById("presenterURL").focus();
	</script>
</body>
</html>`))
