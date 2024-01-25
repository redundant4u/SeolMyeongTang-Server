import { Module } from "@nestjs/common";
import { TerminalModule } from "./terminal.module";

@Module({
    imports: [TerminalModule],
})
export class AppModule {}
