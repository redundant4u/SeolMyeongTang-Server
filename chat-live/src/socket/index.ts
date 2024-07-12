import * as http from "http";
import * as socketio from "socket.io";
import chatHandlers from "./chat";

export const createSocket = (server: http.Server) => {
    return new socketio.Server(server, {
        cors: {},
        transports: ["websocket"],
        httpCompression: true,
    });
};

export const onConnection = (io: socketio.Socket) => {
    const chat = chatHandlers(io);

    io.on("chat:create", chat.createChat);
};
