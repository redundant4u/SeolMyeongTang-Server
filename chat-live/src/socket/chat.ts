import { Socket } from "socket.io";
import { pubClient } from "../utils/redis";
import { query } from "../utils/db";

const chatHandlers = (io: Socket) => {
    const createChat = async (message: any) => {
        const chatCount = await pubClient.rPush("chat", message);

        if (chatCount >= 10) {
            const messages = await pubClient.lRange("chat", 0, 9);
            const values = messages.map((message) => message);
            await query(
                `INSERT INTO chat (content) VALUES ($1), ($2), ($3), ($4), ($5), ($6), ($7), ($8), ($9), ($10)`,
                values
            );
            await pubClient.lPopCount("chat", 10);
        }

        io.broadcast.emit("message", message);
    };

    return {
        createChat,
    };
};

export default chatHandlers;
