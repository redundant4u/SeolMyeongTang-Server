import { Socket } from "socket.io";
import { query } from "./db";
import { QueryResultRow } from "pg";
import { pubClient } from "./redis";

export const wrapSocket = (fn: Function, socket: Socket) => {
    return (...args: any[]) => {
        const fnReturn = fn(...args);
        if (fnReturn instanceof Promise) {
            fnReturn.catch((err: any) => {
                console.error("[ERROR] socket event error ", err);
                socket.emit("error", {});
            });
        }
    };
};

export const wrapChatQuery = async <T extends QueryResultRow>(q: string, values: string[]) => {
    return query<T>(q, values).catch((err) => {
        values.map(async (chat) => await pubClient.rPush("chat", chat));

        console.error("[ERROR] query failed ", err);
        throw new Error("ERROR");
    });
};
