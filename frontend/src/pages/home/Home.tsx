import React from 'react';
import { Guilds } from './components/Guilds';
import { Header } from './components/Header';
import './Home.scss';
import { Music } from './music/Music';

export function Home() {

  return (
    <div>
      <Header />
      <Guilds />

      <div className="main">
        <Music />
      </div>
    </div>
  );
}