import { IPty, spawn } from "node-pty";
import { Socket } from "socket.io";

import { Logger } from "@nestjs/common";
import {
    ConnectedSocket,
    OnGatewayConnection,
    OnGatewayDisconnect,
    SubscribeMessage,
    WebSocketGateway,
} from "@nestjs/websockets";

@WebSocketGateway({
    namespace: "terminal",
    cors: { origin: ["https://redundant4u.com"] },
    transports: ["websocket"],
})
export default class TerminalGateway implements OnGatewayConnection, OnGatewayDisconnect {
    private ptys: Map<string, IPty> = new Map();

    handleConnection(@ConnectedSocket() socket: Socket) {
        Logger.log(`[Terminal Socket Connect]: ${socket.id}`);
    }

    handleDisconnect(@ConnectedSocket() socket: Socket) {
        Logger.log(`[Terminal Socket Disconnect]: ${socket.id}`);

        const { id } = socket;
        const pty = this.ptys.get(id);

        if (pty) {
            pty.kill();
            this.ptys.delete(id);
        }
    }

    @SubscribeMessage("init")
    init(@ConnectedSocket() socket: Socket) {
        const pty = spawn("ssh", ["terminal"], {
            name: "xterm-color",
            cwd: process.env.HOME,
        });

        pty.onData((data) => {
            socket.emit("output", data);
        });

        const { id } = socket;
        this.ptys.set(id, pty);

        socket.on("input", (data) => {
            const pty = this.ptys.get(id);

            if (pty) {
                pty.write(data);
            }
        });
    }
}
