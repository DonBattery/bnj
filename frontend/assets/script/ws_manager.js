"use strict";

// WebSocketManager is the object responsible for WebSocket communication
class WebSocketManager {
  constructor(){
    this.ws = null;
    this.responseListeners = {};
    this.serverUpdates = 0 // number of received server updates

    // onUpdateFn needs to be overriden with the GameWorld's onUpdate method
    this.onUpdateFn = update => { console.log(update); };

    // ready returns true if the WebSocket is ready for read and write
    this.ready          = () => this.ws && this.ws.readyState == WebSocket.OPEN;

    this.initWs         = this.initWs.bind(this);
    this.notify         = this.notify.bind(this);
    this.request        = this.request.bind(this);
    this.msgHandler     = this.msgHandler.bind(this);
    this.handleChat     = this.handleChat.bind(this);
    this.handleResponse = this.handleResponse.bind(this);
  };

  initWs() {
    // Get the WebSocket URL
    let uri = getWSURL();
    console.log(`Initiating WebSocket connection with URL: ${uri}`);
    // Return if already connected
    if (this.ws && this.ws.readyState == 1) {
      console.log("Websocket is already connected")
      return
    }
    // Try to connect
    try {
      this.ws = new WebSocket(uri);
    } catch (exception) {
      console.error("Failed to create WebSocket object", exception)
      return
    }

     /////////////////////////////////////
    // Set up the WebSocket Connection //
   /////////////////////////////////////

    this.ws.onopen = () => {
      // TODO: Implement recconnect
      console.log("WebSocket Connected with ID", ClientID);
    };

    this.ws.onerror = event => {
      // TODO: Implement proper error handling
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
    // Update the status
    Status.update("WS", "❎");
  };

  msgHandler(msg) {
    if (!msg.hasOwnProperty("msg_type")) {
      console.error("Server Message has no type!")
      return
    };
    if (msg.msg_type == "chat") {
      this.handleChat(msg.chat);
      return
    };
    if (msg.msg_type == "update") {
      Status.update("WS Updates", this.serverUpdates++);
      this.onUpdateFn(msg.world_update);
      return
    };
    if (msg.msg_type == "response") {
      this.handleResponse(msg.response);
      return
    };
    console.error("Server Message has unknown type", msg.msg_type)
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

  handleChat(chat) {
    console.log(`CHAT Channel: ${chat.channel} Message: ${chat.message}`);
  };

  // handleResponse calls the registered onResponse function of a previous request
  handleResponse(response) {
    let reqId = response.request_id;
    if (this.responseListeners.hasOwnProperty(reqId)) { // if there is a registered response listener
      this.responseListeners[reqId](response); // call the function
      delete this.responseListeners[reqId]; // then delete the listener
    };
  };

};
