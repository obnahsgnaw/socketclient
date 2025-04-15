interface WebSocketClientOptions {
    url: string;
    reconnectInterval?: number;  // 重连间隔（毫秒）
    debug?: boolean
    heartbeatInterval?: number
    heartbeatData?: string
}

export class WebsocketClient {
    private readonly url: string;
    private reconnectInterval: number;
    private retried: boolean
    private readonly debug: boolean;
    private heartbeatInterval: number
    private heartbeatData: string
    private heartbeatFn: any
    private socket: WebSocket | null
    private connecting: boolean = false;
    private readyHandler: (ws: WebsocketClient) => void
    private closeHandler: () => void
    private connectingHandler: () => void
    private messageHandler: (message: string) => string;
    private errorHandler?: (error: Event) => void;
    private cusClose: boolean

    constructor(options: WebSocketClientOptions) {
        this.url = options.url;
        this.reconnectInterval = options.reconnectInterval || 0;
        this.debug = options.debug || false;
        this.heartbeatInterval = options.heartbeatInterval || 0;
        this.heartbeatData = options.heartbeatData || "";
        this.heartbeatFn = 0
        this.retried = false
        this.cusClose = false;
        this.readyHandler = () => {
        }
        this.connectingHandler = () => {
        }
        this.closeHandler = () => {
        }
        this.messageHandler = (message: string) => {
            return ""
        }
        this.socket = null
    }

    private log(message: string): void {
        if (this.debug) {
            console.log(message);
        }
    }

    public onConnecting(fn: () => void) {
        this.connectingHandler = fn;
    }

    public onReady(fn: (ws: WebsocketClient) => void) {
        this.readyHandler = fn;
    }

    public onClose(fn: () => void) {
        this.closeHandler = fn;
    }

    public onError(fn: (error: Event) => void) {
        this.errorHandler = fn;
    }

    public onMessage(fn: (message: string) => string) {
        this.messageHandler = fn;
    }

    public send(message: string) {
        if (this.socket !== null && !this.connecting) {
            this.log('WebSocket sent message:' + message);
            this.socket.send(message);
        } else {
            this.log('WebSocket send failed, connect not ready');
        }
    }

    public close() {
        if (this.socket !== null) {
            this.cusClose = true
            this.socket.close();
        }
    }

    public pauseHeartbeat() {
        this.log("WebSocket heartbeat paused");
        if (this.heartbeatFn !== 0) {
            clearInterval(this.heartbeatFn);
        }
    }

    public continueHeartbeat() {
        if (!this.connecting && this.socket !== null) {
            this.log("WebSocket heartbeat continue");
            this.heartbeat()
        }
    }

    public connect() {
        if (this.connecting) return;
        this.connecting = true;
        this.log("WebSocket connecting");
        this.connectingHandler()

        this.socket = new WebSocket(this.url);

        this.socket.onopen = () => {
            this.log('WebSocket connection opened');
            this.connecting = false;
            this.readyHandler(this)
            this.heartbeat()
        };

        this.socket.onmessage = (event) => {
            this.log('WebSocket received message:' + event.data);
            const response = this.messageHandler(event.data);
            if (response !== "") {
                this.send(response)
            }
        };

        this.socket.onerror = (error) => {
            this.log('WebSocket error:' + error);
            if (this.errorHandler) this.errorHandler(error);
            this.connecting = false;
        };

        this.socket.onclose = (event) => {
            this.log('WebSocket connection closed:' + event.reason);
            this.connecting = false;
            this.closeHandler();
            this.socket = null;
            if (this.reconnectInterval > 0 && !this.cusClose) {
                if (this.retried) {
                    this.reconnectInterval = this.reconnectInterval * 2
                }
                setTimeout(() => () => {
                    this.connect()
                    this.retried = true
                }, this.reconnectInterval);
            } else {
                this.log('WebSocket reconnect disabled');
            }
        };
    }

    private heartbeat(): void {
        if (this.heartbeatInterval > 0 && this.heartbeatData !== "") {
            this.heartbeatFn = setInterval(() => {
                this.send(this.heartbeatData);
            }, this.heartbeatInterval);
        } else {
            this.log("Websocket heartbeat disabled");
        }
    }
}

export type { WebSocketClientOptions }