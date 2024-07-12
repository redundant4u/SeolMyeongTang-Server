import { Server, Socket } from "socket.io";
import chatHandlers from "./chat";
import { createAdapter } from "@socket.io/redis-adapter";
import { pubClient, subClient } from "../utils/redis";

export const createSocket = () => {
    return new Server({
        cors: {},
        transports: ["websocket"],
        httpCompression: true,
        adapter: createAdapter(pubClient, subClient),
    });
};

export const onConnection = (io: Socket) => {
    const chat = chatHandlers(io);

    io.on("create", chat.createChat);
};
