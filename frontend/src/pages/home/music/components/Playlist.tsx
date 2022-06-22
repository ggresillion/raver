import React, { useState } from 'react';
import {
  DragDropContext, Draggable, Droppable, DropResult,
} from 'react-beautiful-dnd';
import deleteIcon from '../../../../assets/icons/delete_white_24dp.svg';
import playlistIcon from '../../../../assets/icons/queue_music.svg';
import arrowIcon from '../../../../assets/icons/double_arrow.svg';
import { useAppSelector } from '../../../../hooks';
import './Playlist.scss';
import {
  useGetMusicPlayerStateQuery,
  useMoveInPlaylistMutation,
  useRemoveFromPlaylistMutation,
} from '../../../../api/musicApi';

export function Playlist() {
  const [isOpen, setOpen] = useState(false);
  const selectedGuild = useAppSelector((state) => state.selectedGuild);
  const { data: playerState } = useGetMusicPlayerStateQuery(selectedGuild || '');

  const [moveInPlaylist] = useMoveInPlaylistMutation();
  const [removeFromPlaylist] = useRemoveFromPlaylistMutation();

  function onDragEnd(result: DropResult) {
    if (!selectedGuild) return;
    // dropped outside the list
    if (!result.destination) {
      return;
    }
    moveInPlaylist({ guildId: selectedGuild, from: result.source.index + 1, to: result.destination.index + 1 });
  }

  if (!playerState || !playerState.playlist || playerState.playlist.length <= 0) return <div></div>;

  const currentTrack = playerState.playlist[0];

  return (
    <>
      <div className={`playlist ${isOpen ? 'open' : ''}`}>
        <div className="handle" onClick={() => setOpen(!isOpen)}>
          <button className="icon open" style={{ backgroundImage: `url(${arrowIcon})` }}></button>
          <button className="icon closed" style={{ backgroundImage: `url(${playlistIcon})` }}></button>
        </div>
        <h2 className="title">Queue</h2>
        <h3 className="subtitle">Now playing</h3>
        <div className="track">
          <img className="track-thumbnail" src={currentTrack.thumbnail}/>
          <div className="track-info">
            <span className="track-name">{currentTrack.title}</span>
            <span className="track-artist">{currentTrack.artists.map(a => a.name).join(', ')}</span>
          </div>
        </div>
        <h3 className="subtitle">Coming next</h3>
        <DragDropContext onDragEnd={onDragEnd}>
          <Droppable droppableId="droppable">
            {(provided, snapshot) => (
              <div
                className="tracks"
                {...provided.droppableProps}
                ref={provided.innerRef}
              >
                {playerState.playlist.slice(1, playerState.playlist.length).map((track, index) => (
                  <Draggable key={track.id + index} draggableId={track.id} index={index}>
                    {(provided, snapshot) => (
                      <div
                        ref={provided.innerRef}
                        {...provided.draggableProps}
                        {...provided.dragHandleProps}
                      >
                        <div className="track">
                          <img className="track-thumbnail" src={track.thumbnail}/>
                          <div className="track-info">
                            <span className="track-name">{track.title}</span>
                            <span className="track-artist">{track.artists.map(a => a.name).join(', ')}</span>
                          </div>
                          <button className="icon"
                            style={{ backgroundImage: `url(${deleteIcon})` }}
                            onClick={() => !selectedGuild || removeFromPlaylist({
                              guildId: selectedGuild,
                              index: index + 1,
                            })}></button>
                        </div>
                      </div>
                    )}
                  </Draggable>
                ))}
                {provided.placeholder}
              </div>
            )}
          </Droppable>
        </DragDropContext>
      </div>
    </>
  );
}
