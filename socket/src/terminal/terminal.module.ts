import TerminalGateway from "./terminal.gateway";

import { Module } from "@nestjs/common";

@Module({
    providers: [TerminalGateway],
})
export class TerminalModule {}
