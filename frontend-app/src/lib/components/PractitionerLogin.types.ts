import type {CombinedError} from "urql";
import {writable} from "svelte/store";

export interface LoginFormState {
  error?: CombinedError,
  isLoading: boolean,
  emailInput: string,
  passwordInput: string,
}

export const errorKeyToPlaintextMap = {
  "invalid_login_input": "Invalid Username or Password",
}

export const loginFormState = writable<LoginFormState>({
  error: null,
  isLoading: false,
  emailInput: "",
  passwordInput: "",
});