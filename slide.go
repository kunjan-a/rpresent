// Copyright 2014, Chandra Sekar S.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the README.md file.

package main

import (
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"code.google.com/p/go.tools/present"
)

func presentSlide(w http.ResponseWriter, r *http.Request) {
	slideIdParam := strings.SplitN(r.URL.Path, "/", 2)[1]
	slideId := index.getSlideId(slideIdParam)
	if slideId == "" {
		http.NotFound(w, r)
		return
	}

	slideFile := filepath.Join(*slidesDir, slideId, "main.slide")
	f, err := os.Open(slideFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	doc, err := present.Parse(f, slideFile, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := present.Template()
	tmpl.Funcs(template.FuncMap{"playable": func(c present.Code) bool { return false },
		"rSlideId": func() string {
			return slideIdParam
		},
		"userRole": func() string {
			if slideId == slideIdParam {
				return "p"
			} else {
				return "v"
			}
		}})
	_, err = tmpl.New("action").Parse(actionTmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = tmpl.New("slides").Parse(slidesTmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := doc.Render(w, tmpl); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Imported from https://code.google.com/p/go/source/browse/present/templates/slides.tmpl?repo=talks
const slidesTmpl = `
{/* This is the slide template. It defines how presentations are formatted. */}

{{define "root"}}
<!DOCTYPE html>
<html>
  <head>
    <title>{{.Title}}</title>
    <meta charset='utf-8'>
	<script>
	var rSlideId="{{rSlideId}}";
	var userRole="{{userRole}}";
	</script>
    <script src='/static/slides.js'></script>
    <script src='/static/remote.js'></script>
  </head>

  <body style='display: none'>

    <section class='slides layout-widescreen'>

      <article>
	  	{{if eq userRole "v"}}<a class="helpLink" href="help" target="_blank">Help</a>{{end}}
        <h1>{{.Title}}</h1>
        {{with .Subtitle}}<h3>{{.}}</h3>{{end}}
        {{if not .Time.IsZero}}<h3>{{.Time.Format "2 January 2006"}}</h3>{{end}}
        {{range .Authors}}
          <div class="presenter">
            {{range .TextElem}}{{elem $.Template .}}{{end}}
          </div>
        {{end}}
      </article>

  {{range $i, $s := .Sections}}
  <!-- start of slide {{$s.Number}} -->
      <article>
	  {{if eq userRole "v"}}<a class="helpLink" href="help" target="_blank">Help</a>{{end}}
      {{if $s.Elem}}
        <h3>{{$s.Title}}</h3>
        {{range $s.Elem}}{{elem $.Template .}}{{end}}
      {{else}}
        <h2>{{$s.Title}}</h2>
      {{end}}
      </article>
  <!-- end of slide {{$i}} -->
  {{end}}{{/* of Slide block */}}

      <article>
	  	{{if eq userRole "v"}}<a class="helpLink" href="help" target="_blank">Help</a>{{end}}
        <h3>Thank you</h1>
        {{range .Authors}}
          <div class="presenter">
            {{range .Elem}}{{elem $.Template .}}{{end}}
          </div>
        {{end}}
      </article>

  </body>
  {{if .PlayEnabled}}
  <script src='/play.js'></script>
  {{end}}
</html>
{{end}}

{{define "newline"}}
<br>
{{end}}`

// Imported from https://code.google.com/p/go/source/browse/present/templates/action.tmpl?repo=talks
const actionTmpl = `
{/*
This is the action template.
It determines how the formatting actions are rendered.
*/}

{{define "section"}}
  <h{{len .Number}} id="TOC_{{.FormattedNumber}}">{{.FormattedNumber}} {{.Title}}</h{{len .Number}}>
  {{range .Elem}}{{elem $.Template .}}{{end}}
{{end}}

{{define "list"}}
  <ul>
  {{range .Bullet}}
    <li>{{style .}}</li>
  {{end}}
  </ul>
{{end}}

{{define "text"}}
  {{if .Pre}}
  <div class="code"><pre>{{range .Lines}}{{.}}{{end}}</pre></div>
  {{else}}
  <p>
    {{range $i, $l := .Lines}}{{if $i}}{{template "newline"}}
    {{end}}{{style $l}}{{end}}
  </p>
  {{end}}
{{end}}

{{define "code"}}
  <div class="code{{if playable .}} playground{{end}}" contenteditable="true" spellcheck="false">{{.Text}}</div>
{{end}}

{{define "image"}}
<div class="image">
  <img src="{{.URL}}"{{with .Height}} height="{{.}}"{{end}}{{with .Width}} width="{{.}}"{{end}}>
</div>
{{end}}

{{define "iframe"}}
<iframe src="{{.URL}}"{{with .Height}} height="{{.}}"{{end}}{{with .Width}} width="{{.}}"{{end}}></iframe>
{{end}}

{{define "link"}}<p class="link"><a href="{{.URL}}" target="_blank">{{style .Label}}</a></p>{{end}}

{{define "html"}}{{.HTML}}{{end}}`
