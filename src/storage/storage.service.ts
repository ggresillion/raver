import {Injectable, Logger} from '@nestjs/common';
import * as fs from 'fs';

@Injectable()
export class StorageService {

  private readonly STORAGE_PATH = 'uploads/';
  private readonly logger = new Logger(StorageService.name);

  constructor() {
    fs.mkdir('uploads', () => {
      fs.mkdir('uploads/sounds', () => {
        this.logger.log('Storage directories have been created');
      });
    });
  }

  public getFile(filename: string): Promise<Buffer> {
    return new Promise((resolve, reject) => {
      fs.readFile(this.STORAGE_PATH + filename, (err, data) => {
        if (err) {
          reject('Failed to read a file : ' + filename);
        }
        resolve(data);
      });
    });
  }

  public listFiles(dirname: string): Promise<string[]> {
    return new Promise((resolve, reject) => {
      return fs.readdir(this.STORAGE_PATH + '/' + dirname, ((err, files) => {
        if (err) {
          reject('Failed to read dir : ' + dirname);
        }
        resolve(files);
      }));
    });
  }
}
