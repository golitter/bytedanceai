import { create } from 'zustand'

export const useStore = create<Record<string, never>>(() => ({}))
