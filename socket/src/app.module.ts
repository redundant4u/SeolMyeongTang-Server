import { CrdtModule } from "./crdt/crdt.module";
import { TerminalModule } from "./terminal/terminal.module";

import { Module } from "@nestjs/common";

@Module({
    imports: [TerminalModule, CrdtModule],
})
export class AppModule {}
