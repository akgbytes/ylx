import { type ClassValue, clsx } from "clsx";
import { twMerge } from "tailwind-merge";

/** Combines conditional classes and resolves conflicting Tailwind utilities. */
export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
