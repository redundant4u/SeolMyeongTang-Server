import { DynamoDBService } from "@db/dynamodb.service";
import { Module } from "@nestjs/common";
import { ConfigModule } from "@nestjs/config";
import { PostController } from "./presentation/post.controller";
import { PostService } from "./application/post.service";

@Module({
    imports: [ConfigModule],
    providers: [PostService, DynamoDBService],
    controllers: [PostController],
})
export class PostModule {}
