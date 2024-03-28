import { AppModule } from "./app.module";

import { Logger } from "@nestjs/common";
import { NestFactory } from "@nestjs/core";

async function bootstrap() {
    const app = await NestFactory.create(AppModule);

    app.enableCors({
        origin: ["https://redundant4u.com"],
    });

    const port = process.env.APP_PORT || 3000;

    app.setGlobalPrefix("v1");

    await app.listen(port);
    Logger.log(`[Port]: ${port}`);
}
bootstrap();
