import toast from "react-hot-toast";
import axios from "axios";
import { useRemoteServiceStore, useLanguageStore } from "@/stores";
import { httpx } from "wailsjs/go/models";
import { LogError } from "wailsjs/runtime";
import { isWeb } from "./platform";

export interface Response<R> {
  code?: number;
  message?: string;
  result?: R;
}

// useRemoteService is used to check if the remote service is enabled
export function useRemoteService(forceLocal = false): boolean {
  if (isWeb()) return true;

  if (forceLocal) return false;

  const { enable, host } = useRemoteServiceStore.getState();

  return enable && host !== null && !!host;
}

// call is used to handle the results returned by wails.
export async function call<R>(resp: Promise<httpx.Response>): Promise<R | undefined> {
  try {
    let result = await resp
      .then((resp) => {
        if (resp.code !== 0) {
          notificationError(resp?.message ?? undefined);

          return;
        }

        return resp.result;
      })
      .catch((err) => {
        notificationError();
      });

    return result ?? undefined;
  } catch (error: any) {
    LogError(error instanceof Error ? error.toString() : String(error));

    return undefined;
  }
}

// post http request
export async function post<P, R>(url: string, params?: P): Promise<R | undefined> {
  try {
    return await serviceAxios.post(urlHandle(url), params);
  } catch (error: any) {
    return undefined;
  }
}

// get http request
export async function get<P, R>(url: string, params?: P): Promise<R | undefined> {
  try {
    return await serviceAxios.get(urlHandle(url), { params: params });
  } catch (error: any) {
    return undefined;
  }
}

// serviceAxios is used to handle the results returned by http request
const serviceAxios = axios.create({
  timeout: 10000,
  withCredentials: false,
});

// get language
const getLanguage = () => {
  const state = useLanguageStore.getState();
  return state.language;
};

// get service token
const getServiceToken = () => {
  const state = useRemoteServiceStore.getState();
  return state.token;
};

// request interceptor
serviceAxios.interceptors.request.use(
  (config) => {
    // set language
    config.headers["Accept-Language"] = getLanguage();

    // set service token
    const authorizationBasic = btoa("token:" + getServiceToken());

    config.headers["Authorization"] = "Basic " + authorizationBasic;

    return config;
  },
  (error) => {
    notificationError();

    if (isWeb()) console.error(error);
    else LogError(error instanceof Error ? error.toString() : String(error));
  },
);

// response interceptor
serviceAxios.interceptors.response.use(
  (response) => {
    let data = response.data;

    if (data?.code !== 0) {
      if (data?.code === 401 && isWeb()) window.location.href = "/login";

      notificationError(data?.message ?? undefined);

      return;
    }

    return data?.result ?? undefined;
  },
  (error) => {
    notificationError();

    if (isWeb()) console.error(error);
    else LogError(error instanceof Error ? error.toString() : String(error));
  },
);

// get remote service host
const getRemoteServiceHost = () => {
  const state = useRemoteServiceStore.getState();
  return state.host;
};

// urlHandle is used to handle the url
function urlHandle(url: string): string {
  let host = getRemoteServiceHost();

  if (!host) host = import.meta.env.VITE_SERVICE_HOST;

  if (!host) return url;

  return new URL(url, host).href;
}

// notificationError is used to show notification error
function notificationError(message?: string | undefined) {
  toast.error(message ?? "internal error");
}
