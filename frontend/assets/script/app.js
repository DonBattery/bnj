"use strict";

class App {
  constructor() {
    this.ws    = new WebSocketManager();
    this.game  = new GameManager();
    this.login = new LoginManager(this.ws.request, this.game.initGame);

    this.run = () => {
      this.ws.initWs();
      this.login.initLogin();
    };
  };
};

window.addEventListener("load", () => {
  // Create the StatusManager with the StatusBar
  Status = new StatusManager(document.getElementById("StatusBar"));
  Status.update("WS", "❌");
  Status.update("FPS", 0);
  Status.show();

  let app = new App();
  app.run();
});