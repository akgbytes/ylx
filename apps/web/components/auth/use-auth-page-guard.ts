"use client";

import { useEffect } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { useQuery } from "@tanstack/react-query";

import { authQueryKeys, getProfile } from "@/api/auth";
import { authRedirectTarget, sanitizeAuthRedirect } from "@/lib/auth-redirect";

export function useAuthPageGuard() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const redirectTo = sanitizeAuthRedirect(searchParams.get("redirect"));
  const profileQuery = useQuery({
    queryKey: authQueryKeys.me(),
    queryFn: getProfile,
    retry: false,
  });

  useEffect(() => {
    if (profileQuery.data) {
      router.replace(authRedirectTarget(redirectTo));
    }
  }, [profileQuery.data, redirectTo, router]);

  return {
    isRedirecting: Boolean(profileQuery.data),
    redirectTo,
  };
}
