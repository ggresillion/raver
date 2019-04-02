import {Column, Entity, Index, ObjectIdColumn} from 'typeorm';

@Entity()
export class User {
  @ObjectIdColumn()
  public id: number;

  @Column()
  @Index({unique: true})
  public email: string;

  @Column()
  @Index({unique: true})
  public username: string;

  @Column()
  public password: string;
}
