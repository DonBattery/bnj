"use strict";

class GameManager {
  constructor() {
    this.display = null;
    this.world = null;
    this.engine  = null;
    this.input   = null;

    this.initGame = this.initGame.bind(this);
  };

  initGame(configs) {
    console.log("initGame called with Configs:");
    console.log(configs);
  };
};