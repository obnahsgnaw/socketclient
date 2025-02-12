var auth = import('proto/gateway/v1/auth_pb.js')
class Auth{
    /**
     *
     * @param {Client} client
     * @param {Security} security
     * @param {Token} token
     */
    constructor(client, security,token){
        this.client = client;
        this.security = security;
        this.token = token;
        this.ready = ()=>{}
        this.paused = ()=>{}
        this.security.WhenReady(this.start)
        this.security.WhenPaused(this.stop)
    }
    start(){
        this.client.log("auth start")
        // this.client.SendMessage()
        this.ready()
        var pkg = new auth.AuthRequest()
    }
    stop(){
        this.client.log("auth stop")
        this.paused()
    }

    WhenReady(fn){
        if (typeof fn === "function"){
            this.ready = fn
        }
    }

    WhenPaused(fn){
        if (typeof fn === "function"){
            this.paused = fn
        }
    }
}

let Token = {
    AppId:"",
    Token:"",
}