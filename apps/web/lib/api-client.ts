const apiURL = (process.env.NEXT_PUBLIC_API_URL ?? "").replace(/\/$/, "");

type ApiClientOptions = Omit<RequestInit, "body" | "credentials"> & {
  body?: unknown;
  skipAuthRefresh?: boolean;
};

type ApiSuccessResponse<T> = {
  data: T;
  meta?: unknown;
};

type ApiErrorResponse = {
  error: {
    code: string;
    message: string;
    meta?: Record<string, unknown>;
  };
};

export class ApiError extends Error {
  readonly status: number;
  readonly code: string;
  readonly meta?: Record<string, unknown>;

  constructor(
    status: number,
    code: string,
    message: string,
    meta?: Record<string, unknown>
  ) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.code = code;
    this.meta = meta;
  }
}

let refreshRequest: Promise<boolean> | undefined;

export async function apiClient<T>(
  path: `/${string}`,
  options: ApiClientOptions = {}
): Promise<T> {
  if (!apiURL) {
    throw new Error("NEXT_PUBLIC_API_URL is required");
  }

  const {
    body,
    headers: initialHeaders,
    skipAuthRefresh = false,
    ...requestOptions
  } = options;
  const headers = new Headers(initialHeaders);

  if (body !== undefined && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }

  const sendRequest = () =>
    fetch(`${apiURL}${path}`, {
      ...requestOptions,
      credentials: "include",
      headers,
      body: body === undefined ? undefined : JSON.stringify(body),
    });

  let response = await sendRequest();

  if (response.status === 401 && !skipAuthRefresh) {
    const refreshed = await refreshSession();
    if (refreshed) {
      await response.body?.cancel();
      response = await sendRequest();
    }
  }

  if (response.status === 204) {
    return undefined as T;
  }

  const payload = await readJSON(response);

  if (!response.ok) {
    if (isApiErrorResponse(payload)) {
      throw new ApiError(
        response.status,
        payload.error.code,
        payload.error.message,
        payload.error.meta
      );
    }

    throw new ApiError(response.status, "request_failed", "Request failed");
  }

  if (!isApiSuccessResponse(payload)) {
    throw new ApiError(
      response.status,
      "invalid_response",
      "API returned an invalid response"
    );
  }

  return payload.data as T;
}

async function refreshSession(): Promise<boolean> {
  if (typeof window === "undefined") {
    return false;
  }

  refreshRequest ??= fetch(`${apiURL}/auth/refresh`, {
    method: "POST",
    credentials: "include",
  })
    .then((response) => response.ok)
    .catch(() => false)
    .finally(() => {
      refreshRequest = undefined;
    });

  return refreshRequest;
}

async function readJSON(response: Response): Promise<unknown> {
  const text = await response.text();
  if (text === "") {
    return undefined;
  }

  try {
    return JSON.parse(text) as unknown;
  } catch {
    return undefined;
  }
}

function isApiSuccessResponse(
  value: unknown
): value is ApiSuccessResponse<unknown> {
  return isRecord(value) && "data" in value;
}

function isApiErrorResponse(value: unknown): value is ApiErrorResponse {
  if (!isRecord(value) || !isRecord(value.error)) {
    return false;
  }

  return (
    typeof value.error.code === "string" &&
    typeof value.error.message === "string" &&
    (value.error.meta === undefined || isRecord(value.error.meta))
  );
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null;
}
