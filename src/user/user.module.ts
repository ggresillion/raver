import {forwardRef, Module} from '@nestjs/common';
import {UserService} from './user.service';
import {UserController} from './user.controller';
import {AuthModule} from '../auth/auth.module';

@Module({
  imports: [forwardRef(() => AuthModule)],
  providers: [UserService],
  controllers: [UserController],
  exports: [UserService],
})
export class UserModule {
}
