<script context="module" lang="ts">
  import { isAuthenticated } from '$lib/stores/authentication.store';

  export function load({ url }) {
    if (url.pathname === '/login') {
      return {};
    }
    if (isAuthenticated) {
      return {};
    }
    return {
      status: 302,
      redirect: '/login',
    };
  }
</script>

<script>
  import { onMount } from 'svelte';
  import '$lib/styles/index.scss';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';

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
      goto(`?${params.toString()}`);
    }
  });
</script>

<slot/>
