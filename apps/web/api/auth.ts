import { apiClient } from "@/lib/api-client";

import type {
  SignInValues,
  SignUpValues,
  VerificationCodeValues,
} from "@/lib/validation/auth";

export type User = {
  id: string;
  name: string;
  email: string;
  created_at?: string;
};

export const authQueryKeys = {
  all: ["auth"] as const,
  me: () => [...authQueryKeys.all, "me"] as const,
};

export type SignupResponse = {
  retry_at: string;
};

export function signUp(input: SignUpValues) {
  return apiClient<SignupResponse>("/auth/signup", {
    method: "POST",
    body: input,
    skipAuthRefresh: true,
  });
}

export function verifySignUp(input: {
  email: string;
  otp: VerificationCodeValues["otp"];
}) {
  return apiClient<User>("/auth/signup/verify", {
    method: "POST",
    body: input,
    skipAuthRefresh: true,
  });
}

export function resendSignUp(email: string) {
  return apiClient<SignupResponse>("/auth/signup/resend", {
    method: "POST",
    body: { email },
    skipAuthRefresh: true,
  });
}

export function signIn(input: SignInValues) {
  return apiClient<User>("/auth/signin", {
    method: "POST",
    body: input,
    skipAuthRefresh: true,
  });
}

export function getProfile() {
  return apiClient<User>("/auth/me");
}

export function refreshSession() {
  return apiClient<void>("/auth/refresh", {
    method: "POST",
    skipAuthRefresh: true,
  });
}

export function logout() {
  return apiClient<void>("/auth/logout", {
    method: "POST",
    skipAuthRefresh: true,
  });
}
