import {Column, Entity, Index, ManyToOne, ObjectID, OneToOne, PrimaryGeneratedColumn, RelationId, JoinColumn} from 'typeorm';
import {v4 as uuid} from 'uuid';
import {Category} from '../../category/category.entity';
import {Image} from '../../image/entity/image.entity';

@Entity()
export class Sound {

  @PrimaryGeneratedColumn()
  public id: ObjectID;

  @Column()
  public uuid: string = uuid();

  @Column()
  public name: string;

  @ManyToOne(() => Category, {onDelete: 'SET NULL'})
  public category: Category;

  @OneToOne(() => Image)
  @JoinColumn({name: 'imageId'})
  public image: Image;

  @Column({nullable: true})
  public imageId: number | null;

  @Column({nullable: true})
  public categoryId: number | null;

  @Column()
  public guildId: string;
}
