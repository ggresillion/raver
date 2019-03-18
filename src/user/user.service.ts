import {Injectable} from '@nestjs/common';
import {Repository} from 'typeorm';
import {User} from './user.entity';
import {GetUserDTO} from './dto/get-user.dto';
import {CreateUserDTO} from './dto/create-user.dto';
import {InjectRepository} from '@nestjs/typeorm';

@Injectable()
export class UserService {

  constructor(@InjectRepository(User) private readonly userRepository: Repository<User>) {
  }

  public async findAll(): Promise<GetUserDTO[]> {
    return this.userRepository.find()
      .then(users => users.map(GetUserDTO.toDTO));
  }

  public async createUser(userDTO: CreateUserDTO): Promise<GetUserDTO> {
    return this.userRepository.save(this.userRepository.create(userDTO))
      .then(GetUserDTO.toDTO);
  }

  public async findOneByEmail(email: string) {
    return this.userRepository.findOne({where: {email}})
      .then(GetUserDTO.toDTO);
  }
}
