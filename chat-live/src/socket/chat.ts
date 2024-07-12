import { Socket } from "socket.io";

const chatHandlers = (io: Socket) => {
    const createChat = (message: any) => {
        io.broadcast.emit("chat:message", message);
    };

    return {
        createChat,
    };
};

export default chatHandlers;
