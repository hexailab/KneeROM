<script lang="ts">
  import { getContextClient, queryStore } from "@urql/svelte";
  import { browser } from "$app/env";
  import { authStore } from "$lib/auth/";
  import type { UserObject } from "$lib/auth/";
  import { GET_SESSION_ACCOUNT_DETAILS } from "$lib/schema/mutations";

  const client = getContextClient();

  const refreshAccountStore = () => {
    const getAccount = queryStore({
      client,
      variables: { sessionId: $authStore.sessionToken },
      query: GET_SESSION_ACCOUNT_DETAILS,
    });

    getAccount.subscribe((t) => {
      $authStore = {
        onBrowser: true,
        needRefresh: false,
        loading: t.fetching,
        sessionToken: $authStore.sessionToken,
        userCache: (t.error || t.fetching) ? undefined : <UserObject>{
          name: t.data["getPractitionerAccountDetails"]["fullName"],
        }
      };
    })
  }

  authStore.subscribe((t) => {
    if (!browser || !t.onBrowser) return;

    localStorage.setItem("session_id", t.sessionToken);
    if (t.needRefresh) refreshAccountStore();
  })

  if (browser) {
    $authStore = {
      onBrowser: true,
      sessionToken: localStorage.getItem("session_id"),
      needRefresh: true,
      userCache: null,
      loading: true,
    }
  }
</script>