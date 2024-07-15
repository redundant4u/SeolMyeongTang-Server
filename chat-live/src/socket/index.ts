import { Server, Socket } from "socket.io";
import chatHandlers from "./chat";
import { createAdapter } from "@socket.io/redis-adapter";
import { pubClient, subClient } from "../utils/redis";
import { wrapSocket } from "../utils/wrap";

export const createSocket = () => {
    return new Server({
        cors: {},
        transports: ["websocket"],
        httpCompression: true,
        adapter: createAdapter(pubClient, subClient),
    });
};

export const onConnection = (socket: Socket) => {
    const chat = chatHandlers(socket);
    socket.on("create", wrapSocket(chat.createChat, socket));
};
