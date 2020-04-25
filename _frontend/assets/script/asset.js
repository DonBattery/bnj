"use strict";

class ImageAsset {
  constructor(img, name) {
    this.img = img;
    this.name = name;
  };
};

class GameObject {
  constructor(x, y) {

  };
};

class AssetManager {
  constructor() {

  };

  async load() {

  };

  preloadImages(images, callback) {
    let newImages = [];
    let loadedImages = 0;

    function postLoad() {
      loadedImages++
      if (loadedImages == images.length) {
        callback(newImages);
      };
    };

    for (let i = 0; i < images.length; i++) {
        newImages[i] = new Image();
        newImages[i].src = images[i];
        newImages[i].onload = () => { postLoad() };
        newImages[i].onerror = () => { postLoad() };
    };
  };

};
