import { Socket } from "socket.io";

import { Logger } from "@nestjs/common";
import {
    ConnectedSocket,
    MessageBody,
    OnGatewayConnection,
    OnGatewayDisconnect,
    SubscribeMessage,
    WebSocketGateway,
} from "@nestjs/websockets";

type CrdtState = {
    colors: string[];
    data: Array<[timestamp: number, colorIndex: number]>;
};

const w = 64,
    h = 64;

@WebSocketGateway({
    namespace: "crdt",
    cors: { origin: ["https://redundant4u.com"] },
    transports: ["websocket"],
})
export default class CrdtGateway implements OnGatewayConnection, OnGatewayDisconnect {
    private count: number = 0;
    private states: CrdtState = {
        colors: [],
        data: [],
    };

    handleConnection(@ConnectedSocket() socket: Socket) {
        if (this.count > 100) {
            Logger.error("[CRDT Socket Connect]: Too many connections");
            socket.disconnect();
            return;
        }

        this.count++;

        Logger.log(`[CRDT Socket Connect]: ${socket.id} - ${this.count}`);
        socket.emit("init", this.states);
    }

    handleDisconnect(@ConnectedSocket() socket: Socket) {
        this.count--;
        Logger.log(`[CRDT Socket Disconnect]: ${socket.id}`);
    }

    @SubscribeMessage("write")
    write(@ConnectedSocket() socket: Socket, @MessageBody() state: CrdtState) {
        if (state.data.length > w * h) {
            Logger.error("[CRDT Socket Write]: Invalid data");
            return;
        } else if (state.data.length > 0) {
            for (const [timestamp, colorIndex] of state.data) {
                if (timestamp > Number.MAX_SAFE_INTEGER || colorIndex > Number.MAX_SAFE_INTEGER) {
                    Logger.error("[CRDT Socket Write]: Invalid data");
                    return;
                }
            }
        }

        this.states.colors.map((color) => {
            if (color.length !== 6) {
                Logger.error("[CRDT Socket Write]: Invalid state.colors");
                return;
            }
        });

        this.states = state;
        socket.broadcast.emit("merge", state);
    }

    @SubscribeMessage("clear")
    clear(@ConnectedSocket() socket: Socket) {
        this.states = {
            colors: [],
            data: [],
        };

        socket.emit("clear");
        socket.broadcast.emit("clear");
    }
}
