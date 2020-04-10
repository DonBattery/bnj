"use strict";

class LoginManager {
  constructor(onAuth) {
    this.onAuth       = onAuth;
    this.autenticated = false;

    this.loginPage   = document.getElementById("LoginPage");
    this.nameField   = document.getElementById("LoginName");
    this.errorField  = document.getElementById("InputError");
    this.colorField  = document.getElementById("LoginColor");
    this.loginButton = document.getElementById("LoginButton");

    this.validate  = this.validate.bind(this);
    this.initLogin = this.initLogin.bind(this);
    this.auth      = this.auth.bind(this);
    this.showError = this.showError.bind(this);
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
    this.errorField.classList.add("hidden");
    console.log(this.nameField.value);
    console.log(this.colorField.value);
  };

  showError(msg) {
    this.errorField.classList.remove("hidden");
    this.errorField.innerHTML = msg;
  };
};