import { create } from "zustand";

interface remoteServiceState {
  enable: boolean;
  host: string | null;
  enableService: () => void;
  disableService: () => void;
  setServiceHost: (newHost: string | null) => void;
}

export const useRemoteServiceStore = create<remoteServiceState>((set) => ({
  enable: false,
  host: null,
  enableService: () => set({ enable: true }),
  disableService: () => set({ enable: false }),
  setServiceHost: (newHost: string | null) => set((state) => ({ ...state, host: newHost })),
}));
