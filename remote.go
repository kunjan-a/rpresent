// Copyright 2014, Chandra Sekar S.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the README.md file.

package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"code.google.com/p/go.net/websocket"
)

var registry = &listenerRegistry{listeners: make(map[string][]*slideListener)}

func handlePresenter(conn *websocket.Conn) {
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	var slideIdParam string
	if err := websocket.Message.Receive(conn, &slideIdParam); err != nil {
		return
	}

	slideId, _ := index.getIdPair(slideIdParam)
	if slideId == "" {
		return
	}

	for {
		var slide string

		conn.SetDeadline(time.Now().Add(15 * time.Minute))
		if err := websocket.Message.Receive(conn, &slide); err != nil {
			return
		}

		curSlide, _ := strconv.Atoi(slide)
		registry.setSlide(slideId, curSlide)
	}
}

func handleViewer(conn *websocket.Conn) {
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))
	var viewId string
	if err := websocket.Message.Receive(conn, &viewId); err != nil {
		return
	}

	slideId := index.getSlideId(viewId)
	if slideId == "" {
		return
	}

	listener := &slideListener{ch: make(chan int)}
	registry.addListener(slideId, listener)

	for {
		slide := listener.get(1 * time.Minute)
		if slide != 0 {
			conn.SetDeadline(time.Now().Add(10 * time.Second))
			if err := websocket.Message.Send(conn, fmt.Sprintf("%d", slide)); err != nil {
				registry.removeListener(slideId, listener)
				return
			}

			continue
		}

		if err := ping(conn); err != nil {
			registry.removeListener(slideId, listener)
			return
		}
	}
}

func ping(conn *websocket.Conn) error {
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	if err := websocket.Message.Send(conn, "ping"); err != nil {
		return err
	}

	var pong string
	if err := websocket.Message.Receive(conn, &pong); err != nil {
		return err
	}

	return nil
}

type slideListener struct {
	sync.Mutex
	slide int
	ch    chan int
}

func (l *slideListener) set(slide int) {
	select {
	case l.ch <- slide:
	default:
		l.Lock()
		defer l.Unlock()
		l.slide = slide
	}
}

func (l *slideListener) get(timeout time.Duration) int {
	l.Lock()
	curSlide := l.slide
	l.slide = 0
	l.Unlock()

	if curSlide != 0 {
		return curSlide
	}

	timer := time.NewTimer(timeout)
	select {
	case slide := <-l.ch:
		timer.Stop()
		return slide
	case <-timer.C:
		return 0
	}
}

type listenerRegistry struct {
	sync.Mutex
	listeners map[string][]*slideListener
}

func (r *listenerRegistry) addListener(slideId string, listener *slideListener) {
	r.Lock()
	defer r.Unlock()
	r.listeners[slideId] = append(r.listeners[slideId], listener)
}

func (r *listenerRegistry) removeListener(slideId string, listener *slideListener) {
	r.Lock()
	defer r.Unlock()

	var i int
	for i = range r.listeners[slideId] {
		if r.listeners[slideId][i] == listener {
			break
		}
	}

	listeners := r.listeners[slideId]
	listeners = append(listeners[:i], listeners[i+1:]...)

	r.listeners[slideId] = listeners
}

func (r *listenerRegistry) setSlide(slideId string, slide int) {
	r.Lock()
	defer r.Unlock()

	for _, listener := range r.listeners[slideId] {
		listener.set(slide)
	}
}
