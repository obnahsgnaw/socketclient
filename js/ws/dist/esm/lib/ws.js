export class WebsocketClient {
    constructor(options) {
        this.connecting = false;
        this.url = options.url;
        this.reconnectInterval = options.reconnectInterval || 0;
        this.debug = options.debug || false;
        this.heartbeatInterval = options.heartbeatInterval || 0;
        this.heartbeatData = options.heartbeatData || "";
        this.heartbeatFn = 0;
        this.retried = false;
        this.cusClose = false;
        this.readyHandler = () => {
        };
        this.connectingHandler = () => {
        };
        this.closeHandler = () => {
        };
        this.messageHandler = (message) => {
            return "";
        };
        this.socket = null;
    }
    log(message) {
        if (this.debug) {
            console.log(message);
        }
    }
    onConnecting(fn) {
        this.connectingHandler = fn;
    }
    onReady(fn) {
        this.readyHandler = fn;
    }
    onClose(fn) {
        this.closeHandler = fn;
    }
    onError(fn) {
        this.errorHandler = fn;
    }
    onMessage(fn) {
        this.messageHandler = fn;
    }
    send(message) {
        if (this.socket !== null && !this.connecting) {
            this.log('WebSocket sent message:' + message);
            this.socket.send(message);
        }
        else {
            this.log('WebSocket send failed, connect not ready');
        }
    }
    close() {
        if (this.socket !== null) {
            this.cusClose = true;
            this.socket.close();
        }
    }
    pauseHeartbeat() {
        this.log("WebSocket heartbeat paused");
        if (this.heartbeatFn !== 0) {
            clearInterval(this.heartbeatFn);
        }
    }
    continueHeartbeat() {
        if (!this.connecting && this.socket !== null) {
            this.log("WebSocket heartbeat continue");
            this.heartbeat();
        }
    }
    connect() {
        if (this.connecting)
            return;
        this.connecting = true;
        this.log("WebSocket connecting");
        this.connectingHandler();
        this.socket = new WebSocket(this.url);
        this.socket.onopen = () => {
            this.log('WebSocket connection opened');
            this.connecting = false;
            this.readyHandler(this);
            this.heartbeat();
        };
        this.socket.onmessage = (event) => {
            this.log('WebSocket received message:' + event.data);
            const response = this.messageHandler(event.data);
            if (response !== "") {
                this.send(response);
            }
        };
        this.socket.onerror = (error) => {
            this.log('WebSocket error:' + error);
            if (this.errorHandler)
                this.errorHandler(error);
            this.connecting = false;
        };
        this.socket.onclose = (event) => {
            this.log('WebSocket connection closed:' + event.reason);
            this.connecting = false;
            this.closeHandler();
            this.socket = null;
            if (this.reconnectInterval > 0 && !this.cusClose) {
                if (this.retried) {
                    this.reconnectInterval = this.reconnectInterval * 2;
                }
                setTimeout(() => () => {
                    this.connect();
                    this.retried = true;
                }, this.reconnectInterval);
            }
            else {
                this.log('WebSocket reconnect disabled');
            }
        };
    }
    heartbeat() {
        if (this.heartbeatInterval > 0 && this.heartbeatData !== "") {
            this.heartbeatFn = setInterval(() => {
                this.send(this.heartbeatData);
            }, this.heartbeatInterval);
        }
        else {
            this.log("Websocket heartbeat disabled");
        }
    }
}
