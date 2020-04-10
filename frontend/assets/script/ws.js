"use strict";

class WebSocketManager {
  constructor(){
    this.ws = undefined;

    this.connect = this.connect.bind(this);
    this.setup   = this.setup.bind(this);
    this.send    = this.send.bind(this);
  };

  connect() {
    this.ws = new WebSocket(getWSURL());
    this.setup();
  };

  setup() {
    this.ws.onopen = () => {
      console.log('Connected with ID', ClientID);
      this.sendMsg({
        type: "conn",
        payload: {
          event: "connected",
          client_id: ClientID,
        },
      });
    };

    this.ws.onerror = (event) => {
      console.log('Connection error:', event);
    };

    this.ws.onmessage = function(evt) {
      console.log("WS MSG", evt.data);
    };
  };

  send(msg) {
    this.ws.send(JSON.stringify(msg));
  };
};
