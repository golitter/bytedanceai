import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import {
  addToContactGroup,
  createContactGroup,
  deleteContactGroup,
  fetchContactGroups,
  removeFromContactGroup,
  updateContactGroup,
} from '@/lib/api'

export function useContactGroups() {
  return useQuery({
    queryKey: ['contact-groups'],
    queryFn: fetchContactGroups,
  })
}

export function useCreateContactGroup() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (name: string) => createContactGroup(name),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['contact-groups'] })
    },
  })
}

export function useUpdateContactGroup() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ groupId, name }: { groupId: string; name: string }) =>
      updateContactGroup(groupId, name),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['contact-groups'] })
    },
  })
}

export function useDeleteContactGroup() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (groupId: string) => deleteContactGroup(groupId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['contact-groups'] })
    },
  })
}

export function useAddToContactGroup() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ groupId, taskId }: { groupId: string; taskId: string }) =>
      addToContactGroup(groupId, taskId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['contact-groups'] })
    },
  })
}

export function useRemoveFromContactGroup() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ groupId, taskId }: { groupId: string; taskId: string }) =>
      removeFromContactGroup(groupId, taskId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['contact-groups'] })
    },
  })
}
