<script>
  import { getUser } from '$lib/api/user.api.ts';
  import { onMount } from 'svelte';
  import appIcon from '$lib/assets/icons/app_white_24dp.svg';

  const latency = 1;
  let user;

  onMount(async () => {
    user = await getUser();
  });
</script>

<style>
    .header {
        width: 100%;
        height: 100%;
        display: flex;
        align-items: center;
        gap: 16px;
        padding: 0 16px;
        color: white;
        background: #03a9f4;
        background: linear-gradient(90deg, rgba(3, 169, 244, 1) 35%, rgba(109, 5, 255, 0.6194852941176471) 100%);
        box-sizing: border-box;
    }

    .logo {
        height: calc(100% - 12px);
    }

    .latency {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 6px;
    }

    .icon {
        width: 8px;
        height: 8px;
        background-color: rgb(0, 199, 0);
        border-radius: 99px;
    }

    .value {
        font-family: "Roboto-Thin", sans-serif;
        font-size: 11px;
    }

    .user {
        display: flex;
        align-items: center;
        gap: 16px;
        margin-left: auto;
    }

    img {
        border-radius: 99px;
        height: 42px;
        width: 42px;
    }

    .thin {
        font-family: "Roboto-Thin", sans-serif;
    }
</style>

<div class="header">
  <img src={appIcon} alt="logo" class="logo"/>
  <span>Discord Sound Board</span>
  <div class="user">
    <div class="latency">
      <div class="icon"></div>
      <div class="value">{latency} ms</div>
    </div>
    <span><span class="thin">Welcome, </span>{user?.username}</span>
    <img src={`https://cdn.discordapp.com/avatars/${user?.id}/${user?.avatar}.png`} alt="avatar"/>
  </div>
</div>
