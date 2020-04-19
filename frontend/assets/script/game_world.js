"use strict";

class GameWorld {
  constructor(rules, players, worldMap, gameObjects) {
    this.world_rules = rules;
    this.players = players;
    this.world_map = worldMap;
    this.game_objects = gameObjects;

    this.width    = () => this.world_map.width();
    this.height   = () => this.world_map.height();
    this.widthPx  = () => this.world_map.width() * this.world_rules.block_size;
    this.heightPx = () => this.world_map.height() * this.world_rules.block_size;
  };
};

class WorldRules {
  constructor(blockSize, maxPlayer, minPlayer, targetScore, waitTime) {
    this.block_size   = blockSize;
    this.max_player   = maxPlayer;
    this.min_player   = minPlayer;
    this.target_score = targetScore;
    this.wait_time    = waitTime;
  };
};

class WorldMap {
  constructor(background, rows) {
    this.background = background;
    this.rows       = rows;

    this.width  = () => this.rows[0].length;
    this.height = () => this.rows.length;
    this.get    = (x, y) => this.rows[y][x];
  };
};
