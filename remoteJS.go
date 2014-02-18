// Copyright 2014, Chandra Sekar S.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the README.md file.

package main

const remoteJS = `
var remotePaused = false;
var path = window.location.pathname;
var wsURL = "ws://" + window.location.host +  path.substring(0, path.lastIndexOf('/')) + "/" + userRole;
var ws = new WebSocket(wsURL);

document.addEventListener("DOMContentLoaded", function() {
  if(userRole != "v") {
    return;
  }

  document.addEventListener("keypress", function(event) {
    if(event.charCode == 80 || event.charCode == 112) {
      if(!remotePaused) {
        remotePaused = true;
        document.title += " [PAUSED]";
      } else {
        remotePaused = false;
        document.title = document.title.substring(0, document.title.length - " [PAUSED]".length);
      }
    }
  });
});

ws.onmessage = function(event) {
  if(event.data == "ping") {
    ws.send("pong");
    return;
  }

  if(userRole == "v" && !remotePaused) {
    curSlide = Number(event.data) - 1;
    updateSlides();
  }
} 

ws.onopen = function(event) {
  ws.send(rSlideId);
  if(userRole == "p") {
	  ws.send(curSlide + 1);
  }
}

function sendRemote(curSlide) {
  if(userRole == "p") {
    ws.send(curSlide+1 + "");
  }
}`
