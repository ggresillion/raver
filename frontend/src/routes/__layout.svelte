<script context="module" lang="ts">
  import { isAuthenticated } from '$lib/stores/authentication.store';
  import { browser } from '$app/env';
  import { get } from 'svelte/store';

  export async function load({ url }) {
    if (!browser) return {};
    if (url.pathname === '/login' && get(isAuthenticated)) return {
      status: 302,
      redirect: '/',
    };
    if (url.pathname === '/login') return {};
    if (url.pathname === '/callback') return {};
    if (get(isAuthenticated)) return {};
    return {
      status: 302,
      redirect: '/login',
    };
  }
</script>

<script>
  import '$lib/styles/index.scss';
</script>

<slot/>
