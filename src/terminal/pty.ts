import { IPty, spawn } from 'node-pty';
import { Socket } from 'socket.io';

export class Pty {
    private readonly socket: Socket;
    private readonly shell: string;
    private readonly pty: IPty;

    constructor(socket: Socket) {
        this.socket = socket;
        this.shell = 'docker';

        this.pty = spawn(this.shell, ['exec', '-it', 'terminal', 'rbash'], {
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
