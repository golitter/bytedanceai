import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import type { AgentType } from '@/generated/request'
import { createConversation, fetchConversations } from '@/lib/api'

export function useConversations() {
  return useQuery({
    queryKey: ['conversations'],
    queryFn: fetchConversations,
  })
}

export function useCreateConversation() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (params: {
      agents: { type: AgentType; name: string }[]
      repoPath?: string
      title?: string
    }) => createConversation(params.agents, params.repoPath, params.title),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['conversations'] })
    },
  })
}
