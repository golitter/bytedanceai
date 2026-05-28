interface ToolCardProps {
  name?: string
  input?: string
  output?: string
}

export function ToolCard({ name, input, output }: ToolCardProps) {
  return (
    <div className="my-1 rounded-lg border border-border bg-secondary px-3 py-2 text-[13px]">
      {name && <div className="mb-1 text-xs text-muted-foreground">🔧 {name}</div>}
      {input && <code className="rounded bg-muted px-1.5 py-0.5 font-mono text-xs">{input}</code>}
      {output && <div className="mt-1 text-xs text-agent-opencode">{output}</div>}
    </div>
  )
}
