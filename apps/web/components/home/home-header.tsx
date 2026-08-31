"use client";

import Image from "next/image";
import Link from "next/link";
import { useQuery } from "@tanstack/react-query";

import { buttonVariants } from "@ylx/ui/components/button";

import { authQueryKeys, getProfile } from "@/api/auth";
import { UserButton } from "@/components/home/user-button";

export function HomeHeader() {
  const profileQuery = useQuery({
    queryKey: authQueryKeys.me(),
    queryFn: getProfile,
    retry: false,
  });
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
            <UserButton user={profileQuery.data} />
          ) : (
            <Link
              href="/sign-in"
              className={buttonVariants({ variant: "outline", size: "sm" })}
            >
              Sign in
            </Link>
          )}
        </div>
      </div>
    </header>
  );
}
