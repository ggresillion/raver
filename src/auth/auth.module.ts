import {Module} from '@nestjs/common';
import {JwtModule} from '@nestjs/jwt';
import {AuthService} from './auth.service';
import {JwtStrategy} from './jwt.strategy';
import {UserModule} from '../user/user.module';
import {PassportModule} from '@nestjs/passport';
import {AuthController} from './auth.controller';

@Module({
  imports: [
    PassportModule.register({defaultStrategy: 'jwt'}),
    JwtModule.register({
      secretOrPrivateKey: 'secretKey',
      signOptions: {
        expiresIn: 3600,
      },
    }),
    UserModule,
  ],
  providers: [AuthService, JwtStrategy],
  controllers: [AuthController],
})
export class AuthModule {
}
