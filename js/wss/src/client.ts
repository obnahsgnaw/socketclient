interface Options {
    url: string;
    protocols?: string | string[];
    debug: boolean;
    reconnectTimeoutMs: number
    retryMax: number ;
}

export class Client {
    private option:Options
    private socket: WebSocket | null = null;
    private connected: boolean = false;
    private retries: number = 0;
    private cusClose: boolean = false;
    private sendInterceptors: any[] = [];
    private receiveInterceptors: any[] = [];
    private _onConnected: Function[] = [];
    private _onDisconnected: Function[] = [];
    private _onErrored: Function = function (){};
    private _onMessaged: (message: string) => void = function (message:string ){};
    constructor(options: Options ) {
        this.option = options;
        if (this.option.reconnectTimeoutMs === 0) {
            this.option.reconnectTimeoutMs = 5000
        }
        if (this.option.retryMax === 0) {
            this.option.retryMax = 5
        }
    }

    public log(msg:string){
        const that = this;
        if (that.option.debug){
            console.log(that.option.url+":"+msg);
        }
    }

    public open(){
        const that = this;
        if (that.socket == null){
            if (that.retries > 0){
                that.log("reconnecting...")
            }else{
                that.log("connecting...")
            }
            that.socket = new WebSocket(that.option.url,that.option.protocols);
            that.socket.onopen = function (e){
                that.connected = true;
                that.log("connected")
                that._onConnected.forEach(fn => fn(e));
            }
            that.socket.onclose = function (e){
                that.connected = false;
                that.log("disconnected")
                that._onDisconnected.forEach(fn => fn(e));
                if (!that.cusClose && that.retries < that.option.retryMax){
                    setTimeout(()=>{
                        that.open()
                        that.retries++
                    },that.option.reconnectTimeoutMs)
                }
            }
            that.socket.onmessage = function (e){
                that.log("received message:"+e.data)
                let data = e.data
                const il = that.receiveInterceptors.length
                if (il >0){
                    for (let i = 0; i < il; i++){
                        data = that.receiveInterceptors[i](data)
                        if (data === ""){
                            return
                        }
                    }
                }
                that._onMessaged(data);
            }
            that.socket.onerror = function (e){
                that.log("errored")
                that._onErrored(e);
                if (!that.cusClose && that.retries < that.option.retryMax){
                    setTimeout(()=>{
                        that.open()
                        that.retries++
                    },that.option.reconnectTimeoutMs)
                }
            }
        }
    }

    public close(){
        const that = this;
        that.cusClose = true;
        if (that.socket != null){
            that.socket.close()
        }
        that.log("closed")
    }

    public send(message:string, raw:boolean = false):string{
        const that = this;
        if (that.socket != null && that.connected){
            try{
                if (!raw){
                    const il = that.sendInterceptors.length
                    if (il >0){
                        for (let i = il-1; i >0; i--){
                            message = that.sendInterceptors[i](message)
                            if (message === ""){
                                return ""
                            }
                        }
                    }
                }
                that.socket.send(message)
                this.log("sent message:"+message)
                return ""
            }catch(e:any){
                return e.message
            }
        }
        return "not connected"
    }

    public ready(fn:Function){
        const that = this;
        that._onConnected.push(fn)
    }

    public paused(fn:Function){
        const that = this;
        that._onDisconnected.push(fn)
    }

    public errored(fn:Function){
        const that = this;
        that._onErrored = fn
    }

    public messaged(fn:(message: string) => void){
        const that = this;
        that._onMessaged = fn
    }

    public addSendInterceptor(fn:(message: string) => string){
        const that = this;
        that.sendInterceptors[that.sendInterceptors.length] = fn
    }

    public addReceiveInterceptor(fn:(message: string) => string){
        const that = this;
        that.receiveInterceptors[that.receiveInterceptors.length] = fn
    }
}