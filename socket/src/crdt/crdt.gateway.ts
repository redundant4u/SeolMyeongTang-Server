import { Logger } from "@nestjs/common";
import {
    ConnectedSocket,
    MessageBody,
    OnGatewayConnection,
    OnGatewayDisconnect,
    SubscribeMessage,
    WebSocketGateway,
} from "@nestjs/websockets";
import { Socket } from "socket.io";

type RGB = [number, number, number];
type CrdtState = [string, [peer: string, timestamp: number, value: RGB]];

@WebSocketGateway({
    namespace: "crdt",
    cors: { origin: ["https://redundant4u.com"] },
    transports: ["websocket"],
})
export default class CrdtGateway implements OnGatewayConnection, OnGatewayDisconnect {
    private count: number = 0;
    private states: CrdtState = ["", ["", 0, [0, 0, 0]]];

    handleConnection(@ConnectedSocket() socket: Socket) {
        if (this.count > 50) {
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
    write(@ConnectedSocket() socket: Socket, @MessageBody() data: CrdtState) {
        this.states = data;
        socket.broadcast.emit("merge", data);
    }

    @SubscribeMessage("clear")
    clear(@ConnectedSocket() socket: Socket) {
        this.states = ["", ["", 0, [0, 0, 0]]];

        socket.emit("clear");
        socket.broadcast.emit("clear");
    }
}
