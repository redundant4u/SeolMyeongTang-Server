import CrdtGateway from "./crdt.gateway";

import { Module } from "@nestjs/common";

@Module({
    providers: [CrdtGateway],
})
export class CrdtModule {}
