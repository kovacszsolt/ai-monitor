export type ApiError = { code: string; message: string };

export type ApiResult<T> = {
  ok: boolean;
  data?: T;
  error?: ApiError;
};

export function unwrap<T>(result: ApiResult<T>): T {
  if (!result.ok) {
    throw new Error(result.error?.message ?? "Request failed");
  }
  if (result.data === undefined) {
    throw new Error(result.error?.message ?? "Request failed");
  }
  return result.data;
}
