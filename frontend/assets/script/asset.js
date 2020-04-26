"use strict";

class WorldObject {
  constructor(x, y, obj_type, anim, flip_x, flip_y) {
    this.x = x;
    this.y = y;
    this.obj_type = obj_type;
    this.anim = anim;
    this.flip_x = flip_x;
    this.flip_y = flip_y;
  };
};

class AssetManager {
  constructor() {
    this.assets = {};

    document.querySelectorAll(".imgAsset").forEach(asset => {
      this.assets[asset.getAttribute("data-name")] = asset;
    });

    this.getAll = () => this.assets;
  };
};
