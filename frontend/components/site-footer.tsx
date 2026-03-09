import Link from "next/link"
import { Github, Twitter, Mail } from "lucide-react"

const footerSections = [
  {
    title: "\u793e\u533a",
    links: [
      { label: "\u4e3b\u9898\u8ba8\u8bba", href: "/topics" },
      { label: "\u6280\u672f\u6587\u7ae0", href: "/articles" },
      { label: "\u5f00\u6e90\u9879\u76ee", href: "/projects" },
      { label: "\u8d44\u6e90\u5206\u4eab", href: "/resources" },
      { label: "\u9177\u5de5\u4f5c", href: "/jobs" },
    ],
  },
  {
    title: "\u5b66\u4e60",
    links: [
      { label: "Go \u56fe\u4e66", href: "/books" },
      { label: "Go \u6307\u5357", href: "/docs/guide" },
      { label: "\u6807\u51c6\u5e93\u6587\u6863", href: "/docs/stdlib" },
      { label: "\u6bcf\u65e5\u9762\u8bd5\u9898", href: "/interview" },
    ],
  },
  {
    title: "\u5173\u4e8e",
    links: [
      { label: "\u5173\u4e8e\u6211\u4eec", href: "/about" },
      { label: "\u53cd\u9988\u5efa\u8bae", href: "/feedback" },
      { label: "\u5fd7\u613f\u8005\u62db\u52df", href: "/volunteer" },
      { label: "\u5e7f\u544a\u5408\u4f5c", href: "/advertise" },
    ],
  },
]

export function SiteFooter() {
  return (
    <footer className="border-t border-border bg-card">
      <div className="mx-auto max-w-7xl px-4 py-10 lg:px-6">
        <div className="grid grid-cols-2 gap-8 md:grid-cols-4">
          {/* Brand */}
          <div className="col-span-2 md:col-span-1">
            <Link href="/" className="flex items-center gap-2.5">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary">
                <span className="text-sm font-bold text-primary-foreground">Go</span>
              </div>
              <span className="text-lg font-bold tracking-tight text-foreground">
                Go<span className="text-primary">{"\u4e2d\u6587\u7f51"}</span>
              </span>
            </Link>
            <p className="mt-3 text-sm leading-relaxed text-muted-foreground">
              {"\u4e2d\u56fd\u6700\u5927\u7684 Go \u8bed\u8a00\u793e\u533a\uff0c\u81f4\u529b\u4e8e\u6784\u5efa\u5b8c\u5584\u7684 Go \u8bed\u8a00\u4e2d\u6587\u751f\u6001\u3002"}
            </p>
            <div className="mt-4 flex items-center gap-3">
              <a
                href="https://github.com/studygolang"
                target="_blank"
                rel="noopener noreferrer"
                className="flex h-8 w-8 items-center justify-center rounded-md bg-secondary text-muted-foreground transition-colors hover:bg-primary/10 hover:text-primary"
                aria-label="GitHub"
              >
                <Github className="h-4 w-4" />
              </a>
              <a
                href="#"
                className="flex h-8 w-8 items-center justify-center rounded-md bg-secondary text-muted-foreground transition-colors hover:bg-primary/10 hover:text-primary"
                aria-label="Twitter"
              >
                <Twitter className="h-4 w-4" />
              </a>
              <a
                href="mailto:contact@studygolang.com"
                className="flex h-8 w-8 items-center justify-center rounded-md bg-secondary text-muted-foreground transition-colors hover:bg-primary/10 hover:text-primary"
                aria-label={"\u90ae\u4ef6"}
              >
                <Mail className="h-4 w-4" />
              </a>
            </div>
          </div>

          {/* Links */}
          {footerSections.map((section) => (
            <div key={section.title}>
              <h3 className="mb-3 text-sm font-semibold text-foreground">
                {section.title}
              </h3>
              <ul className="space-y-2">
                {section.links.map((link) => (
                  <li key={link.href}>
                    <Link
                      href={link.href}
                      className="text-sm text-muted-foreground transition-colors hover:text-primary"
                    >
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>

        <div className="mt-8 border-t border-border pt-6">
          <div className="flex flex-col items-center justify-between gap-3 sm:flex-row">
            <p className="text-xs text-muted-foreground">
              {"2013 - 2026 Go\u8bed\u8a00\u4e2d\u6587\u7f51 | Golang\u4e2d\u6587\u793e\u533a | "}
              <a href="https://beian.miit.gov.cn/" target="_blank" rel="noopener noreferrer" className="hover:text-primary">
                {"ICP\u5907\u6848\u53f7"}
              </a>
            </p>
            <p className="text-xs text-muted-foreground">
              {"Powered by studygolang.com"}
            </p>
          </div>
        </div>
      </div>
    </footer>
  )
}
