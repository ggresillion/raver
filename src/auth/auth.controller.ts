import {Body, Controller, Post} from '@nestjs/common';
import {AuthService} from './auth.service';
import {CredentialsDto} from './dto/credentials.dto';

@Controller('auth')
export class AuthController {
  constructor(private readonly authService: AuthService) {
  }

  @Post('login')
  public async createToken(@Body() credentials: CredentialsDto): Promise<any> {
    return await this.authService.signIn(credentials.email, credentials.password);
  }
}
