import {Column, Entity, PrimaryGeneratedColumn} from 'typeorm';

@Entity()
export class Category {

  @PrimaryGeneratedColumn()
  public id: number;

  @Column()
  public name: string;
}
