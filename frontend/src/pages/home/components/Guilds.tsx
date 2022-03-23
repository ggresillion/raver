import React, { useEffect } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { Loader } from '../../../components/Loader';
import { joinGuild } from '../../../services/guilds';
import { HttpClient } from '../../../services/http';
import { Guild } from '../../../services/model/guild';
import { setGuild, setGuilds } from '../../../slices/guild';
import { RootState } from '../../../store';
import './Guilds.scss';

export function Guilds() {

  const dispatch = useDispatch();
  const { selectedGuild, guilds } = useSelector((state: RootState) => state.guild);

  useEffect(() => {
    const getGuilds = async () => {
      const http = new HttpClient();
      const res = await http.get<Guild[]>('/guilds');
      dispatch(setGuilds({ guilds: res }));

      const guildId = localStorage.getItem('selectedGuildId');
      if (guildId) {
        const guild = res?.find(g => g.id === guildId);
        if (guild) {
          selectGuild(guild);
          return;
        }
      }
      dispatch(setGuild({ selectedGuild: undefined }));
    };
    getGuilds();
  }, []);

  async function selectGuild(guild: Guild) {
    try {
      await joinGuild(guild);
      dispatch(setGuild({ selectedGuild: guild }));
      localStorage.setItem('selectedGuildId', guild.id);

    } catch (e) {
      console.log(e);
    }
  }

  if (!guilds) {
    return (
      <div className="guilds">
        <Loader />
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
          <div className={g.id === selectedGuild?.id ? 'guild selected' : 'guild'} key={g.id} style={style} title={g.name} onClick={() => selectGuild(g)}>
            {
              !g.icon ?
                <span>{g.name.charAt(0).toUpperCase()}</span> :
                <img className="guild-icon" src={`https://cdn.discordapp.com/icons/${g.id}/${g.icon}`} />
            }
          </div>
        );
      })}
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