import { Socket } from "socket.io";
import { pubClient } from "../utils/redis";
import { wrapChatQuery } from "../utils/wrap";

type Chat = {
    id: number;
    content: string;
    createdAt: string;
};

const chatHandlers = (socket: Socket) => {
    const createChat = async (chat: string) => {
        /*
        const chatCount = await pubClient.rPush("chat", chat);

        if (chatCount >= 10) {
            const chats = await pubClient.lRange("chat", 0, 9);
            const values = chats.map((chat) => chat);
            await query(
                `INSERT INTO chat (content) VALUES ($1), ($2), ($3), ($4), ($5), ($6), ($7), ($8), ($9), ($10)`,
                values
            );
            await pubClient.lPopCount("chat", 10);
        }
        */

        await pubClient.rPush("chat", chat);
        const mqChat = await pubClient.lRange("chat", 0, 1).then((chats) => chats[0]);
        await pubClient.lPopCount("chat", 1);

        const newChat = await wrapChatQuery<Chat>("INSERT INTO chats (content) VALUES ($1) RETURNING *", [mqChat]).then(
            (result) => result[0]
        );

        const data = { content: newChat.content, createdAt: newChat.createdAt };
        socket.emit("chat", data);
        socket.broadcast.emit("chat", data);
    };

    return {
        createChat,
    };
};

export default chatHandlers;
