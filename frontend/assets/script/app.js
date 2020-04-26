"use strict";

// App is the top level application it is invoked when the page is fully loaded
class App {
  constructor() {

    // Create the time Engine with 33 millisecond step (yielding ~30 FPS)
    this.engine = new Engine(33);

    // Create the AssetManager
    this.assets = new AssetManager();

    // Create the WebSocketManager
    this.ws = new WebSocketManager();

    // Create the LoginManager
    this.login = new LoginManager();

    // Create the InputManager
    this.input = new InputManager();

    // The Display will be created upon suddesful login
    this.display = null;

    // The GameWorld will be created upon suddesful login
    this.world = null;

    // The LoginManager can use the WebSocketManager's request method
    // to request login access to the game
    this.login.requestFn = this.ws.request;

    // The LoginManager will create the GameWorld and the Display upon successful login
    // using the World Data retrieved from the server
    this.login.onSuccessFn = world_data => {
      this.world = new GameWorld(world_data);
      this.display = new Display(this.world, this.assets.getAll());
      this.engine.initEngine(this.display.render);
    };

    // On incoming server updates the GameWorld will be updated accordingly
    // and and drawn to the Display
    this.ws.onUpdateFn = update => {
      this.world.updateWorld(update);
      this.display.drawWorld(this.world);
    };

    this.run = () => {
      this.ws.initWs();
      this.login.initLogin();
    };

    this.onLogin = (world_data) => {
      this.wolrd = new GameWorld(world_data);
    };
  };
};

function setupApp() {
  // Create the StatusManager with the StatusBar
  Status = new StatusManager(document.getElementById("StatusBar"));
  Status.update("WS", "❌");
  Status.update("FPS", 0);
  Status.show();
  // Create the App and run it
  let app = new App();
  app.run();
};

window.addEventListener("load", setupApp);
