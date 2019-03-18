import {createParamDecorator, UnauthorizedException} from '@nestjs/common';

export const User: () => ParameterDecorator = createParamDecorator((args, req) => {
  const user = req.user;
  if (!user) {
    throw new UnauthorizedException();
  }
  return user;
});
