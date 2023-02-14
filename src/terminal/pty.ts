import { IPty, spawn } from 'node-pty';
import { Socket } from 'socket.io';

export class Pty {
    socket: Socket;
    shell: string;
    pty: IPty;

    constructor(socket: Socket) {
        this.socket = socket;
        this.shell = process.platform === 'win32' ? 'cmd.exe' : 'bash';

        this.pty = spawn(this.shell, [], {
            name: 'terminal',
            cwd: process.env.HOME,
        });
        this.pty.onData((data) => this.send(data));
    }

    write(data: string) {
        this.pty.write(data);
    }

    send(data: string) {
        this.socket.emit('output', data);
    }
}
