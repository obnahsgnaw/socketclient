import {WebSocket} from "ws";

export class Client {
    /**
     *
     * @param {(String|URL)} address The URL to which to connect
     */
    constructor(address) {
        this.address = address;
        this.socket = null;
        this.connected = false
        this.retries = 0
        this.retryMax = 5
        this.cusClose = false
        this.debug = true
        this.sendInterceptors = []
        this.receiveInterceptors = []
        this._onConnected = function (){}
        this._onDisconnected = function (){}
        this._onErrored = function (e){}
        this._onMessaged = function (data){}
    }

    log(msg){
        const that = this;
        if (that.debug){
            console.log(that.address+":"+msg);
        }
    }

    DisableDebug(){
        this.debug=false;
    }

    /**
     *
     * @param {String} message
     * @returns {String}
     * @constructor
     */
    SendMessage(message){
        const that = this;
        if (that.socket != null && that.connected){
            try{
                const il = that.sendInterceptors.length
                if (il >0){
                    for (let i = il-1; i >0; i--){
                        message = that.sendInterceptors[i](message)
                        if (message === ""){
                            return ""
                        }
                    }
                }
                that.socket.send(message)
                this.log("sent message:"+message)
                return ""
            }catch(e){
                return e.message
            }
        }
    }

    /**
     *
     * @param {String} message
     * @returns {String}
     * @constructor
     */
    SendRawMessage(message){
        const that = this;
        if (that.socket != null && that.connected){
            try{
                that.socket.send(message)
                this.log("sent message:"+message)
                return ""
            }catch(e){
                return e.message
            }
        }
    }

    WhenReady(fn){
        const that = this;
        if (typeof fn === "function"){
            that._onConnected = fn
        }
    }

    WhenPaused(fn){
        const that = this;
        if (typeof fn === "function"){
            that._onDisconnected = fn
        }
    }

    OnError(fn){
        const that = this;
        if (typeof fn === "function"){
            that._onErrored = fn
        }
    }

    /**
     *
     * @param {function(String)} fn
     * @constructor
     */
    OnMessage(fn){
        const that = this;
        if (typeof fn === "function"){
            that._onMessaged = function (data){
                fn(data)
            }
        }
    }

    Stop(){
        const that = this;
        that.cusClose = true;
        if (that.socket != null){
            that.socket.close()
        }
    }

    /**
     *
     * @param {function(String)} fn
     * @constructor
     */
    AddSendInterceptor(fn){
        const that = this;
        if (typeof fn === "function"){
            that.sendInterceptors[that.sendInterceptors.length] = fn
        }
    }

    /**
     *
     * @param {function(String)} fn
     * @constructor
     */
    AddReceiveInterceptor(fn){
        const that = this;
        if (typeof fn === "function"){
            that.receiveInterceptors[that.receiveInterceptors.length] = fn
        }
    }

    Start(){
        const that = this;
        if (that.socket == null){
            if (that.retries > 0){
                that.log("reconnecting...")
            }else{
                that.log("connecting...")
            }
            that.socket = new WebSocket(that.address);
            that.socket.onopen = function (e){
                that.connected = true;
                that.log("connected")
                that._onConnected(e);
            }
            that.socket.onclose = function (e){
                that.connected = false;
                that.log("disconnected")
                that._onDisconnected(e);
                if (!that.cusClose && that.retries < that.retryMax){
                    that.Open()
                    that.retries++
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
                that.log("errored,err="+e.message)
                that._onErrored(e);
            }
        }
    }
}