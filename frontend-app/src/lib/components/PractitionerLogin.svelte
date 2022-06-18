<script lang="ts">
  import {getContextClient, mutationStore} from "@urql/svelte";
  import {LOGIN_TO_PRACTITIONER_ACCOUNT} from "$lib/schema/mutations";
  import {authStore} from "$lib/auth";
  import {loginFormState, errorKeyToPlaintextMap} from "./PractitionerLogin.types";

  let contextClient = getContextClient();

  const onLoginClick = () => {
    let loginMutation = mutationStore({
      client: contextClient,
      query: LOGIN_TO_PRACTITIONER_ACCOUNT,
      variables: {
        password: $loginFormState.passwordInput,
        email: $loginFormState.emailInput
      }
    })

    loginMutation.subscribe(t => {
      // If the login attempt is loading on the backend, or it's errored, update that on the state and return.
      if (t.fetching || t.error) {
        loginFormState.update(log => {
          log.isLoading = t.fetching;
          log.error = t.error;

          return log;
        });

        return;
      }

      // If the login was successful, i.e. not loading and no errors, set the new session token and refresh.
      authStore.update(auth => {
        auth.sessionToken = t.data["loginToAccount"]["sessionId"];
        auth.needRefresh = true;

        return auth;
      });
    });
  }
</script>

<input
    type="text"
    placeholder="email"
    class="border p-2"
    bind:value={$loginFormState.emailInput}
/>
<input
    type="password"
    placeholder="password"
    class="border p-2"
    bind:value={$loginFormState.passwordInput}
/>
<button class="bg-gray-200 p-2" on:click={onLoginClick}>Login</button>

{#if $loginFormState.error != null}
  <span>Error: {errorKeyToPlaintextMap[$loginFormState.error.graphQLErrors[0].message]}</span>
{/if}
{#if $loginFormState.isLoading}
  <span>Loading...</span>
{/if}