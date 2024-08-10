import { Server, Socket } from "socket.io";
import chatHandlers from "./chat";
import { createAdapter } from "@socket.io/redis-adapter";
import { pubClient, subClient } from "../utils/redis";
import { wrapSocket } from "../utils/wrap";

export const createSocket = () => {
    return new Server({
        cors: {
            origin: ["https://redundant4u.com"],
            methods: ["GET", "OPTIONS"],
            allowedHeaders: ["Origin", "Authorization", "Content-Type"],
        },
        transports: ["websocket"],
        httpCompression: true,
        adapter: createAdapter(pubClient, subClient),
    });
};

export const onConnection = (socket: Socket) => {
    console.log(socket.id);
    const chat = chatHandlers(socket);
    socket.on("create", wrapSocket(chat.createChat, socket));
};
