"use strict";

// Engine is animation manager that calls the Display's render then
// the GameManager's update method on every step in time
class Engine {
  constructor(step, update, render) {
    this.step            = step,
    this.update          = (update) ? update : () => {};
    this.render          = (render) ? render : () => {};
    this.timeSinceUpdate = 0;
    this.frame           = null,
    this.engineTime      = null,
    this.updated         = false;

    this.initEngine  = this.initEngine.bind(this);
    this.haltEngine  = this.haltEngine.bind(this);
    this.run = this.run.bind(this);
  };

  initEngine() {
    fixWindow();
    this.timeSinceUpdate = this.step;
    this.engineTime = window.performance.now();
    this.frame = window.requestAnimFrame(this.run);
  };

  haltEngine() {
    window.cancelAnimationFrame(this.frame);
  };

  run(rightNow) {
    this.frame = window.requestAnimFrame(this.run);

    let deltaTime = rightNow - this.engineTime;
    this.timeSinceUpdate += deltaTime;
    this.engineTime = rightNow;

    if (this.timeSinceUpdate >= this.step) {
      Status.update("FPS", 1 / (this.timeSinceUpdate / 1000)); // Calculate FPS
      this.render();
      this.timeSinceUpdate = 0;
    };

    this.update();
  };
};
