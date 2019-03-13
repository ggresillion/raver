import {Injectable, Logger} from '@nestjs/common';
import * as fs from 'fs';

@Injectable()
export class StorageService {

  private readonly STORAGE_PATH = 'uploads/';
  private readonly logger = new Logger(StorageService.name);

  constructor() {
    fs.mkdir('uploads', () => {
        this.logger.log('Storage directories have been created');
    });
  }

  public async getFile(filename: string): Promise<Buffer> {
    return new Promise((resolve, reject) => {
      fs.readFile(this.STORAGE_PATH + filename, (err, data) => {
        return err ? reject('Failed to read a file : ' + filename) : resolve(data);
      });
    });
  }

  public async listFiles(): Promise<string[]> {
    return new Promise((resolve, reject) => {
      return fs.readdir(this.STORAGE_PATH, ((err, files) => {
        return err ? reject('Failed to read dir : ' + this.STORAGE_PATH) : resolve(files);
      }));
    });
  }

  public async saveFile(filename: string, buffer: Buffer): Promise<void> {
    return new Promise((resolve, reject) =>
      fs.writeFile(this.STORAGE_PATH + filename, buffer, (err) => {
        return err ? reject(err) : resolve();
      }));
  }
}
