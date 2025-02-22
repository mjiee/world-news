import { create } from "zustand";
import { getRemoteService, saveRemoteService, SystemConfigKey } from "@/services";

interface RemoteServiceState {
  enable: boolean;
  host?: string;
  token?: string;
  getService: () => void;
  saveService: (enable: boolean, host: string | undefined) => void;
}

interface RemoteServiceValue {
  enable: boolean;
  host?: string;
  token?: string;
}

export const useRemoteServiceStore = create<RemoteServiceState>((set) => ({
  enable: false,
  host: "http://localhost:9010",
  getService: async () => {
    const resp = await getRemoteService<RemoteServiceValue>({ key: SystemConfigKey.RemoteService });

    if (!resp || !resp.value) return;

    set(() => {
      return { ...resp.value };
    });
  },
  saveService: async (enable: boolean, host: string | undefined) =>
    set((state) => {
      let data = {
        ...state,
        enable: enable,
        host: host,
      };

      saveRemoteService({ key: SystemConfigKey.RemoteService, value: data });

      return data;
    }),
}));
