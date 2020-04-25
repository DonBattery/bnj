"use strict";

class LoginRequest {
  constructor(name, color) {
    this.name = name;
    this.color = color;
  };
};

class LoginManager {
  constructor(requestFn, onSuccessFn) {
    this.authenticated = false;
    this.requestFn     = requestFn;
    this.onSuccessFn   = onSuccessFn;

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
    return this.nameField.value.length >= 3 && this.nameField.value.length <= 16;
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