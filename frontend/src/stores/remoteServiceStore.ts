import { create } from "zustand";
import { getRemoteService, saveRemoteService, SystemConfigKey } from "@/services";

interface RemoteServiceState {
  enable: boolean;
  host?: string;
  token?: string;
  getService: () => void;
  setToken: (token: string) => void;
  saveService: (enable: boolean, host: string | undefined, token: string | undefined) => void;
}

interface RemoteServiceValue {
  enable: boolean;
  host?: string;
  token?: string;
}

export const useRemoteServiceStore = create<RemoteServiceState>((set) => ({
  enable: false,
  host: "http://localhost:9010",
  token: "",
  getService: async () => {
    const resp = await getRemoteService<RemoteServiceValue>({ key: SystemConfigKey.RemoteService });

    if (!resp || !resp.value) return;

    set(() => {
      return { ...resp.value };
    });
  },
  setToken: (token: string) => set((state) => ({ ...state, token: token })),
  saveService: async (enable: boolean, host: string | undefined, token: string | undefined) =>
    set((state) => {
      let data = {
        ...state,
        enable: enable,
        host: host,
        token: token,
      };

      saveRemoteService({ key: SystemConfigKey.RemoteService, value: data });

      return data;
    }),
}));
