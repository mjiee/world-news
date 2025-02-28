import { create } from "zustand";
import { getSystemConfig, saveSystemConfig, SystemConfigKey } from "@/services";
import { isWeb } from "@/utils/platform";

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
    if (isWeb()) return;

    const resp = await getSystemConfig<RemoteServiceValue>({ key: SystemConfigKey.RemoteService }, true);

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

      if (!isWeb()) saveSystemConfig({ key: SystemConfigKey.RemoteService, value: data }, true);

      return data;
    }),
}));

// useRemoteService is used to check if the remote service is enabled
export function useRemoteService(forceLocal = false): boolean {
  if (isWeb()) return true;

  if (forceLocal) return false;

  const { enable, host } = useRemoteServiceStore.getState();

  return enable && host !== null && !!host;
}

// get service token
export function useServiceToken(): string {
  const state = useRemoteServiceStore.getState();

  return state.token ?? "";
}

// get service host
export function useServiceHost(): string | undefined {
  const state = useRemoteServiceStore.getState();
  return state.host;
}
