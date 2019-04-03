import {Column, Entity, Index, ObjectID, PrimaryGeneratedColumn} from 'typeorm';
import {v4 as uuid} from 'uuid';

@Entity()
export class Sound {

  @PrimaryGeneratedColumn()
  public id: ObjectID;

  @Column()
  public uuid: string = uuid();

  @Column()
  @Index({unique: true})
  public name: string;

}
