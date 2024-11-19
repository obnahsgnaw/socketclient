import {WebSocket} from "ws";
// 创建一个 WebSocket 连接
const socket = new WebSocket('ws://127.0.0.1:1803');

// 当连接打开时触发
socket.addEventListener('open', () => {
    console.log('WebSocket connection opened');

    // 发送消息
    socket.send('Hello from client!');
});

// 当收到消息时触发
socket.addEventListener('message', (event) => {
    console.log('Received message: ', event.data);
});

// 当连接关闭时触发
socket.addEventListener('close', () => {
    console.log('WebSocket connection closed');
});

// 当连接发生错误时触发
socket.addEventListener('error', (error) => {
    console.error('WebSocket error: ', error);
});