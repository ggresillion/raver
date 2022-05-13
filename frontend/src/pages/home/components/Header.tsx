import React, { useEffect, useState } from 'react';
import logo from '../../../assets/icons/app_white_24dp.svg';
import { getUser } from '../../../services/auth';
import { getLatencyInMillis } from '../../../services/bot';
import { User } from '../../../services/model/user';
import './Header.scss';

export function Header() {

  const [user, setUser] = useState<User>();
  const [latency, setLatency] = useState<number>();

  useEffect(() => {

    const fetchUser = async () => {
      const user = await getUser();
      setUser(user);
    };

    const fetchLatency = async () => {
      const latency = await getLatencyInMillis();
      console.log(latency);
      setLatency(latency);
    };

    fetchUser();
    fetchLatency();
    setInterval(() => fetchLatency(), 5000);
  }, []);

  return (
    <div className="header">
      <img src={logo} alt="logo" className="logo"></img>
      <span>Discord Sound Board</span>
      <div className='user'>
        <div className='latency'>
          <div className='icon'></div>
          <div className='value'>{latency} ms</div>
        </div>
        <span><span className='thin'>Welcome, </span>{user?.username}</span>
        <img src={`https://cdn.discordapp.com/avatars/${user?.id}/${user?.avatar}.png`}></img>
      </div>
    </div>
  );
}