import { Controller, Get, HttpCode, HttpStatus } from "@nestjs/common";

@Controller("healthz")
export class HealthzController {
    @Get()
    @HttpCode(HttpStatus.OK)
    check() {
        return;
    }
}
