import { QueryClient } from "@tanstack/react-query";

import { ApiError } from "@/lib/api-client";

export function createQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
        retry: (failureCount, error) => {
          if (
            error instanceof ApiError &&
            error.status >= 400 &&
            error.status < 500
          ) {
            return false;
          }

          return failureCount < 2;
        },
        staleTime: 30_000,
      },
      mutations: {
        retry: false,
      },
    },
  });
}
