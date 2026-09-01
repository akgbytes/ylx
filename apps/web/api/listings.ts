import { apiClient } from "@/lib/api-client";

export type Listing = {
  id: string;
  title: string;
  description: string;
  price: number;
  city: string;
  seller_id: string;
  category_id: string;
  category_name: string;
  created_at: string;
  updated_at: string;
};

export const listingsQueryKeys = {
  all: ["listings"] as const,
  list: () => [...listingsQueryKeys.all, "list"] as const,
};

export function getListings() {
  return apiClient<Listing[]>("/listings", { skipAuthRefresh: true });
}
