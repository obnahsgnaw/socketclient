interface WebSocketClientOptions {
    url: string;
    reconnectInterval?: number;
    debug?: boolean;
    heartbeatInterval?: number;
    heartbeatData?: string;
}
export declare class WebsocketClient {
    private readonly url;
    private reconnectInterval;
    private retried;
    private readonly debug;
    private heartbeatInterval;
    private heartbeatData;
    private heartbeatFn;
    private socket;
    private connecting;
    private readyHandler;
    private closeHandler;
    private connectingHandler;
    private messageHandler;
    private errorHandler?;
    private cusClose;
    constructor(options: WebSocketClientOptions);
    private log;
    onConnecting(fn: () => void): void;
    onReady(fn: (ws: WebsocketClient) => void): void;
    onClose(fn: () => void): void;
    onError(fn: (error: Event) => void): void;
    onMessage(fn: (message: string) => string): void;
    send(message: string): void;
    close(): void;
    pauseHeartbeat(): void;
    continueHeartbeat(): void;
    connect(): void;
    private heartbeat;
}
export type { WebSocketClientOptions };
