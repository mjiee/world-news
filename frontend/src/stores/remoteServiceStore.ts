import { create } from "zustand";
import { getSystemConfig } from "@/services";
import { isWeb } from "@/utils/platform";

interface RemoteServiceState {
  enable: boolean;
  host?: string;
  token?: string;
  setEnable: (newEnable: boolean) => void;
  setHost: (newHost: string) => void;
  setToken: (newToken: string) => void;
  setService: (data: RemoteServiceValue) => void;
}

interface RemoteServiceValue {
  enable: boolean;
  host?: string;
  token?: string;
}

const remoteServiceKey = "remoteService";

export const useRemoteServiceStore = create<RemoteServiceState>((set) => ({
  enable: false,
  setEnable: (newEnable: boolean) => set((state) => ({ ...state, enable: newEnable })),
  setHost: (newHost: string) => set((state) => ({ ...state, host: newHost })),
  setToken: (newToken: string) => set((state) => ({ ...state, token: newToken })),
  setService: (data: RemoteServiceValue) => set(data),
}));

// initRemoteService is used to initialize the remote service
export async function initRemoteService() {
  if (isWeb()) return;

  const setService = useRemoteServiceStore((state) => state.setService);

  const resp = await getSystemConfig<RemoteServiceValue>({ key: remoteServiceKey });

  if (!resp || !resp.value) return;

  setService(resp.value);
}
