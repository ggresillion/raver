<script lang="ts">
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { isAuthenticated } from '$lib/stores/authentication.store';

  onMount(() => {
    let params = new URLSearchParams($page.url.searchParams);
    if (params.get('accessToken')) {
      const accessToken = params.get('accessToken');
      const refreshToken = params.get('refreshToken');
      localStorage.setItem('accessToken', accessToken);
      localStorage.setItem('refreshToken', refreshToken);
      isAuthenticated.set(true);
      params.delete('accessToken');
      params.delete('refreshToken');
      goto(`/?${params.toString()}`);
    }
  });
</script>
