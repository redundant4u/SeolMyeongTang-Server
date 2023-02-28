import { Logger } from '@nestjs/common';
import {
    ConnectedSocket,
    MessageBody,
    OnGatewayConnection,
    OnGatewayDisconnect,
    SubscribeMessage,
    WebSocketGateway,
} from '@nestjs/websockets';
import { Socket } from 'socket.io';
import { Pty } from './pty';

@WebSocketGateway({
    namespace: 'terminal',
    cors: { origin: ['https://redundant4u.com'] },
    transports: ['websocket'],
})
export default class TerminalGateway implements OnGatewayConnection, OnGatewayDisconnect {
    pty: Pty;

    handleConnection(@ConnectedSocket() socket: Socket) {
        Logger.log(`connect: ${socket.id}`);
    }

    handleDisconnect(@ConnectedSocket() socket: Socket) {
        Logger.log(`disconnect: ${socket.id}`);
    }

    @SubscribeMessage('init')
    init(@ConnectedSocket() socket: Socket) {
        this.pty = new Pty(socket);
    }

    @SubscribeMessage('input')
    write(_, @MessageBody() data: string) {
        this.pty.write(data);
    }
}
