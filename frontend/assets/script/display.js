"use strict";

class Display {
  constructor(configs, assets) {
    this.updater       = () => {};
    this.configs       = configs;
    this.assets        = assets;
    this.width         = configs.getPxWidth();
    this.height        = configs.getPxHeight();
    this.widthToHeight = this.width / this.height;
    this.gamePage      = document.getElementById("GamePage");
    this.statusBar     = document.getElementById("StatusBar");
    this.FPSMeter      = document.getElementById("FPSMeter");
    this.gameCanvas    = document.getElementById("GameCanvas");
    this.ctx           = this.gameCanvas.getContext("2d");
    this.buffer        = document.createElement("canvas").getContext("2d");


    this.buffer.canvas.height = this.height;
    this.buffer.canvas.width  = this.width;

    this.initDisplay = this.initDisplay.bind(this);
    this.setUpdater  = this.setUpdater.bind(this);
    this.resize      = this.resize.bind(this);
    this.render      = this.render.bind(this);
    this.drawObject  = this.drawObject.bind(this);
    this.drawBunny   = this.drawBunny.bind(this);
  };

  setUpdater(updaterFn) {
    this.updater = updaterFn;
  };

  resize() {
    let newWidth = window.innerWidth - 4;
    let newHeight = window.innerHeight - 4;
    let newWidthToHeight = newWidth / newHeight;

    if (newWidthToHeight > this.widthToHeight) {
        newWidth = newHeight * this.widthToHeight;
    } else {
        newHeight = newWidth / this.widthToHeight;
    }

    this.gamePage.style.width = newWidth + "px";
    this.gamePage.style.height = newHeight + "px";

    this.gamePage.style.marginTop = (-newHeight / 2) + "px";
    this.gamePage.style.marginLeft = (-newWidth / 2) + "px";

    this.gameCanvas.width = newWidth;
    this.gameCanvas.height = newHeight;

    this.gamePage.style.fontSize = (newWidth / this.width / 2) + 'em';
  };

  initDisplay(configs, assets) {
    this.configs = configs;
    this.assets = assets;
    return this.draw;
  };

  render(timeStamp, fps) {
    this.FPSMeter.innerHTML = "FPS: " + fps.toFixed(2);
    this.ctx.drawImage(this.buffer.canvas,
      0, 0, this.buffer.canvas.width, this.buffer.canvas.height,
      0, 0, this.ctx.canvas.width, this.ctx.canvas.height);
  };

  // draw everything according to the configs, assets and state
  draw(state) {
    // trigger the updater function (this should flip the engine's updated switch)
    this.updater();
  };

  drawObject(image, source_x, source_y, destination_x, destination_y, width, height) {
    this.buffer.drawImage(image,
      source_x, source_y, width, height,
      Math.round(destination_x), Math.round(destination_y), width, height);
  };

  drawBunny(image, source_x, source_y, destination_x, destination_y, width, height) {
    this.buffer.drawImage(image,
      source_x, source_y, width, height,
      Math.round(destination_x), Math.round(destination_y), width, height);
  };
};
