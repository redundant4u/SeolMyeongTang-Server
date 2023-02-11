import { Logger } from '@nestjs/common';
import {
    ConnectedSocket,
    OnGatewayConnection,
    OnGatewayDisconnect,
    SubscribeMessage,
    WebSocketGateway,
    WebSocketServer,
} from '@nestjs/websockets';
import { Server, Socket } from 'socket.io';

@WebSocketGateway(3002, {
    namespace: 'terminal',
    cors: { origin: ['http://localhost:3001'] },
    transports: ['websocket'],
})
export class TerminalGateway implements OnGatewayConnection, OnGatewayDisconnect {
    @WebSocketServer()
    server: Server;

    handleConnection(@ConnectedSocket() socket: Socket) {
        Logger.log(`hello ${socket}`);
    }

    handleDisconnect(@ConnectedSocket() socket: Socket) {
        Logger.log(`bye ${socket}`);
    }

    @SubscribeMessage('init')
    handleInit(@ConnectedSocket() socket: Socket) {
        Logger.log(`init ${socket}`);
    }
}
