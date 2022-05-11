import React from 'react';
import './ProgressBar.scss';

export function ProgressBar(props: { disabled?: boolean, progress?: number, onProgressChange?: (progress: number) => void }) {

  const progressRef = React.createRef<HTMLDivElement>();

  function onProgressChange(e: React.MouseEvent) {
    if(props.disabled) return;
    if (!props.onProgressChange) return;
    e.preventDefault();
    const rect = progressRef.current?.getBoundingClientRect();
    if (!rect) return;
    const x = e.clientX - rect.left;
    const xp = x / rect.width * 100;
    props.onProgressChange(xp);
  }

  return (
    <div className={`progress-outer ${props.disabled ? 'disabled': ''}`} onClick={onProgressChange} ref={progressRef}>
      <div className='progress-inner' style={{ minWidth: `calc(${props.progress}% - 3px)` }}></div>
      {!!props.onProgressChange && <div className='progress-button'></div>}
    </div>
  );
}