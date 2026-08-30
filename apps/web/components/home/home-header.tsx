"use client";

import Image from "next/image";
import Link from "next/link";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { Button } from "@ylx/ui/components/button";

import { authQueryKeys, getProfile, logout } from "@/api/auth";
import { ApiError } from "@/lib/api-client";

export function HomeHeader() {
  const queryClient = useQueryClient();
  const profileQuery = useQuery({
    queryKey: authQueryKeys.me(),
    queryFn: getProfile,
    retry: false,
  });
  const logoutMutation = useMutation({
    mutationFn: logout,
    onSuccess: () => {
      queryClient.removeQueries({ queryKey: authQueryKeys.all });
      window.location.reload();
    },
  });
  const isUnauthenticated =
    profileQuery.error instanceof ApiError && profileQuery.error.status === 401;
  const hasProfileError = profileQuery.isError && !isUnauthenticated;

  return (
    <header className="border-b bg-card/80 px-6 py-4 backdrop-blur sm:px-10">
      <div className="mx-auto flex max-w-5xl items-center gap-3">
        <Link href="/" className="flex items-center gap-3">
          <Image
            src="/brand/favicon-64.png"
            alt="YLX"
            className="size-10"
            width="40"
            height="40"
          />
          <div>
            <p className="font-heading text-xl font-bold leading-none">YLX</p>
            <p className="mt-1 text-xs text-muted-foreground">
              The Marketplace for Developers
            </p>
          </div>
        </Link>

        <div className="ml-auto flex items-center gap-3">
          {profileQuery.data ? (
            <>
              <span className="hidden text-sm text-muted-foreground sm:inline">
                {profileQuery.data.name}
              </span>
              <Button
                variant="outline"
                size="sm"
                onClick={() => logoutMutation.mutate()}
                disabled={logoutMutation.isPending}
              >
                {logoutMutation.isPending ? "Logging out…" : "Log out"}
              </Button>
              {logoutMutation.isError && (
                <p className="text-sm text-destructive" role="alert">
                  {logoutMutation.error.message}
                </p>
              )}
            </>
          ) : hasProfileError ? (
            <div className="flex items-center gap-3">
              <p className="text-sm text-destructive" role="alert">
                Unable to load your session.
              </p>
              <Button
                variant="outline"
                size="sm"
                onClick={() => profileQuery.refetch()}
              >
                Try again
              </Button>
            </div>
          ) : (
            <Link
              href="/sign-in"
              className="inline-flex h-8 items-center justify-center rounded-md border bg-background px-3 text-xs font-medium shadow-xs transition-colors hover:bg-accent hover:text-accent-foreground"
            >
              Sign in
            </Link>
          )}
        </div>
      </div>
    </header>
  );
}
