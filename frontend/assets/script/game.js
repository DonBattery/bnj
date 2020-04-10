"use strict";

class GameManager {
  constructor() {
    this.display = undefined;
    this.engine  = new Engine();
    // this.assets  = new AssetManager();
    this.input   = new InputManager();
    this.ws      = new WebSocketManager();

    this.initGame = this.initGame.bind(this);
  };

  initGame(configs) {
    // console.log(configs);

    // Load the assets
    // await this.assets.load();

    this.display = new Display(configs, this.assets);

    // Init the display with the configs and assets and get the reDraw function
    let reDraw = this.display.initDisplay(configs, this.assets);

    // Set the engine's render function to the display's render funtion
    this.engine.setRender(this.display.render);

    // Start the engine and get the onUpdate function
    let onUpdate = this.engine.initEngine();

    // Set the display finished with drawing it should set the updated value in the engine to ture.
    this.display.setUpdater(onUpdate);

    // Init the WebSocket manager with the onUpdate functions and get the onInput function
    let onInput = this.ws.initWS(onUpdate);

    // Set up the controls and pass the onInput callback function
    this.input.listen(onInput);
  };
};