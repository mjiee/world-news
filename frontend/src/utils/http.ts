import axios, { AxiosResponse } from "axios";
import toast from "react-hot-toast";
import { useRemoteServiceStore } from "@/stores";
import { httpx } from "wailsjs/go/models";
import { isWeb } from "./platform";

export interface Response<R> {
  code?: number;
  message?: string;
  result?: R;
}

// useRemoteService is used to check if the remote service is enabled
export function useRemoteService(): boolean {
  const { enable, host } = useRemoteServiceStore.getState();

  return isWeb() || (enable && host !== null && !!host);
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
  } catch (error) {
    return undefined;
  }
}

// post http request
export async function post<P, R>(url: string, params?: P): Promise<R | undefined> {
  return await httpResultHandle(axios.post<Response<R>>(urlHandle(url), params));
}

// get http request
export async function get<P, R>(url: string, params?: P): Promise<R | undefined> {
  return await httpResultHandle(axios.get<Response<R>>(urlHandle(url), { params: params }));
}

// urlHandle is used to handle the url
function urlHandle(url: string): string {
  if (isWeb()) return url;

  const host = useRemoteServiceStore((state) => state.host);

  return new URL(url, host as string).href;
}

// httpResultHandle is used to handle the results returned by http request
async function httpResultHandle<R>(resp: Promise<AxiosResponse<Response<R>, any>>): Promise<R | undefined> {
  let data = await resp
    .then((resp) => {
      if (resp?.status !== 200 && resp?.data?.code !== 0) {
        notificationError(resp?.data?.message ?? undefined);

        return;
      }

      return resp.data.result;
    })
    .catch((err) => {
      notificationError();
    });

  return data ?? undefined;
}

// notificationError is used to show notification error
function notificationError(message?: string | undefined) {
  toast.error(message ?? "internal error");
}
