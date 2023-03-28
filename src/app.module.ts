import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";
import { HealthzModule } from "@health/healthz.module";
import { TerminalModule } from "./terminal/terminal.module";
import { PostModule } from "./post/post.module";
import { DynamoDBModule } from "@db/dynamodb.module";
import * as Joi from "joi";

@Module({
    imports: [
        ConfigModule.forRoot({
            isGlobal: true,
            envFilePath: [".env"],
            validationSchema: Joi.object({
                AWS_ACCESS_KEY: Joi.string().required(),
                AWS_SECRET_KEY: Joi.string().required(),
                AWS_REGION: Joi.string().required(),
                DYNAMODB_TABLE: Joi.string().required(),
            }),
        }),
        HealthzModule,
        TerminalModule,
        PostModule,
        DynamoDBModule,
    ],
})
export class AppModule {}
