import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import { createTask, fetchTask, patchSession } from '@/lib/api'

export function useTask(taskId: string) {
  return useQuery({
    queryKey: ['task', taskId],
    queryFn: () => fetchTask(taskId),
    enabled: !!taskId,
  })
}

export function useSessions(taskId: string) {
  const { data, ...rest } = useTask(taskId)
  return {
    ...rest,
    data: data?.sessions ?? [],
    task: data?.task,
  }
}

export function useCreateTask() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (title: string) => createTask(title),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] })
    },
  })
}

export function useDeleteTask() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (taskId: string) =>
      fetch(`/api/tasks/${taskId}`, { method: 'DELETE' }).then((r) => {
        if (!r.ok) throw new Error('Failed to delete task')
      }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['tasks'] })
    },
  })
}

export function useDeactivateSession(taskId: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (sessionId: string) => patchSession(sessionId, 'inactive'),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['task', taskId] })
    },
  })
}
