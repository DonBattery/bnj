"use strict";

// ChatNotify is sent to the server when the user creates a chat message
class ChatNotify {
  constructor(channel, message) {
    this.channel = channel;
    this.message = message;
  };
};

// ControlNotify is sent to the server when the user initiates a control command
class ControlNotify {
  constructor(controlType, controlKey) {
    this.control_type = controlType;
    this.contorl_key = controlKey;
  };
};

// ClientNotify is a message sent to the server without the need of a direct response
class ClientNotify {
  constructor(notifyType, opts){
    opts = opts || {};
    this.notify_type = notifyType;
    if (opts.hasOwnProperty("chat")) {
      this.chat = opts.chat;
    };
    if (opts.hasOwnProperty("control")) {
      this.control = opts.control;
    };
  };
};

// ClientRequest is the kind of ClientMsg that requires a direct response from the server
class ClientRequest {
  constructor(requestId, requestType, requestBody) {
    this.request_id = requestId;
    this.request_type = requestType;
    this.request_body = requestBody;
  };
};

// ClientMsg is the actual object that will be sent to the server on the WebSocket
class ClientMsg {
  constructor(msgType, opts) {
    opts = opts || {};
    this.msg_type = msgType;
    if (opts.hasOwnProperty("request")) {
      this.request = opts.request;
    };
    if (opts.hasOwnProperty("notify")) {
      this.notify = opts.notify;
    };
  };
};

// ChatMessage creates a chat type ClientMsg
function ChatMessage(channel, message) {
  return new ClientMsg("notify", {
    notify: new ClientNotify("chat", {
      chat: new ChatNotify(channel, message),
    }),
  });
};

// ControlMessage creates a control type ClientMsg
function ControlMessage(controlType, controlKey) {
  return new ClientMsg("notify", {
    notify: new ClientNotify("control", {
      contrl: new ControlNotify(controlType, controlKey),
    }),
  });
};

// RequestMessage creates a request type ClientMsg
function RequestMessage(requestId, requestType, requestBody) {
  return new ClientMsg("request", {
    request: new ClientRequest(requestId, requestType, requestBody),
  });
};
