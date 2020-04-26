"use strict";

// StatusManager manages the status bar
class StatusManager {
  constructor(statusBar) {
    this.statusBar = statusBar;
    this.statuses = [];

    this.hide = () => {
      this.statusBar.classList.add("hidden");
    };

    this.show = () => {
      this.statusBar.classList.remove("hidden");
    };

    this.update = (name, value) => {
      let status = this.statuses.find(status => status.name == name);
      if (status) {
        status.value = value;
      } else {
        let newItem = document.createElement("div");
        newItem.classList.add("basicBorder", "statusItem");
        newItem.id = `${name}StatusItem`;
        this.statusBar.appendChild(newItem);
        this.statuses.push({
          name: name,
          value: value,
        });
      };
      this.render();
    };

    this.render = () => {
      this.statuses.forEach(status => {
        document.getElementById(`${status.name}StatusItem`).innerHTML = `${status.name}: ${status.value}`;
      });
    };
  };
};
