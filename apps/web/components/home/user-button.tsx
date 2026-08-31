"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";

import { Avatar, AvatarFallback, AvatarImage } from "@ylx/ui/components/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@ylx/ui/components/dropdown-menu";
import { IconLogout } from "@ylx/ui/icons";

import { authQueryKeys, logout, type User } from "@/api/auth";

const profileImageURL =
  "https://res.cloudinary.com/dmnh10etf/image/upload/v1750270944/default_epnleu.png";

export function UserButton({ user }: { user: User }) {
  const queryClient = useQueryClient();
  const logoutMutation = useMutation({
    mutationFn: logout,
    onSuccess: () => {
      queryClient.removeQueries({ queryKey: authQueryKeys.all });
      window.location.reload();
    },
  });

  return (
    <div className="flex items-center gap-3">
      <DropdownMenu>
        <DropdownMenuTrigger aria-label={`Open ${user.name}'s account menu`}>
          <Avatar className="size-9">
            <AvatarImage src={profileImageURL} alt="" />
            <AvatarFallback>{getUserInitials(user.name)}</AvatarFallback>
          </Avatar>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" className="w-64">
          <DropdownMenuGroup>
            <DropdownMenuLabel className="p-1 font-normal">
              <UserIdentity user={user} />
            </DropdownMenuLabel>
          </DropdownMenuGroup>
          <DropdownMenuSeparator />
          <DropdownMenuItem
            disabled={logoutMutation.isPending}
            onClick={() => logoutMutation.mutate()}
          >
            <IconLogout aria-hidden="true" />
            {logoutMutation.isPending ? "Logging out…" : "Log out"}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      {logoutMutation.isError && (
        <p className="text-sm text-destructive" role="alert">
          {logoutMutation.error.message}
        </p>
      )}
    </div>
  );
}

function UserIdentity({ user }: { user: User }) {
  return (
    <div className="flex items-center gap-2">
      <Avatar className="size-9 rounded-lg">
        <AvatarImage src={profileImageURL} alt="" />
        <AvatarFallback className="rounded-lg">
          {getUserInitials(user.name)}
        </AvatarFallback>
      </Avatar>
      <span className="grid min-w-0 flex-1 text-left text-sm leading-tight">
        <span className="truncate font-medium">{user.name}</span>
        <span className="truncate text-xs text-muted-foreground">
          {user.email}
        </span>
      </span>
    </div>
  );
}

function getUserInitials(name: string) {
  return (
    name
      .split(" ")
      .filter(Boolean)
      .slice(0, 2)
      .map((part) => part[0])
      .join("")
      .toUpperCase() || "YLX"
  );
}
