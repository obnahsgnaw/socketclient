import JSEncrypt from 'jsencrypt'
import {Client} from "./client";

type DataType = 'json' | 'proto'
type EsType = 'aes256' | 'aes128' | 'des';

class Security {
    private id: string = ""
    private type: string = "user"
    private client: Client
    private dataType: DataType = 'json'
    private esType: EsType = 'aes256'
    private rsa: JSEncrypt | null = null
    private aes: JSEncrypt | null = null
    private key: string = ""
    private initialized: boolean = false
    private encryptEnable: boolean = false
    private err: (msg: string) => void = msg => {
        this.client.log(msg)
    }
    private _ready: Function[] = [];
    private _paused: Function[] = [];

    constructor(client: Client, publicKey: string, id: string, type: string = "user") {
        this.client = client;
        this.id = id;
        this.type = type;
        this.client = client;
        this.setPublicKey(publicKey);
        this.client.addReceiveInterceptor((message: string) => {
            if (this.initialized) {
                return message
            }
            if (message === "000" || message === "111") {
                this.initialized = true
                this.encryptEnable = message === "111"
                this.client.log("authenticate initialized success")
                this._ready.forEach(fn => fn());
            } else {
                this.client.log("authenticate initialized failed with:" + message)
                this.err("authenticated init failed with:" + message)
            }
            return ""
        })
        this.client.addReceiveInterceptor((message: string) => {
            if (this.encryptEnable) {
                // 解密 TODO
                return message
            }
            return message
        })
        this.client.addSendInterceptor((message: string) => {
            if (this.encryptEnable) {
                // 加密 TODO
                return message
            }
            return message
        })
        this.client.paused(this.stop)
        this.client.ready(this.start)
    }

    private setPublicKey(publicKey: string) {
        this.rsa = new JSEncrypt()
        this.aes = new JSEncrypt() // TODO
        this.rsa.setPublicKey(publicKey)
    }

    private start() {
        this.client.log("authenticate start")
        this.client.send(this.authenticateMessage(), true)
    }

    private stop() {
        this.initialized = false
        this.encryptEnable = false
        this.client.log("authenticate stop")
        this._paused.forEach(fn => fn())
    }

    public error(fn: (msg: string) => void) {
        this.err = fn
    }

    public ready(fn: () => void) {
        this._ready.push(fn)
    }

    public paused(fn: () => void) {
        this._paused.push(fn)
    }

    private authenticateMessage() {
        return `${this.type}@${this.id}@${this.dataType}::${this.encryptedKey()}`
    }

    private encryptedKey() {
        if (this.rsa != null) {
            this.randEsKey()
            const str = this.rsa.encrypt(this.key + Math.floor(Date.now() / 1000))
            if (str === false) {
                this.err(`encrypted key failed: ${str}.`)
                this.client.close()
                return ""
            }
            return str
        }
        return ""
    }

    private randEsKey() {
        let len = 8
        if (this.esType === "des") {
            len = 8
        } else {
            len = parseInt(this.esType) / 8
        }
        this.key = this.rand(len)
    }

    private rand(num: number) {
        let chars = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"];
        let res = ""
        for (let i = 0; i < num; i++) {
            let id = Math.floor(Math.random() * (chars.length))
            res += chars[id]
        }
        return res
    }
}