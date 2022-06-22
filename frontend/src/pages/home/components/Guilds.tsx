import React from 'react';
import { useDispatch } from 'react-redux';
import { Loader } from '../../../components/Loader';
import { Guild } from '../../../api/model/guild';
import addIcon from '../../../assets/icons/add_white_24dp.svg';
import './Guilds.scss';
import { useGetGuildsQuery } from '../../../api/guildAPI';
import { useAppSelector } from '../../../hooks';
import { setSelectedGuild } from '../../../slices/selectedGuildSlice';

export function addGuild(): void {
  window.open('http://localhost:8080/api/bot/guilds/add');
}

export function Guilds() {

  const { selectedGuild } = useAppSelector(state => state);
  const dispatch = useDispatch();
  const { data: guilds } = useGetGuildsQuery();

  async function selectGuild(guild: Guild) {
    try {
      dispatch(setSelectedGuild(guild.id));

    } catch (e) {
      console.log(e);
    }
  }

  if (!guilds) {
    return (
      <div className="guilds">
        <Loader/>
      </div>
    );
  }

  return (
    <div className="guilds">
      {guilds.map(g => {
        const style = {
          backgroundColor: stringToRGB(g.name),
        };
        return (
          <div className={g.id === selectedGuild ? 'guild selected' : 'guild'}
            key={g.id}
            style={style}
            title={g.name}
            onClick={() => selectGuild(g)}>
            {
              !g.icon ?
                <span>{g.name.charAt(0).toUpperCase()}</span> :
                <img className="guild-icon" src={`https://cdn.discordapp.com/icons/${g.id}/${g.icon}`}/>
            }
          </div>
        );
      })}
      <div className="guild" key="add" title="Add to guild" onClick={() => addGuild()}>
        <img className="guild-icon" src={addIcon}/>
      </div>
    </div>
  );
}

function stringToRGB(str: string): string {
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }
  const c = (hash & 0x00FFFFFF)
    .toString(16)
    .toUpperCase();
  const color = '00000'.substring(0, 6 - c.length) + c;
  return '#' + color;
}
