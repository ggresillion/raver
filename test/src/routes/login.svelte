<script lang="ts">
  import { config } from '../lib/api/config.api';

  function login() {
    window.location.href = `${config.apiUrl}/auth/login`;
  }
</script>

<h1>login page</h1>
<button on:click={login}>
  Login
</button>
