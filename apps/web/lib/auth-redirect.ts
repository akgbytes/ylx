const authPaths = ["/sign-in", "/sign-up", "/verify-sign-up"];

export function sanitizeAuthRedirect(value: string | null) {
  if (!value || !value.startsWith("/") || value.startsWith("//")) {
    return undefined;
  }

  try {
    const url = new URL(value, "http://ylx.local");
    if (url.origin !== "http://ylx.local" || isAuthPath(url.pathname)) {
      return undefined;
    }

    return `${url.pathname}${url.search}${url.hash}`;
  } catch {
    return undefined;
  }
}

export function authHref(path: string, redirectTo: string | undefined) {
  if (!redirectTo) return path;

  const params = new URLSearchParams({ redirect: redirectTo });
  return `${path}?${params.toString()}`;
}

export function authRedirectTarget(redirectTo: string | undefined) {
  return redirectTo ?? "/";
}

function isAuthPath(pathname: string) {
  return authPaths.some(
    (path) => pathname === path || pathname.startsWith(`${path}/`)
  );
}
