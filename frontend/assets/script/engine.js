"use strict";

class Engine {
  constructor(step, update, render) {
    this.step            = step,
    this.render          = () => {};
    this.timeSinceUpdate = 0;
    this.frame           = undefined,
    this.engineTime      = undefined,
    this.updated         = false;

    this.run        = this.run.bind(this);
    this.update     = this.update.bind(this);
    this.setRender  = this.setRender.bind(this);
    this.initEngine = this.initEngine.bind(this);
    this.stop       = this.stop.bind(this);
  };

  initEngine() {
    this.timeSinceUpdate = this.step;
    this.engineTime = window.performance.now();
    this.frame = window.requestAnimationFrame(this.run);
    return this.update;
  };

  run(rightNow) {
    this.frame = window.requestAnimationFrame(this.run);

    let deltaTime = rightNow - this.engineTime;
    this.timeSinceUpdate += deltaTime;
    this.engineTime = rightNow;

    if (this.updated && this.timeSinceUpdate >= this.step) {
      this.updated = false;
      this.render(1 / (this.timeSinceUpdate / 1000)); // pass the FPS to the render function
      this.timeSinceUpdate = 0;
    };
  };

  setRender(renderFn) {
    this.render = renderFn;
  };

  update() {
    this.updated = true;
  };


  stop() {
    window.cancelAnimationFrame(this.frame);
  };
};
