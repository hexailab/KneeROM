<script lang="ts">
  import '../app.scss';
  import Header from '$lib/components/Header.svelte';
  import BackendWrapper from "$lib/components/BackendWrapper.svelte";
  import PractitionerLogin from "$lib/components/PractitionerLogin.svelte";
  import { authStore } from "$lib/auth";
  import { initContextClient } from "@urql/svelte";

  const GRAPHQL_API_HREF = 'http://localhost:8080/graphql';
  initContextClient({ url: GRAPHQL_API_HREF });
</script>


<BackendWrapper />
<main>
  <Header name={$authStore.userCache?.name} />
  <section class="c-main">
    {#if $authStore.loading}
      <p class="text-gray-600">Loading...</p>
    {:else if $authStore.userCache}
      <slot />
    {:else}
      <PractitionerLogin />
    {/if}
  </section>
</main>