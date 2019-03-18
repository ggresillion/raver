import {ExecutionContext, Injectable} from '@nestjs/common';
import {AuthGuard} from '@nestjs/passport';

@Injectable()
export class UserGuard extends AuthGuard('jwt') {

  public canActivate(context: ExecutionContext) {
    return super.canActivate(context);
  }
}
