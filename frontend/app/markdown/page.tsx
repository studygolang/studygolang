import type { Metadata } from 'next'
import { PageLayout } from '@/components/page-layout'
import { PageHeader } from '@/components/page-header'

export const metadata: Metadata = {
  title: 'Markdown 教程 - Go语言中文网',
  description: 'Markdown 语法教程，学习如何在文章和话题中使用 Markdown 排版',
}

export default function MarkdownPage() {
  return (
    <PageLayout sidebar={false}>
      <PageHeader
        title="Markdown 教程"
        description="学习如何在文章和话题中使用 Markdown 排版"
        breadcrumbs={[{ label: 'Markdown 教程' }]}
      />

      <div className="space-y-8">
        <section>
          <h2 className="text-2xl font-bold mb-4">语法指导</h2>

          <h3 className="text-xl font-semibold mb-3">普通内容</h3>
          <p className="text-muted-foreground mb-4">
            这段内容展示了在内容里面一些小的格式，比如：
          </p>
          <ul className="list-disc list-inside space-y-1 text-muted-foreground mb-6">
            <li><strong>加粗</strong> - <code className="bg-muted px-1 rounded">**加粗**</code></li>
            <li><em>倾斜</em> - <code className="bg-muted px-1 rounded">*倾斜*</code></li>
            <li><del>删除线</del> - <code className="bg-muted px-1 rounded">~~删除线~~</code></li>
            <li><code className="bg-muted px-1 rounded">Code 标记</code> - <code className="bg-muted px-1 rounded">`Code 标记`</code></li>
            <li>
              <a href="http://github.com" className="text-primary hover:underline" target="_blank" rel="noopener noreferrer">超级链接</a>
              {' '}- <code className="bg-muted px-1 rounded">[超级链接](http://github.com)</code>
            </li>
          </ul>

          <h3 className="text-xl font-semibold mb-3">提及用户</h3>
          <p className="text-muted-foreground mb-6">
            通过 <code className="bg-muted px-1 rounded">@</code> 可以在发帖和回帖里面提及用户，
            信息提交以后，被提及的用户将会收到系统通知。
          </p>

          <h3 className="text-xl font-semibold mb-3">表情符号 Emoji</h3>
          <p className="text-muted-foreground mb-6">
            支持表情符号，你可以用系统默认的 Emoji 符号，也可以用图片的表情，输入{' '}
            <code className="bg-muted px-1 rounded">:</code> 将会出现智能提示。
          </p>
        </section>

        <section>
          <h2 className="text-2xl font-bold mb-4">标题</h2>
          <p className="text-muted-foreground mb-4">
            你可以选择使用 H2 至 H6，使用 <code className="bg-muted px-1 rounded">##(N)</code> 打头，
            H1 不能使用，会自动转换成 H2。
          </p>
          <div className="bg-muted rounded-lg p-4 mb-6">
            <pre className="text-sm"><code>{`## Heading 2
### Heading 3
#### Heading 4
##### Heading 5
###### Heading 6`}</code></pre>
          </div>
        </section>

        <section>
          <h2 className="text-2xl font-bold mb-4">图片</h2>
          <div className="bg-muted rounded-lg p-4 mb-6">
            <pre className="text-sm"><code>{`![alt 文本](http://image-path.png)
![alt 文本](http://image-path.png "图片 Title 值")
![设置图片宽度高度](http://image-path.png =300x200)`}</code></pre>
          </div>
        </section>

        <section>
          <h2 className="text-2xl font-bold mb-4">代码块</h2>
          <p className="text-muted-foreground mb-4">
            如果在 <code className="bg-muted px-1 rounded">```</code> 后面跟随语言名称，可以有语法高亮的效果。
          </p>
          <div className="bg-muted rounded-lg p-4 mb-4">
            <pre className="text-sm"><code>{`// Go 代码高亮
\`\`\`go
package main

import "fmt"

func main() {
    fmt.Println("Hello World!")
}
\`\`\``}</code></pre>
          </div>
          <div className="bg-muted rounded-lg p-4 mb-6">
            <pre className="text-sm"><code>{`// JSON 代码高亮
\`\`\`json
{"name":"Go语言中文网","url":"https://studygolang.com"}
\`\`\``}</code></pre>
          </div>
        </section>

        <section>
          <h2 className="text-2xl font-bold mb-4">有序、无序列表</h2>
          <div className="grid md:grid-cols-2 gap-6 mb-6">
            <div>
              <h3 className="text-lg font-semibold mb-2">无序列表</h3>
              <div className="bg-muted rounded-lg p-4">
                <pre className="text-sm"><code>{`- Go
  - Gofmt
  - Revel
  - Gin
- PHP
  - Laravel
  - ThinkPHP`}</code></pre>
              </div>
            </div>
            <div>
              <h3 className="text-lg font-semibold mb-2">有序列表</h3>
              <div className="bg-muted rounded-lg p-4">
                <pre className="text-sm"><code>{`1. Go
    1. Gofmt
    2. Revel
2. PHP
    1. Laravel
3. Java`}</code></pre>
              </div>
            </div>
          </div>
        </section>

        <section>
          <h2 className="text-2xl font-bold mb-4">表格</h2>
          <p className="text-muted-foreground mb-4">如果需要展示数据，可以选择使用表格。</p>
          <div className="bg-muted rounded-lg p-4 mb-6">
            <pre className="text-sm"><code>{`| header 1 | header 2 |
| -------- | -------- |
| cell 1   | cell 2   |
| cell 3   | cell 4   |`}</code></pre>
          </div>
        </section>

        <section>
          <h2 className="text-2xl font-bold mb-4">引用</h2>
          <div className="bg-muted rounded-lg p-4 mb-6">
            <pre className="text-sm"><code>{`> 引用文本：Markdown is a text formatting syntax inspired`}</code></pre>
          </div>
        </section>

        <section>
          <h2 className="text-2xl font-bold mb-4">段落</h2>
          <p className="text-muted-foreground">
            留空白的换行，将会被自动转换成一个段落，会有一定的段落间距，便于阅读。
          </p>
        </section>
      </div>
    </PageLayout>
  )
}
