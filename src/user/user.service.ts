import {Injectable} from '@nestjs/common';
import {Repository} from 'typeorm';
import {User} from './user.entity';
import {GetUserDTO} from './dto/get-user.dto';
import {CreateUserDto} from './dto/create-user.dto';
import {InjectRepository} from '@nestjs/typeorm';
import {compare} from 'bcrypt';

@Injectable()
export class UserService {

  constructor(@InjectRepository(User) private readonly userRepository: Repository<User>) {
  }

  public async findAll(): Promise<GetUserDTO[]> {
    return this.userRepository.find()
      .then(users => users.map(GetUserDTO.toDTO));
  }

  public async createUser(userDTO: CreateUserDto): Promise<GetUserDTO> {
    return this.userRepository.save(this.userRepository.create(userDTO))
      .then(GetUserDTO.toDTO);
  }

  public async findOneByEmail(email: string): Promise<GetUserDTO> {
    return this.userRepository.findOne({where: {email}})
      .then(GetUserDTO.toDTO);
  }

  public async findOneByEmailAndPassword(email: string, password: string): Promise<GetUserDTO> {
    const user = await this.userRepository.findOne({where: {email}});
    if (await compare(password, user.password)) {
      return GetUserDTO.toDTO(user);
    }
    return null;
  }
}
