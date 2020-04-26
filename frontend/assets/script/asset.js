"use strict";

class AssetManager {
  constructor() {
    this.assets = {};

    document.querySelectorAll(".imgAsset").forEach(asset => {
      this.assets[asset.getAttribute("data-name")] = asset;
    });

    this.getAll = () => this.assets;
  };
};
