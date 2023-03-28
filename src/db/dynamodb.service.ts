import { DynamoDBClient, QueryCommand } from "@aws-sdk/client-dynamodb";
import { unmarshall } from "@aws-sdk/util-dynamodb";
import { Injectable } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";
import { Post } from "@post/domain/post.interface";
import { EnvConfig } from "@common/envConfig";

@Injectable()
export class DynamoDBService {
    private readonly dynamoDB: DynamoDBClient;

    constructor(private readonly configService: ConfigService<EnvConfig, true>) {
        this.dynamoDB = new DynamoDBClient({
            region: this.configService.get("AWS_REGION"),
            credentials: {
                accessKeyId: this.configService.get("AWS_ACCESS_KEY"),
                secretAccessKey: this.configService.get("AWS_SECRET_KEY"),
            },
        });
    }

    async findPosts() {
        const command = new QueryCommand({
            TableName: this.configService.get("DYNAMODB_TABLE"),
            KeyConditionExpression: "#PK = :PK",
            ExpressionAttributeNames: {
                "#PK": "PK",
            },
            ExpressionAttributeValues: {
                ":PK": {
                    S: "post",
                },
            },
        });

        const posts = await this.dynamoDB.send(command);
        const unmarshalled = posts.Items?.map((post) => unmarshall(post)) as Post[] | null;

        return unmarshalled;
    }

    async findPost(postId: string) {
        const command = new QueryCommand({
            TableName: this.configService.get("DYNAMODB_TABLE"),
            KeyConditionExpression: "#PK = :PK AND #SK = :SK",
            ExpressionAttributeNames: {
                "#PK": "PK",
                "#SK": "SK",
            },
            ExpressionAttributeValues: {
                ":PK": {
                    S: "post",
                },
                ":SK": {
                    S: postId,
                },
            },
        });

        const post = await this.dynamoDB.send(command);
        const unmarshalled = post.Items?.map((post) => unmarshall(post)) as Post[] | null;

        return unmarshalled ? unmarshalled[0] : null;
    }
}
