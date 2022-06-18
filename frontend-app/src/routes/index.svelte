<script lang="ts">
  import { goto } from '$app/navigation';
  import { logOutAccount } from "$lib/auth";

  let number = 0;
</script>

<h1 class="text-xl">Hello.</h1>
<p>Press the blue button to make its number go up:</p>
<button class="c-button" on:click={() => { number++; }}>
  Number: {number}
</button>
<button class="c-button c-button--secondary" on:click={() => { goto("/another-page/"); }}>
  Or, click here to go to another page.
</button>
<button class="c-button c-button--secondary" on:click={logOutAccount}>
  Log Out
</button>