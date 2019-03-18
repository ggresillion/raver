import {Body, Controller, Post} from '@nestjs/common';
import {AuthService} from './auth.service';
import {CredentialsDto} from './dto/credentials.dto';
import {CreateUserDto} from '../user/dto/create-user.dto';

@Controller('auth')
export class AuthController {
  constructor(private readonly authService: AuthService) {
  }

  @Post('login')
  public async login(@Body() credentials: CredentialsDto): Promise<any> {
    return await this.authService.login(credentials.email, credentials.password);
  }

  @Post('register')
  public async register(@Body()userDto: CreateUserDto) {
    return await this.authService.register(userDto);
  }
}
