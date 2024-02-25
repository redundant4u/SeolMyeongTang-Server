import { Module } from "@nestjs/common";
import { TerminalModule } from "./terminal/terminal.module";
import { CrdtModule } from "./crdt/crdt.module";

@Module({
    imports: [TerminalModule, CrdtModule],
})
export class AppModule {}
