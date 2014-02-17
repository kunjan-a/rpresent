// Copyright 2014, Chandra Sekar S.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the README.md file.

package main

import (
	"encoding/json"
	"os"
)
import "sync"

var (
	indexMutex sync.RWMutex
	slides     = make(map[string]string)
	views      = make(map[string]string)
)

func loadIndex(file string) error {
	f, err := os.Open(file)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return nil
	}

	defer f.Close()

	if err := json.NewDecoder(f).Decode(&slides); err != nil {
		return err
	}

	for k, v := range slides {
		views[v] = k
	}
	return nil
}

func addSlide(slideId, viewId string) {
	indexMutex.Lock()
	defer indexMutex.Unlock()

	slides[slideId] = viewId
	views[viewId] = slideId
}

func saveIndex(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()

	return json.NewEncoder(f).Encode(slides)
}

func getSlideId(slideOrViewId string) string {
	indexMutex.RLock()
	defer indexMutex.RUnlock()

	if slideId, ok := views[slideOrViewId]; ok {
		return slideId
	}

	if _, ok := slides[slideOrViewId]; ok {
		return slideOrViewId
	}

	return ""
}

func getIdPair(slideIdParam string) (slideId, viewId string) {
	indexMutex.RLock()
	defer indexMutex.RUnlock()

	if viewId, ok := slides[slideIdParam]; ok {
		return slideIdParam, viewId
	}
	return "", ""
}
