import { HeadContent, Scripts, createRootRoute } from '@tanstack/react-router'
import { TanStackRouterDevtoolsPanel } from '@tanstack/react-router-devtools'
import { TanStackDevtools } from '@tanstack/react-devtools'

import appCss from '../styles.css?url'

const site = {
  title: 'YLX - The Marketplace for Developers',
  description:
    'A developer-first marketplace for buying and selling tech gear.',
  image: '/brand/ylx-thumbnail.png',
} as const

export const Route = createRootRoute({
  head: () => ({
    meta: [
      {
        charSet: 'utf-8',
      },
      {
        name: 'viewport',
        content: 'width=device-width, initial-scale=1',
      },
      {
        title: site.title,
      },
      {
        name: 'description',
        content: site.description,
      },
      {
        name: 'theme-color',
        content: '#0068f2',
      },
      {
        property: 'og:title',
        content: site.title,
      },
      {
        property: 'og:description',
        content: site.description,
      },
      {
        property: 'og:image',
        content: site.image,
      },
      {
        property: 'og:type',
        content: 'website',
      },
      {
        name: 'twitter:card',
        content: 'summary_large_image',
      },
      {
        name: 'twitter:title',
        content: site.title,
      },
      {
        name: 'twitter:description',
        content: site.description,
      },
      {
        name: 'twitter:image',
        content: site.image,
      },
    ],
    links: [
      {
        rel: 'stylesheet',
        href: appCss,
      },
      {
        rel: 'icon',
        href: '/brand/favicon.ico',
      },
      {
        rel: 'icon',
        type: 'image/png',
        sizes: '32x32',
        href: '/brand/favicon-32.png',
      },
      {
        rel: 'icon',
        type: 'image/png',
        sizes: '16x16',
        href: '/brand/favicon-16.png',
      },
      {
        rel: 'apple-touch-icon',
        href: '/brand/favicon-192.png',
      },
      {
        rel: 'manifest',
        href: '/site.webmanifest',
      },
    ],
  }),
  shellComponent: RootDocument,
})

function RootDocument({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <head>
        <HeadContent />
      </head>
      <body>
        {children}
        <TanStackDevtools
          config={{
            position: 'bottom-right',
          }}
          plugins={[
            {
              name: 'Tanstack Router',
              render: <TanStackRouterDevtoolsPanel />,
            },
          ]}
        />
        <Scripts />
      </body>
    </html>
  )
}
