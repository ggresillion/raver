import {createParamDecorator, UnauthorizedException} from '@nestjs/common';

export const ConnectedUser: () => ParameterDecorator = createParamDecorator((args, req) => {
  const user = req.switchToHttp().getRequest().params.user;
  if (!user) {
    throw new UnauthorizedException();
  }
  return user;
});
