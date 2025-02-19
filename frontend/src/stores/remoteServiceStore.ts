import { create } from "zustand";
import { getRemoteService, saveRemoteService } from "@/services";
import { isWeb } from "@/utils/platform";

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

const remoteServiceKey = "remoteService";

export const useRemoteServiceStore = create<RemoteServiceState>((set) => ({
  enable: false,
  host: "http://localhost:9010",
  getService: async () => {
    if (isWeb()) return;

    const resp = await getRemoteService<RemoteServiceValue>({ key: remoteServiceKey });

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

      if (!isWeb()) saveRemoteService({ key: remoteServiceKey, value: data });

      return data;
    }),
}));
