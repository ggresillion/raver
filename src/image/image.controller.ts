import { Controller, Get, Param, Res } from '@nestjs/common';
import { ImageService } from './image.service';
import { Response } from 'express';

@Controller('image')
export class ImageController {

    constructor(private readonly imageService: ImageService) { }

    @Get(':uuid')
    public async getImageByUUID(@Param('uuid') uuid: string, @Res() res: Response) {
        this.imageService.getImageByUUID(uuid).pipe(res);
    }
}
