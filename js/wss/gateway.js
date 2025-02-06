
class Gateway{
    /**
     *
     * @param {Client} client
     * @param {Security} security
     * @param {Auth} auth
     */
    constructor(client, security, auth){
        this.client = client;
        this.security = security;
        this.auth = auth;
        this.ready = ()=>{}
        this.paused = ()=>{}
    }

    start(){
        this.client.log("gateway start")
        this.ready()
    }
    stop(){
        this.client.log("gateway stop")
        this.paused()
    }

    Start(){
        this.client.Start()
    }
    Stop(){
        this.client.Stop()
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