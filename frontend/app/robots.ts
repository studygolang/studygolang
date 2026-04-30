import type { MetadataRoute } from 'next'

export default function robots(): MetadataRoute.Robots {
  return {
    rules: [
      {
        userAgent: '*',
        allow: '/',
        disallow: [
          '/account/',
          '/admin/',
          '/api/',
          '/message/',
          '/messages/',
          '/publish/',
          '/balance/',
          '/favorites/',
          '/oauth/',
        ],
      },
    ],
    sitemap: 'https://studygolang.com/sitemap/sitemapindex.xml',
  }
}
