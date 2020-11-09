import { Injectable } from '@nestjs/common';
import { StorageService } from '../storage/storage.service';
import { Bucket } from '../storage/bucket.enum';

@Injectable()
export class ImageService {

    constructor(private readonly storageService: StorageService) { }

    public getImageByUUID(uuid: string) {
        return this.storageService.getFile(Bucket.IMAGES, uuid);
    }
}
