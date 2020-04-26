"use strict";

// GameWorld representes the static game world (platforms)
// and the dynamic elements in it, like the players and game objects
class GameWorld {
  constructor(opts) {
    this.world_rules   = new WorldRules(opts.world_rules || {});
    this.world_map     = new WorldMap(opts.world_map || {});
    this.players       = opts.players       || [];
    this.world_objects = opts.world_objects || [];

    this.updateWorld = opts => {
      this.players       = opts.players       || this.players;
      this.world_objects = opts.world_objects || this.world_objects;
    };

    this.width    = () => this.world_map.width();
    this.height   = () => this.world_map.height();
    this.widthPx  = () => this.world_map.width() * this.world_rules.block_size;
    this.heightPx = () => this.world_map.height() * this.world_rules.block_size;
  };
};

// WorldRules are the rules of the GameWorld
class WorldRules {
  constructor(opts) {
    this.block_size   = opts.block_size || 16;
    this.max_player   = opts.max_player;
    this.min_player   = opts.min_player;
    this.target_score = opts.target_score;
    this.wait_time    = opts.wait_time;
    this.gravity      = opts.gravity;
    this.friction     = opts.friction;
  };
};

// WorldMap is the map of the GameWorld
class WorldMap {
  constructor(opts) {
    this.background = opts.background || "black";
    this.rows       = opts.rows || [""];

    this.width  = () => this.rows[0].length;
    this.height = () => this.rows.length;
    this.get    = (x, y) => this.rows[y][x];
  };
};
