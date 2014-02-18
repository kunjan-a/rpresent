// Copyright 2014, Chandra Sekar S.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the README.md file.

package main

import (
	"encoding/json"
	"os"
	"sync"
)

var index = &slideIndex{slides: make(map[string]string), views: make(map[string]string)}

type slideIndex struct {
	sync.RWMutex
	slides map[string]string
	views  map[string]string
}

func (idx *slideIndex) load(file string) error {
	f, err := os.Open(file)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return nil
	}

	defer f.Close()

	if err := json.NewDecoder(f).Decode(&idx.slides); err != nil {
		return err
	}

	for k, v := range idx.slides {
		idx.views[v] = k
	}
	return nil
}

func (idx *slideIndex) addSlide(slideId, viewId string) {
	idx.Lock()
	defer idx.Unlock()

	idx.slides[slideId] = viewId
	idx.views[viewId] = slideId
}

func (idx *slideIndex) save(file string) error {
	idx.Lock()
	defer idx.Unlock()

	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()

	return json.NewEncoder(f).Encode(idx.slides)
}

func (idx *slideIndex) getSlideId(slideOrViewId string) string {
	idx.RLock()
	defer idx.RUnlock()

	if slideId, ok := idx.views[slideOrViewId]; ok {
		return slideId
	}

	if _, ok := idx.slides[slideOrViewId]; ok {
		return slideOrViewId
	}

	return ""
}

func (idx *slideIndex) getIdPair(slideIdParam string) (slideId, viewId string) {
	idx.RLock()
	defer idx.RUnlock()

	if viewId, ok := idx.slides[slideIdParam]; ok {
		return slideIdParam, viewId
	}
	return "", ""
}
