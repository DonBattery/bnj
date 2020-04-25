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
    // Create the WorldMap from the configs
    let world_map = new WorldMap(configs.world_map.background, configs.world_map.rows);

    // Create the GameWorld from the configs
    this.world = new GameWorld(configs.world_rules, configs.players, world_map, configs.world_objects);

    // Create the Display from the world
    this.display = new Display(this.world);

    // Create the Engine with 33 millisecond step yielding ~ 30 FPS
    // and pass the display's update and render method in it
    // this.engine = new Engine(33, this.world.update, this.display.render);
    this.engine = new Engine(33, this.display.drawAll, this.display.render);

    // Initialize the Engine
    this.engine.initEngine();
  };
};