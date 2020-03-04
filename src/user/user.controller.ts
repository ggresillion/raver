import {Controller, Get, UseGuards} from '@nestjs/common';
import {UserDTO} from './dto/user.dto';
import {UserGuard} from '../auth/guards/user.guard';
import {ConnectedUser} from '../auth/decorators/connected-user.decorator';

@Controller('users')
@UseGuards(UserGuard)
export class UserController {

  @Get('/me')
  public async getMe(@ConnectedUser() user): Promise<UserDTO[]> {
    return user;
  }
}
