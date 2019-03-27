import {Controller, Get, UseGuards} from '@nestjs/common';
import {UserService} from './user.service';
import {GetUserDTO} from './dto/get-user.dto';
import {UserGuard} from '../auth/guards/user.guard';
import {DiscordService} from '../discord/discord.service';
import {ConnectedUser} from '../auth/decorators/connected-user.decorator';

@Controller('users')
@UseGuards(UserGuard)
export class UserController {

  constructor(private readonly userService: UserService,
              private readonly discordService: DiscordService) {
  }

  @Get()
  @UseGuards(UserGuard)
  public async getAllUsers(): Promise<GetUserDTO[]> {
    return await this.userService.findAll();
  }

  @Get('/me')
  @UseGuards(UserGuard)
  public getMe(@ConnectedUser() user): Promise<GetUserDTO[]> {
    return user;
  }
}
