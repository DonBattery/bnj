"use strict";

// WebSocketManager is the object responsible for WebSocket communication
class WebSocketManager {
  constructor(){
    this.ws = null;
    this.responseListeners = {};

    this.ready          = () => this.ws && this.ws.readyState == WebSocket.OPEN;
    this.initWs         = this.initWs.bind(this);
    this.setup          = this.setup.bind(this);
    this.notify         = this.notify.bind(this);
    this.request        = this.request.bind(this);
    this.msgHandler     = this.msgHandler.bind(this);
    this.handleChat     = this.handleChat.bind(this);
    this.handleUpdate   = this.handleUpdate.bind(this);
    this.handleResponse = this.handleResponse.bind(this);
  };

  initWs() {
    if (this.ws && this.ws.readyState == 1) {
      console.log("Websocket is already connected")
      return
    }
    let uri = getWSURL();
    console.log(`Initiating WebSocket with connection string: ${uri}`);

    try {
      this.ws = new WebSocket(uri);
    } catch (exception) {
      console.error("Failed to create WebSocket object", exception)
      return
    }

    Status.update("WS", "❎");

    this.setup();
  };

  setup() {
    if (!this.ws) {
      console.error("WebSocketManager.setup called, but ws is not connected");
      return
    };

    this.ws.onopen = () => {
      // TODO: Implement recconnect
      console.log("WebSocket Connected with ID", ClientId);
    };

    this.ws.onerror = event => {
      console.error("Connection error:", event);
    };

    this.ws.onmessage = event => {
      var parsed;
      try {
        parsed = JSON.parse(event.data);
      } catch (exception) {
        console.error("Failed to JSON parse incoming Server Message", exception);
        return
      };
      console.log(parsed);
      this.msgHandler(parsed);
    };
  };

  notify(msg) {
    if (!this.ready()) {
      console.error("WebSocketManager.notify called, but ws is not open");
      return
    };
    try {
      this.ws.send(JSON.stringify(msg));
    } catch (exception) {
      console.error("Failed to send WebSocket notification message", exception);
    };
  };

  request(requestType, requestBody, onResponse) {
    if (!this.ready()) {
      console.error("WebSocketManager.request called, but ws is not open");
      return
    };

    let requestId = randomId();

    this.responseListeners[requestId] = onResponse;

    try {
      this.ws.send(JSON.stringify(new RequestMessage(requestId, requestType, requestBody)));
    } catch (exception) {
      console.error("Failed to send WebSocket notification message", exception);
    };
  };

  msgHandler(msg) {
    if (!msg.hasOwnProperty("msg_type")) {
      console.error("Server Message has unknown msg_type")
      return
    };
    if (msg.msg_type == "chat") {
      this.handleChat(msg.chat);
    };
    if (msg.msg_type == "update") {
      this.handleUpdate(msg.objects);
    };
    if (msg.msg_type == "response") {
      this.handleResponse(msg.response);
    };
  };

  handleChat(chat) {
    console.log(`CHAT Channel: ${chat.channel} Message: ${chat.message}`);
  }

  handleUpdate(objects) {
    console.log("Update", objects);
  }

  // handleResponse calls the registered onResponse function of a previous request
  handleResponse(response) {
    let reqId = response.request_id;
    if (this.responseListeners.hasOwnProperty(reqId)) { // if there is a registered response listener
      this.responseListeners[reqId](response); // call the function
      delete this.responseListeners[reqId]; // then delete the listener
    };
  };

};
