import JSEncrypt from 'jsencrypt'

class Security {
    /**
     *
     * @param {Client} client
     */
    constructor(client) {
        this.client = client;
        this.id = ""
        this.dataType = "json"
        this.esType = "aes256"
        this.rsa = null
        this.ades = null
        this.key = ""
        this.initialized = false
        this.encryptEnable = false
        this.err = (msg)=>{
            this.client.log(msg)
        }
        this.ready = ()=>{}
        this.paused = ()=>{}
        this.client.AddReceiveInterceptor( (message)=>{
            if (this.initialized){
                return message
            }
            if (message === "000" || message === "111"){
                this.initialized = true
                this.encryptEnable = message === "111"
                this.client.log("authenticate initialized success")
                this.ready()
            }
            this.client.log("authenticate initialized failed with:"+message)
            this.err("authenticated init failed with:"+message)
        })
        this.client.AddReceiveInterceptor( (message)=>{
            if (this.encryptEnable){
                // 解密 TODO
            }
            return message
        })
        this.client.AddSendInterceptor((message)=>{
            if (this.encryptEnable){
                // 加密 TODO
            }
            return message
        })
        this.client.WhenPaused(this.stop)
        this.client.WhenReady( ()=>{
               this.start()
        })
    }

    SetPublicKey(publicKey){
        this.rsa = new JSEncrypt()
        this.rsa.SetPublicKey(publicKey)
    }

    start(){
        this.client.log("authenticate start")
        this.client.SendRawMessage(this.authenticateMessage())
    }

    stop(){
        this.paused()
        this.initialized = false
        this.encryptEnable = false
        this.client.log("authenticate stop")
    }

    Error(fn){
        if (typeof fn === "function"){
            this.err = fn
        }
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

    authenticateMessage(){
        return `user@${this.id}@${this.dataType}::${this.encryptedKey()}`
    }
    encryptedKey(){
        if (this.rsa != null){
            this.randEsKey()
            return this.rsa.encrypt(this.key+Math.floor(Date.now() / 1000))
        }
        return ""
    }
    randEsKey(){
        let len = 8
        if (this.esType === "des"){
            len = 8
        }else{
            len = parseInt(this.esType)/8
        }
        this.key = this.rand(len)
    }
    rand(num){
        let chars = ["0","1","2","3","4","5","6","7","8","9","A","B","C","D","E","F","G","H","I","J","K","L","M","N","O","P","Q","R","S","T","U","V","W","X","Y","Z","a","b","c","d","e","f","g","h","i","j","k","l","m","n","o","p","q","r","s","t","u","v","w","x","y","z"];
        let res = ""
        for (let i=0;i<num;i++){
            let id = Math.floor(Math.random()*(chars.length))
            res += chars[id]
        }
        return res
    }
}