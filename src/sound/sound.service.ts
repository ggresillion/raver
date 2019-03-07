import {Injectable} from '@nestjs/common';
import {StorageService} from '../storage/storage.service';

@Injectable()
export class SoundService {

  constructor(
    private readonly storageService: StorageService) {
  }

  public async getSounds(): Promise<string[]> {
    return this.storageService.listFiles('sounds');
  }
}
