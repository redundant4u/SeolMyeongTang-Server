import { FindPostResponseDto } from "@post/application/dto/findPost.dto";
import { Post } from "./post.interface";

export class PostDtoService {
    public static convertPostsToFindPostResponseDto(posts: Post[] | null) {
        const result: FindPostResponseDto[] = [];

        if (!posts) {
            return [];
        }

        for (const post of posts) {
            const { title, content, SK, createdAt } = post;

            result.push(
                new FindPostResponseDto({
                    title,
                    content,
                    SK,
                    createdAt,
                })
            );
        }

        return result;
    }
}
