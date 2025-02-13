import axios, { AxiosResponse } from "axios";
import { useTranslation } from "react-i18next";
import { useNotificationStore, useRemoteServiceStore } from "@/stores";
import { httpx } from "wailsjs/go/models";
import { isWeb } from "./platform";

export interface Response<T> {
  code?: number;
  message?: string;
  result?: T;
}

// useRemoteService is used to check if the remote service is enabled
export function useRemoteService(): boolean {
  const { enable, host } = useRemoteServiceStore.getState();

  return isWeb() || (enable && host !== null && !!host);
}

// post http request
export async function post<T, K>(url: string, params?: T): Promise<K | undefined> {
  return await resultHandle(axios.post<Response<K>>(urlHandle(url), params));
}

// get http request
export async function get<T, K>(url: string, params?: T): Promise<K | undefined> {
  return await resultHandle(axios.get<Response<K>>(urlHandle(url), { params: params }));
}

// urlHandle is used to handle the url
function urlHandle(url: string): string {
  if (isWeb()) return url;

  const host = useRemoteServiceStore((state) => state.host);

  return new URL(url, host as string).href;
}

// call is used to handle the results returned by wails.
export async function call<T>(resp: Promise<httpx.Response>): Promise<T | undefined> {
  return await resultHandle(resp);
}

// resultHandle is used to handle the results returned by axios or wails.
async function resultHandle<T>(resp: Promise<AxiosResponse<Response<T>, any>> | Promise<httpx.Response>): Promise<T | undefined> {
  let data = await resp
    .then((resp) => {
      if (resp instanceof httpx.Response) {
        if (resp?.code !== 0) {
          notificationError(resp.message);

          return undefined;
        }

        return resp.result;
      } else {
        if (resp?.status !== 200 && resp?.data?.code !== 0) {
          notificationError(resp.data.message ?? undefined);

          return undefined;
        }

        return resp.data.result;
      }
    })
    .catch((err) => {
      notificationError();

      return undefined;
    });

  return data ?? undefined;
}

// notificationError is used to show notification error
function notificationError(message?: string | undefined) {
  const { t } = useTranslation();
  const showNotification = useNotificationStore((state) => state.showNotification);

  showNotification(message ?? t("message.internal_error"));
}
