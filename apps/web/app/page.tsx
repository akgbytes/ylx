import Link from "next/link";
import { IconArrowRight } from "@ylx/ui/icons";

import { buttonVariants } from "@ylx/ui/components/button";

import { HomeHeader } from "@/components/home/home-header";
import { ListingsGrid } from "@/components/listings/listings-grid";

export default function Home() {
  return (
    <>
      <HomeHeader />

      <main id="main-content" className="px-6 py-16 sm:px-10">
        <section className="mx-auto max-w-5xl">
          <div className="max-w-xl space-y-6">
            <p className="text-sm font-semibold tracking-wide text-primary uppercase">
              Developer to developer
            </p>
            <h1 className="max-w-xl text-4xl text-balance sm:text-5xl">
              Great tech deserves its next developer.
            </h1>
            <p className="max-w-lg text-lg text-muted-foreground">
              Buy and sell laptops, components, peripherals, and homelab gear
              directly with other developers.
            </p>
            <Link href="/listings" className={buttonVariants({ size: "lg" })}>
              Explore gear <IconArrowRight aria-hidden="true" />
            </Link>
          </div>
        </section>

        <section
          className="mx-auto mt-16 max-w-5xl"
          aria-labelledby="latest-gear"
        >
          <div className="mb-6 flex items-center justify-between gap-4">
            <h2
              id="latest-gear"
              className="font-heading text-2xl font-semibold"
            >
              Latest gear
            </h2>
            <Link
              href="/listings"
              className={buttonVariants({ variant: "link", size: "sm" })}
            >
              View all
            </Link>
          </div>

          <ListingsGrid maxItems={3} />
        </section>
      </main>
    </>
  );
}
