import React from 'react';
import {
  DragDropContext, Draggable, Droppable, DropResult
} from 'react-beautiful-dnd';
import { useAppSelector } from '../../../../hooks';
import { moveInPlaylist, removeFromPlaylist } from '../../../../services/musicPlayer';
import './Playlist.scss';
import { Thumbnail } from './Thumbnail';

export function Playlist() {
  const musicPlayer = useAppSelector((state) => state.musicPlayer);
  const selectedGuild = useAppSelector((state) => state.guild.selectedGuild);

  function onDragEnd(result: DropResult) {
    if (!selectedGuild) return;
    // dropped outside the list
    if (!result.destination) {
      removeFromPlaylist(selectedGuild.id, result.source.index);
      return;
    }
    moveInPlaylist(selectedGuild.id, result.source.index, result.destination.index);
  }

  return (
    <div className="playlist-container">
      <h2 className="side-title">__ Playlist __</h2>
      <div className="table-container">
        <DragDropContext onDragEnd={onDragEnd}>
          <Droppable droppableId="droppable">
            {(provided) => (
              <table className="playlist-table" ref={provided.innerRef}>
                <thead>
                  <tr>
                    <th>#</th>
                    <th>Title</th>
                    <th>Album</th>
                  </tr>
                </thead>
                <tbody>
                  {musicPlayer.playerState?.playlist.map((item, index) => (
                    <Draggable key={index} draggableId={String(index)} index={index}>
                      {(provided) => (
                        <tr
                          ref={provided.innerRef}
                          {...provided.draggableProps}
                          {...provided.dragHandleProps}
                        >
                          <td>{index + 1}</td>
                          <td>
                            <div className="track">
                              <Thumbnail url={`https://img.youtube.com/vi/${item.id}/0.jpg`} />
                              <div className="track-info">
                                <span className="track-name">{item.title}</span>
                                <span className="track-artist">{item.artist}</span>
                              </div>
                            </div>
                          </td>
                          <td>
                            <span className="track-album">{item.album}</span>
                          </td>
                        </tr>
                      )}
                    </Draggable>
                  ))}
                  {provided.placeholder}
                </tbody>
              </table>
            )}
          </Droppable>
        </DragDropContext>
      </div>
    </div>
  );
}
