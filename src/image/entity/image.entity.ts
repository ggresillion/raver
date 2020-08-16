import {Column, Entity, ObjectID, PrimaryGeneratedColumn} from 'typeorm';
import {v4 as uuid} from 'uuid';

@Entity()
export class Image {

  @PrimaryGeneratedColumn()
  public id: ObjectID;

  @Column()
  public uuid: string = uuid();

}
