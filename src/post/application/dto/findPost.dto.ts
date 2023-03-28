import { Post } from "@post/domain/post.interface";

export class FindPostResponseDto {
    title: string;
    content: string;
    link: string;
    createdAt: string;

    constructor(post: Post) {
        this.title = post.title;
        this.content = post.content;
        this.link = post.SK;
        this.createdAt = post.createdAt;
    }
}

export class FindPostsResponseDto {
    posts: FindPostResponseDto[] | null;

    constructor(posts: FindPostResponseDto[] | null) {
        this.posts = posts;
    }
}
