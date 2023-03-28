import { Injectable, NotFoundException } from "@nestjs/common";
import { DynamoDBService } from "@db/dynamodb.service";
import { FindPostResponseDto, FindPostsResponseDto } from "./dto/findPost.dto";
import { PostDtoService } from "@post/domain/postDto.service";
import { NotFoundPost } from "./dto/error.dto";

@Injectable()
export class PostService {
    constructor(private readonly dynamoDBService: DynamoDBService) {}

    async findPosts() {
        const posts = await this.dynamoDBService.findPosts();
        const convertedDto = PostDtoService.convertPostsToFindPostResponseDto(posts);

        return new FindPostsResponseDto(convertedDto);
    }

    async findPost(postId: string) {
        const post = await this.dynamoDBService.findPost(postId);

        if (!post) {
            throw new NotFoundException(new NotFoundPost());
        }

        return new FindPostResponseDto(post);
    }
}
