import {Injectable} from '@nestjs/common';
import {Repository} from 'typeorm';
import {User} from './user.entity';
import {GetUserDTO} from './dto/get-user.dto';
import {InjectRepository} from '@nestjs/typeorm';

@Injectable()
export class UserService {

  constructor(@InjectRepository(User) private readonly userRepository: Repository<User>) {
  }

  public async findAll(): Promise<GetUserDTO[]> {
    return this.userRepository.find()
      .then(users => users.map(GetUserDTO.toDTO));
  }
}
