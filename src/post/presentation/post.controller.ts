import { Controller, Get, HttpCode, HttpStatus, Param } from "@nestjs/common";
import { PostService } from "@post/application/post.service";

@Controller("post")
export class PostController {
    constructor(private readonly postService: PostService) {}

    @Get()
    @HttpCode(HttpStatus.OK)
    async getPosts() {
        return await this.postService.findPosts();
    }

    @Get(":postId")
    @HttpCode(HttpStatus.OK)
    async getPost(@Param("postId") postId: string) {
        return await this.postService.findPost(postId);
    }
}
