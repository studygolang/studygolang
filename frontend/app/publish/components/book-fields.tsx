import { Label } from "@/components/ui/label"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"

interface BookFieldsProps {
  bookAuthor: string
  setBookAuthor: (author: string) => void
  bookTranslator: string
  setBookTranslator: (translator: string) => void
  bookCover: string
  setBookCover: (cover: string) => void
  bookPubDate: string
  setBookPubDate: (date: string) => void
  bookLang: string
  setBookLang: (lang: string) => void
  bookIsFree: boolean
  setBookIsFree: (isFree: boolean) => void
  bookOnlineUrl: string
  setBookOnlineUrl: (url: string) => void
  bookDownloadUrl: string
  setBookDownloadUrl: (url: string) => void
  bookBuyUrl: string
  setBookBuyUrl: (url: string) => void
  bookPrice: string
  setBookPrice: (price: string) => void
}

export function BookFields({
  bookAuthor,
  setBookAuthor,
  bookTranslator,
  setBookTranslator,
  bookCover,
  setBookCover,
  bookPubDate,
  setBookPubDate,
  bookLang,
  setBookLang,
  bookIsFree,
  setBookIsFree,
  bookOnlineUrl,
  setBookOnlineUrl,
  bookDownloadUrl,
  setBookDownloadUrl,
  bookBuyUrl,
  setBookBuyUrl,
  bookPrice,
  setBookPrice,
}: BookFieldsProps) {
  return (
    <>
      {/* 作者 + 译者 */}
      <div className="flex gap-4">
        <div className="flex-1 space-y-1.5">
          <Label htmlFor="bookAuthor" className="text-sm font-medium">
            作者 <span className="text-destructive">*</span>
          </Label>
          <Input
            id="bookAuthor"
            placeholder="作者姓名"
            value={bookAuthor}
            onChange={(e) => setBookAuthor(e.target.value)}
            className="h-10"
          />
        </div>
        <div className="flex-1 space-y-1.5">
          <Label htmlFor="bookTranslator" className="text-sm font-medium">译者</Label>
          <Input
            id="bookTranslator"
            placeholder="译者姓名（如有）"
            value={bookTranslator}
            onChange={(e) => setBookTranslator(e.target.value)}
            className="h-10"
          />
        </div>
      </div>

      {/* 封面 + 出版日期 */}
      <div className="flex gap-4">
        <div className="flex-1 space-y-1.5">
          <Label htmlFor="bookCover" className="text-sm font-medium">封面图片</Label>
          <Input
            id="bookCover"
            placeholder="封面图片URL"
            value={bookCover}
            onChange={(e) => setBookCover(e.target.value)}
            className="h-10"
          />
        </div>
        <div className="flex-1 space-y-1.5">
          <Label htmlFor="bookPubDate" className="text-sm font-medium">出版日期</Label>
          <Input
            id="bookPubDate"
            placeholder="例如：2024-01"
            value={bookPubDate}
            onChange={(e) => setBookPubDate(e.target.value)}
            className="h-10"
          />
        </div>
      </div>

      {/* 语言 + 是否免费 */}
      <div className="flex gap-4">
        <div className="flex-1 space-y-1.5">
          <Label htmlFor="bookLang" className="text-sm font-medium">语言</Label>
          <Select value={bookLang} onValueChange={setBookLang}>
            <SelectTrigger className="h-10">
              <SelectValue placeholder="选择语言" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="中文">中文</SelectItem>
              <SelectItem value="英文">英文</SelectItem>
              <SelectItem value="中英双语">中英双语</SelectItem>
              <SelectItem value="其他">其他</SelectItem>
            </SelectContent>
          </Select>
        </div>
        <div className="flex-1 space-y-1.5">
          <Label className="text-sm font-medium">是否免费</Label>
          <label className="flex items-center gap-2 h-10 cursor-pointer">
            <input
              type="checkbox"
              checked={bookIsFree}
              onChange={(e) => setBookIsFree(e.target.checked)}
              className="h-4 w-4"
            />
            <span className="text-sm">免费资源</span>
          </label>
        </div>
      </div>

      {/* 在线阅读 + 下载地址 */}
      <div className="flex gap-4">
        <div className="flex-1 space-y-1.5">
          <Label htmlFor="bookOnlineUrl" className="text-sm font-medium">在线阅读</Label>
          <Input
            id="bookOnlineUrl"
            placeholder="在线阅读地址"
            value={bookOnlineUrl}
            onChange={(e) => setBookOnlineUrl(e.target.value)}
            className="h-10"
          />
        </div>
        <div className="flex-1 space-y-1.5">
          <Label htmlFor="bookDownloadUrl" className="text-sm font-medium">下载地址</Label>
          <Input
            id="bookDownloadUrl"
            placeholder="电子版下载地址"
            value={bookDownloadUrl}
            onChange={(e) => setBookDownloadUrl(e.target.value)}
            className="h-10"
          />
        </div>
      </div>

      {/* 购买地址 + 价格 */}
      <div className="flex gap-4">
        <div className="flex-1 space-y-1.5">
          <Label htmlFor="bookBuyUrl" className="text-sm font-medium">购买地址</Label>
          <Input
            id="bookBuyUrl"
            placeholder="购买链接"
            value={bookBuyUrl}
            onChange={(e) => setBookBuyUrl(e.target.value)}
            className="h-10"
          />
        </div>
        <div className="flex-1 space-y-1.5">
          <Label htmlFor="bookPrice" className="text-sm font-medium">价格</Label>
          <Input
            id="bookPrice"
            placeholder="例如：99.00"
            value={bookPrice}
            onChange={(e) => setBookPrice(e.target.value)}
            className="h-10"
          />
        </div>
      </div>
    </>
  )
}
