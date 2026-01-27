import { AxiosError } from "axios";
import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

export const parseErrorMessage = (e: AxiosError | Error): string => {
  const error = e as AxiosError<{ error: string }>;

  return error.response?.data?.error ?? "";
};
