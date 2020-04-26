"use strict";

// LoginRequest is the body of a login request sent to the server
class LoginRequest {
  constructor(name, color) {
    this.name = name;
    this.color = color;
  };
};

// LoginManager is responsible to manage to login page and sending
// login requests to the server and trigger onSuccessFn with the world_data
class LoginManager {
  constructor() {
    this.authenticated = false;

    this.requestFn     = () => {};
    this.onSuccessFn   = () => {};

    this.loginPage   = document.getElementById("LoginPage");
    this.nameField   = document.getElementById("LoginName");
    this.errorField  = document.getElementById("InputError");
    this.colorField  = document.getElementById("LoginColor");
    this.loginButton = document.getElementById("LoginButton");

    this.initLogin  = this.initLogin.bind(this);
    this.auth       = this.auth.bind(this);
    this.validate   = this.validate.bind(this);
    this.onResponse = this.onResponse.bind(this);
    this.showError  = this.showError.bind(this);
  };

  initLogin() {
    this.loginPage.classList.add("activePage");
    this.loginButton.addEventListener("click", this.auth);
  };

  validate() {
    return this.nameField.value.length >= 3 && this.nameField.value.length <= 12;
  };

  auth() {
    if (!this.validate()) {
      this.showError("Invalid Name");
      return
    };
    this.requestFn("login", JSON.stringify(new LoginRequest(this.nameField.value, this.colorField.value)), this.onResponse);
  };

  onResponse(resp) {
    if (resp.status != 202) {
      // console.error("Login failed", resp);
      this.showError(resp.payload);
      return
    };
    resp.payload = JSON.parse(resp.payload);
    this.loginPage.classList.remove("activePage");
    this.authenticated = true;
    this.onSuccessFn(resp.payload);
  };

  showError(msg) {
    this.errorField.classList.remove("hidden");
    this.errorField.innerHTML = msg;
  };
};