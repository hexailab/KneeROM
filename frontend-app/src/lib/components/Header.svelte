<script lang="ts">
  import './Header.scss';
  import { goto } from "$app/navigation";

  export let name: string;
</script>

<header class="c-header">
  <h1 class="c-header__brand-title" on:click={() => { goto("/"); }}>
    Range of Motion Prototype
  </h1>

  <div class="flex-grow"></div>

  {#if name === undefined}
    <p>Not logged in.</p>
  {:else}
    <p>Logged in as {name}.</p>
  {/if}
</header>