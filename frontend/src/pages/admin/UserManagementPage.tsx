import { useQuery } from '@tanstack/react-query'
import { Camera } from 'lucide-react'
import { useRef, useState } from 'react'

import { getAdminAvatar, updateAdminAvatar, uploadAvatar } from '@/lib/api'
import { CURRENT_USER_NAME } from '@/lib/constants'
import { useAdminStore } from '@/stores/admin'

export function UserManagementPage() {
  const { data, isLoading } = useQuery({
    queryKey: ['admin-avatar'],
    queryFn: getAdminAvatar,
    staleTime: 30_000,
  })
  // Track a locally overridden URL (from upload); fall back to query data
  const [localAvatarUrl, setLocalAvatarUrl] = useState<string | null>(null)
  const setAdminAvatarUrl = useAdminStore((s) => s.setAdminAvatarUrl)
  const [uploading, setUploading] = useState(false)
  const fileRef = useRef<HTMLInputElement>(null)

  const avatarUrl = localAvatarUrl ?? data?.url ?? ''

  const handleAvatarChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (!file) return
    setUploading(true)
    try {
      const url = await uploadAvatar(file)
      await updateAdminAvatar(url)
      setLocalAvatarUrl(url)
      setAdminAvatarUrl(url)
    } catch {
      // ignore
    } finally {
      setUploading(false)
    }
  }

  if (isLoading && !localAvatarUrl) {
    return (
      <div className="flex h-full items-center justify-center">
        <span className="text-sm text-tertiary">加载中...</span>
      </div>
    )
  }

  return (
    <div className="p-6">
      <h2 className="mb-6 text-lg font-semibold text-foreground">用户管理</h2>

      <div className="rounded-lg border border-border bg-card p-6">
        <h3 className="mb-4 text-sm font-medium text-text-secondary">管理员头像</h3>

        <div className="flex items-center gap-4">
          <div className="group relative">
            <div className="h-20 w-20 overflow-hidden rounded-lg">
              <img src={avatarUrl} alt="Admin Avatar" className="h-full w-full object-cover" />
            </div>
            <button
              className="absolute inset-0 flex items-center justify-center rounded-lg bg-black/40 opacity-0 transition-opacity group-hover:opacity-100"
              onClick={() => fileRef.current?.click()}
              disabled={uploading}
            >
              <Camera className="h-5 w-5 text-primary-foreground" strokeWidth={1.25} />
            </button>
            <input
              ref={fileRef}
              type="file"
              accept="image/jpeg,image/png,image/gif,image/webp"
              className="hidden"
              onChange={handleAvatarChange}
            />
          </div>

          <div>
            <p className="text-sm font-medium text-foreground">{CURRENT_USER_NAME}</p>
            <p className="mt-0.5 text-xs text-tertiary">管理员</p>
            <button
              className="mt-2 text-xs text-brand transition-[transform,opacity]"
              onClick={() => fileRef.current?.click()}
              disabled={uploading}
            >
              {uploading ? '上传中...' : '更换头像'}
            </button>
          </div>
        </div>
      </div>
    </div>
  )
}
