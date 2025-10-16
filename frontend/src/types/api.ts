export type TApiResponse<T = unknown> = {
  success: boolean;
  data?: T;
  errors?: unknown;
};
