import {Column, Entity, Index, ObjectID, ObjectIdColumn} from 'typeorm';
import {v4 as uuid} from 'uuid';

@Entity()
export class Sound {

  @ObjectIdColumn()
  public id: ObjectID;

  @Column()
  public uuid: string = uuid();

  @Column()
  @Index({unique: true})
  public name: string;

  @Column()
  public category: string;

}
