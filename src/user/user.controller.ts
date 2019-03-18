import {Controller, Get, UseGuards} from '@nestjs/common';
import {UserService} from './user.service';
import {GetUserDTO} from './dto/get-user.dto';
import {UserGuard} from '../auth/guards/user.guard';
import {AuthGuard} from '@nestjs/passport';

@Controller('users')
@UseGuards(UserGuard)
export class UserController {

  constructor(private readonly userService: UserService) {
  }

  @Get()
  @UseGuards(AuthGuard())
  public async getAllUsers(): Promise<GetUserDTO[]> {
    return await this.userService.findAll();
  }
}
