import {WebsocketClient} from "../dist/esm/index.js";
const ws = new WebsocketClient({url:"ws://127.0.0.1:28088/wss",debug:true});
ws.connect()