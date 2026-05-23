import { create } from 'zustand'

interface ChatNavState {
  currentSessionId: string | null
  currentTaskId: string | null
  setCurrentSession: (sessionId: string) => void
  setCurrentTask: (taskId: string) => void
  clearNavigation: () => void
}

export const useChatNav = create<ChatNavState>((set) => ({
  currentSessionId: null,
  currentTaskId: null,
  setCurrentSession: (sessionId) => set({ currentSessionId: sessionId }),
  setCurrentTask: (taskId) => set({ currentTaskId: taskId }),
  clearNavigation: () => set({ currentSessionId: null, currentTaskId: null }),
}))
