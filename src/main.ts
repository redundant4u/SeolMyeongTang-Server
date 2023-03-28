import { Logger, ValidationPipe } from "@nestjs/common";
import { NestFactory } from "@nestjs/core";
import { AppModule } from "./app.module";

async function bootstrap() {
    const app = await NestFactory.create(AppModule);

    app.useGlobalPipes(new ValidationPipe({ transform: true, stopAtFirstError: true }));
    app.enableCors({
        origin: ["https://redundant4u.com"],
    });

    const port = process.env.APP_PORT || 3000;

    app.setGlobalPrefix("v1");

    await app.listen(port);
    Logger.log(`[Port]: ${port}`);
}
bootstrap();
