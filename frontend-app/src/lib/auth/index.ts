import { writable } from "svelte/store";

// This can potentially hold much more information on the user.
export interface UserObject {
  name: string;
}

export interface AuthenticationStore {
  needRefresh: boolean,
  loading: boolean,
  onBrowser: boolean,
  userCache?: UserObject,
  sessionToken?: string,
}

export let authStore = writable<AuthenticationStore>({
  onBrowser: false,
  needRefresh: false,
  loading: true,
  sessionToken: undefined,
  userCache: undefined,
})

export const logOutAccount = () => {
  authStore.set({
    sessionToken: "",
    userCache: undefined,
    loading: false,
    needRefresh: true,
    onBrowser: true,
  });
}