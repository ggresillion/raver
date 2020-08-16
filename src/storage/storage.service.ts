import { Injectable, Logger, OnModuleInit } from '@nestjs/common';
import * as fs from 'fs';

@Injectable()
export class StorageService implements OnModuleInit {

  private readonly STORAGE_PATH = './uploads/';
  private readonly logger = new Logger(StorageService.name);

  public onModuleInit() {
    if (!fs.existsSync(this.STORAGE_PATH)) {
      fs.mkdir(this.STORAGE_PATH, (err) => {
        if (err) {
          this.logger.error(err);
          return;
        }
        fs.mkdir(this.STORAGE_PATH + 'tmp', (err) => {
          if (err) {
            this.logger.error(err);
          }
        });
        fs.mkdir(this.STORAGE_PATH + 'sounds', (err) => {
          if (err) {
            this.logger.error(err);
          }
        });
        fs.mkdir(this.STORAGE_PATH + 'images', (err) => {
          if (err) {
            this.logger.error(err);
          }
        });
      });
    }
  }

  public async getFile(bucket: string, filename: string): Promise<Buffer> {
    return new Promise((resolve, reject) => {
      fs.readFile(`${this.STORAGE_PATH}/${bucket}/${filename}`, (err, data) => {
        return err ? reject('Failed to read a file : ' + filename) : resolve(data);
      });
    });
  }

  public async listFiles(bucket: string): Promise<string[]> {
    return new Promise((resolve, reject) => {
      return fs.readdir(`${this.STORAGE_PATH}/${bucket}`, ((err, files) => {
        return err ? reject('Failed to read dir : ' + `${this.STORAGE_PATH}/${bucket}`) : resolve(files);
      }));
    });
  }

  public async saveFile(bucket: string, filename: string, buffer: Buffer): Promise<void> {
    return new Promise((resolve, reject) =>
      fs.writeFile(`${this.STORAGE_PATH}/${bucket}/${filename}`, buffer, (err) => {
        return err ? reject(err) : resolve();
      }));
  }

  public async removeFile(bucket: string, filename: string): Promise<void> {
    return new Promise((resolve, reject) =>
      fs.unlink(`${this.STORAGE_PATH}/${bucket}/${filename}`, (err) => {
        return err ? reject(err) : resolve();
      }));
  }

  public getPathFromUUID(bucket:string, uuid: string) {
    return `./${this.STORAGE_PATH}/${bucket}/${uuid}`;
  }

  public getUploadDir(bucket: string) {
    return `${this.STORAGE_PATH}/${bucket}`;
  }
}
