export function useHoverStyle(hoverBg = 'var(--accent)', normalBg = 'transparent') {
  return {
    onMouseEnter: (e: React.MouseEvent<HTMLElement>) => {
      e.currentTarget.style.backgroundColor = hoverBg
    },
    onMouseLeave: (e: React.MouseEvent<HTMLElement>) => {
      e.currentTarget.style.backgroundColor = normalBg
    },
  }
}
