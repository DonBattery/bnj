"use strict";

function randRange(min, max) {
  return Math.floor(Math.random() * (max - min + 1) ) + min;
};

function randInt(max) {
  return randRange(0, max - 1);
};

function randomId(length) {
  return Math.random().toString(36).substr(2, length);
};

function getWSURL() {
  let loc = window.location;
  let uri = 'ws:';
  if (loc.protocol === 'https:') {
    uri = 'wss:';
  }
  uri += '//' + loc.host;
  uri += loc.pathname + 'hub';

  return uri;
};

function fixWindow() {
  window.requestAnimFrame = (function(){
    return  window.requestAnimationFrame       ||
            window.webkitRequestAnimationFrame ||
            window.mozRequestAnimationFrame    ||
            window.ieRequestAnimationFrame     ||
            function( callback ){
              window.setTimeout(callback, 1000 / 60);
            };
          })();
};