import {Injectable, Logger, OnModuleInit} from '@nestjs/common';
import * as fs from 'fs';

@Injectable()
export class StorageService implements OnModuleInit {

  private readonly STORAGE_PATH = 'uploads/';
  private readonly logger = new Logger(StorageService.name);

  public onModuleInit() {
    if (!fs.existsSync(this.STORAGE_PATH)) {
      fs.mkdir(this.STORAGE_PATH, (err) => this.logger.error(err));
    }
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

  public async removeFile(filename: string): Promise<void> {
    return new Promise((resolve, reject) =>
      fs.unlink(this.STORAGE_PATH + filename, (err) => {
        return err ? reject(err) : resolve();
      }));
  }

  public getPathFromUUID(uuid: string) {
    return this.STORAGE_PATH + uuid;
  }

  public getUploadDir() {
    return this.STORAGE_PATH;
  }
}
