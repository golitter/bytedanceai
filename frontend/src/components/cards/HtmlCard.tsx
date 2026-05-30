interface HtmlCardProps {
  content: string
  expanded?: boolean
}

export function HtmlCard({ content, expanded }: HtmlCardProps) {
  return (
    <div className="my-2 overflow-hidden rounded-lg border border-border">
      <iframe
        sandbox=""
        srcDoc={content}
        className={expanded ? 'h-[min(72vh,760px)] w-full border-0' : 'h-64 w-full border-0'}
        title="HTML Preview"
      />
    </div>
  )
}
