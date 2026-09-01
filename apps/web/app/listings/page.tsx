import type { Metadata } from "next";

import { HomeHeader } from "@/components/home/home-header";
import { ListingsGrid } from "@/components/listings/listings-grid";

export const metadata: Metadata = {
  title: "Explore gear | YLX",
  description: "Browse developer-owned tech gear on YLX.",
};

export default function ListingsPage() {
  return (
    <>
      <HomeHeader />

      <main id="main-content" className="px-6 py-16 sm:px-10">
        <section className="mx-auto max-w-5xl">
          <p className="text-sm font-semibold tracking-wide text-primary uppercase">
            Developer to developer
          </p>
          <h1 className="mt-3 text-4xl sm:text-5xl">Explore gear</h1>
          <p className="mt-4 max-w-lg text-lg text-muted-foreground">
            Find your next laptop, component, peripheral, or homelab upgrade.
          </p>

          <div className="mt-10">
            <ListingsGrid />
          </div>
        </section>
      </main>
    </>
  );
}
