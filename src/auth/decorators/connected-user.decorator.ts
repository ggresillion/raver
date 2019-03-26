import {createParamDecorator, UnauthorizedException} from '@nestjs/common';

export const ConnectedUser: () => ParameterDecorator = createParamDecorator((args, req) => {
  const user = req.params.user;
  if (!user) {
    throw new UnauthorizedException();
  }
  return user;
});
