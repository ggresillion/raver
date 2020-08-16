import { Controller, Get, Query } from '@nestjs/common';
import { ImageService } from './image.service';

@Controller('image')
export class ImageController {

    constructor(private readonly imageService: ImageService) { }

    @Get(':uuid')
    public async getImageByUUID(@Query('uuid') uuid: string) {
        return this.imageService.getImageByUUID(uuid);
    }
}
