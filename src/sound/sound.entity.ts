import {Column, Entity, Index, ManyToOne, ObjectID, PrimaryGeneratedColumn} from 'typeorm';
import {v4 as uuid} from 'uuid';
import {Category} from '../category/category.entity';

@Entity()
export class Sound {

  @PrimaryGeneratedColumn()
  public id: ObjectID;

  @Column()
  public uuid: string = uuid();

  @Column()
  @Index({unique: true})
  public name: string;

  @ManyToOne(type => Category)
  public category: Category;

}
