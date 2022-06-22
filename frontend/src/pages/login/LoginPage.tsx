import React from 'react';
import { config } from '../../config';

export function LoginPage() {

  function login() {
    window.location.href = `${config.apiUrl}auth/login`;
  }

  return <div>
    <div>Please login to discord</div>
    <button onClick={login}>Login</button>
  </div>;
}
