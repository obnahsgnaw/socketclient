// import {Client} from "./client.js";
// 创建一个 WebSocket 连接
// const conn = new Client('ws://127.0.0.1:28088/wss');

// conn.Start()
import rq from  'proto/auth.js'
let pkg = rq()
pkg.setToken("xcx")

console.log(pkg)
console.log(pkg.serializeBinary())