import React from 'react';
import './Soundwave.scss';

export function Soundwave(props: {play: boolean}) {
  return (
    <div className={props.play ? 'play bars': 'bars'}>
      <div className='bar'></div>
      <div className='bar'></div>
      <div className='bar'></div>
      <div className='bar'></div>
      <div className='bar'></div>
      <div className='bar'></div>
      <div className='bar'></div>
      <div className='bar'></div>
      <div className='bar'></div>
      <div className='bar'></div>
    </div>
  );
}