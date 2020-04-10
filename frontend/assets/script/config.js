"use strict";

class WorldMap {
  constructor(worldData) {
    this.data = worldData || getDefaultLevel();
    this.squareSize = 16;

    this.getWidth  = this.getWidth.bind(this);
    this.getHeight = this.getHeight.bind(this);
  };

  getWidth() {
    return 22;
  };

  getHeight() {
    return 16;
  };
};

class ConfigManager {
  constructor(worldMap) {
    this.worldMap = worldMap;

    this.getPxWidth  = this.getPxWidth.bind(this);
    this.getPxHeight = this.getPxHeight.bind(this);
  };

  getPxWidth() {
    return this.worldMap.getWidth() * this.worldMap.squareSize;
  };

  getPxHeight() {
    return this.worldMap.getHeight() * this.worldMap.squareSize;
  };
};