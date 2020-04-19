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
    this.world = new GameWorld(
      configs.world_rules,
      configs.players,
      new WorldMap(
        configs.world_map.background,
        configs.world_map.rows),
      configs.world_objects)

    this.display = new Display(this.world);

    // this.engine = new Engine(33, this.world.update, this.display.render);
    this.engine = new Engine(33, this.display.drawAll, this.display.render);
    this.engine.initEngine();
  };
};