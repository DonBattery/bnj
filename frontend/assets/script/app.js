"use strict";

class App {
  constructor() {
    this.game  = new GameManager();
    this.login = new LoginManager(this.game.initGame);

    this.run = () => { this.login.initLogin(); };
  };
};

window.addEventListener("load", () => {
  let app = new App();
  app.run();
});