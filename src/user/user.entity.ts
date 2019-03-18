import {BeforeInsert, BeforeUpdate, Column, Entity, Index, ObjectIdColumn} from 'typeorm';
import {hash} from 'bcrypt';

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

  @BeforeInsert()
  @BeforeUpdate()
  private cipherPassword() {
    if (this.password) {
      return hash(this.password, 10)
        .then(hPassword => this.password = hPassword);
    }
  }
}
