import { Module } from "@nestjs/common";
import CrdtGateway from "./crdt.gateway";

@Module({
    providers: [CrdtGateway],
})
export class CrdtModule {}
