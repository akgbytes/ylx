import { createFileRoute } from '@tanstack/react-router'
import { Button } from '@ylx/ui/components/button'
import { IconArrowRight, IconMapPin } from '@ylx/ui/icons'

export const Route = createFileRoute('/')({ component: Home })

function Home() {
  return (
    <>
      <header className="border-b bg-card/80 px-6 py-4 backdrop-blur sm:px-10">
        <div className="mx-auto flex max-w-5xl items-center gap-3">
          <img
            src="/brand/favicon-64.png"
            alt=""
            className="size-10"
            width="40"
            height="40"
          />
          <div>
            <p className="font-heading text-xl font-bold leading-none">YLX</p>
            <p className="mt-1 text-xs text-muted-foreground">
              The Marketplace for Developers
            </p>
          </div>
        </div>
      </header>

      <main className="px-6 py-16 sm:px-10">
        <section className="mx-auto grid max-w-5xl gap-12 lg:grid-cols-[1fr_0.8fr] lg:items-center">
          <div className="space-y-6">
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
            <Button size="lg">
              Explore gear <IconArrowRight aria-hidden="true" />
            </Button>
          </div>

          <article className="rounded-xl border bg-card p-5 text-card-foreground shadow-sm">
            <div className="mb-16 grid aspect-4/3 place-items-center rounded-lg bg-accent">
              <img
                src="/brand/favicon-192.png"
                alt=""
                className="size-28 drop-shadow-sm"
                width="112"
                height="112"
              />
            </div>
            <div className="space-y-3">
              <div className="flex items-start justify-between gap-4">
                <div>
                  <p className="font-mono text-xs text-muted-foreground">
                    DEV-GEAR · FW13-R7
                  </p>
                  <h2 className="mt-1 text-xl">
                    Framework Laptop 13 · Ryzen 7
                  </h2>
                </div>
                <p className="shrink-0 text-xl font-bold tabular-nums">
                  ₹82,000
                </p>
              </div>
              <p className="flex items-center gap-1.5 text-sm text-muted-foreground">
                <IconMapPin className="size-4" aria-hidden="true" /> Bengaluru,
                Karnataka
              </p>
            </div>
          </article>
        </section>
      </main>
    </>
  )
}
