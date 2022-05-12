import React, { useEffect, useState } from 'react';
import logo from '../../../assets/icons/app_white_24dp.svg';
import { getUser } from '../../../services/auth';
import { User } from '../../../services/model/user';
import './Header.scss';

export function Header() {

  const [user, setUser] = useState<User>();

  useEffect(() => {

    const fetchUser = async () => {
      const user = await getUser();
      setUser(user);
    };

    fetchUser();
  }, []);

  return (
    <div className="header">
      <img src={logo} alt="logo" className="logo"></img>
      <span>Discord Sound Board</span>
      <div className='user'>
        <span><span className='thin'>Welcome, </span>{user?.username}</span>
        <img src={`https://cdn.discordapp.com/avatars/${user?.id}/${user?.avatar}.png`}></img>
      </div>
    </div>
  );
}