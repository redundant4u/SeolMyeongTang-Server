import { IsNotEmpty, IsUUID } from 'class-validator';

export class GetPageRequestDto {
    // @IsUUID()
    // @IsNotEmpty()
    pageId: string;
}
