"use client";

import Image from "next/image";
import { useQuery } from "@tanstack/react-query";

import { Button } from "@ylx/ui/components/button";
import { IconMapPin } from "@ylx/ui/icons";

import { getListings, listingsQueryKeys, type Listing } from "@/api/listings";

const priceFormatter = new Intl.NumberFormat("en-IN", {
  style: "currency",
  currency: "INR",
  maximumFractionDigits: 0,
});

export function ListingsGrid({ maxItems }: { maxItems?: number }) {
  const listingsQuery = useQuery({
    queryKey: listingsQueryKeys.list(),
    queryFn: getListings,
  });

  if (listingsQuery.isPending) {
    return <ListingsLoadingState />;
  }

  if (listingsQuery.isError) {
    return (
      <section
        className="rounded-xl border bg-card p-6 text-card-foreground"
        role="alert"
      >
        <h2 className="font-heading text-xl">Unable to load listings</h2>
        <p className="mt-2 text-muted-foreground">
          Please check your connection and try again.
        </p>
        <Button
          className="mt-5"
          onClick={() => {
            void listingsQuery.refetch();
          }}
        >
          Try again
        </Button>
      </section>
    );
  }

  const listings =
    maxItems === undefined
      ? listingsQuery.data
      : listingsQuery.data.slice(0, maxItems);

  if (listings.length === 0) {
    return (
      <section className="rounded-xl border bg-card p-6 text-card-foreground">
        <h2 className="font-heading text-xl">No listings yet</h2>
        <p className="mt-2 text-muted-foreground">
          Check back shortly for developer-owned tech gear.
        </p>
      </section>
    );
  }

  return (
    <ul className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
      {listings.map((listing) => (
        <li key={listing.id}>
          <ListingCard listing={listing} />
        </li>
      ))}
    </ul>
  );
}

function ListingCard({ listing }: { listing: Listing }) {
  return (
    <article className="flex h-full flex-col rounded-xl border bg-card p-4 text-card-foreground shadow-sm">
      <div className="relative grid aspect-4/3 place-items-center rounded-lg bg-accent">
        <span className="absolute top-3 left-3 max-w-[calc(100%-1.5rem)] truncate rounded-full bg-background/90 px-2.5 py-1 text-xs font-semibold text-foreground shadow-xs backdrop-blur-sm">
          {listing.category_name}
        </span>
        <Image
          src="/brand/favicon-192.png"
          alt=""
          className="size-20 drop-shadow-sm"
          width={80}
          height={80}
        />
      </div>

      <div className="mt-4 flex flex-1 flex-col">
        <div className="flex items-start justify-between gap-4">
          <h2 className="min-w-0 flex-1 truncate font-heading text-lg leading-snug font-medium">
            {listing.title}
          </h2>
          <p className="shrink-0 font-semibold tabular-nums">
            {formatPrice(listing.price)}
          </p>
        </div>

        <p className="mt-3 line-clamp-2 min-h-10 text-sm leading-5 text-muted-foreground">
          {listing.description}
        </p>

        <p className="mt-auto flex items-center gap-1.5 pt-3 text-sm text-muted-foreground">
          <IconMapPin className="size-4 shrink-0" aria-hidden="true" />
          <span className="truncate">{listing.city}</span>
        </p>
      </div>
    </article>
  );
}

function ListingsLoadingState() {
  return (
    <section aria-busy="true" aria-live="polite">
      <p className="sr-only">Loading listings…</p>
      <div
        className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3"
        aria-hidden="true"
      >
        {[0, 1, 2].map((index) => (
          <div key={index} className="rounded-xl border bg-card p-4">
            <div className="aspect-4/3 animate-pulse rounded-lg bg-accent" />
            <div className="mt-4 space-y-3">
              <div className="h-5 w-3/4 animate-pulse rounded bg-accent" />
              <div className="h-4 w-full animate-pulse rounded bg-accent" />
              <div className="h-4 w-1/3 animate-pulse rounded bg-accent" />
            </div>
          </div>
        ))}
      </div>
    </section>
  );
}

function formatPrice(priceInPaise: number) {
  return priceFormatter.format(priceInPaise / 100);
}
