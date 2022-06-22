import React from 'react';
import logo from '../../../assets/icons/app_white_24dp.svg';
import './Header.scss';
import { useGetUserQuery } from '../../../api/authAPI';
import { useGetLatencyInMillisQuery } from '../../../api/botAPI';

export function Header() {

  const {
    data: user,
  } = useGetUserQuery();

  const {
    data: latency,
  } = useGetLatencyInMillisQuery();

  return (
    <div className="header">
      <img src={logo} alt="logo" className="logo"></img>
      <span>Discord Sound Board</span>
      <div className="user">
        <div className="latency">
          <div className="icon"></div>
          <div className="value">{latency} ms</div>
        </div>
        <span><span className="thin">Welcome, </span>{user?.username}</span>
        <img src={`https://cdn.discordapp.com/avatars/${user?.id}/${user?.avatar}.png`}></img>
      </div>
    </div>
  );
}
