import { Logger } from '@nestjs/common';
import { ConnectedSocket, OnGatewayConnection, OnGatewayDisconnect, WebSocketGateway } from '@nestjs/websockets';
import { IPty, spawn } from 'node-pty';
import { Socket } from 'socket.io';

@WebSocketGateway({
    namespace: 'terminal',
    cors: { origin: ['https://redundant4u.com'] },
    transports: ['websocket'],
})
export default class TerminalGateway implements OnGatewayConnection, OnGatewayDisconnect {
    private ptys: Map<string, IPty> = new Map();

    handleConnection(@ConnectedSocket() socket: Socket) {
        Logger.log(`connect: ${socket.id}`);

        const pty = spawn('docker', ['exec', '-it', 'terminal', 'bash'], {
            name: 'xterm-color',
            cwd: process.env.HOME,
        });

        pty.onData((data) => {
            socket.emit('output', data);
        });

        const { id } = socket;
        this.ptys.set(id, pty);

        socket.on('input', (data) => {
            const pty = this.ptys.get(id);

            if (pty) {
                pty.write(data);
            }
        });
    }

    handleDisconnect(@ConnectedSocket() socket: Socket) {
        Logger.log(`disconnect: ${socket.id}`);

        const { id } = socket;
        const pty = this.ptys.get(id);

        if (pty) {
            pty.kill();
            this.ptys.delete(id);
        }
    }
}
