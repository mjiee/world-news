import toast from "react-hot-toast";
import axios from "axios";
import { GolbalLanguage, useServiceToken } from "@/stores";
import { httpx } from "wailsjs/go/models";
import { LogError } from "wailsjs/runtime";
import { isWeb } from "./platform";
import { setHost } from "./url";

export interface Response<R> {
  code?: number;
  message?: string;
  result?: R;
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
    return await serviceAxios.post(setHost(url), params);
  } catch (error: any) {
    return undefined;
  }
}

// get http request
export async function get<P, R>(url: string, params?: P): Promise<R | undefined> {
  try {
    return await serviceAxios.get(setHost(url), { params: params });
  } catch (error: any) {
    return undefined;
  }
}

// serviceAxios is used to handle the results returned by http request
const serviceAxios = axios.create({
  timeout: 120000,
  withCredentials: false,
});

// request interceptor
serviceAxios.interceptors.request.use(
  (config) => {
    // set language
    config.headers["Accept-Language"] = GolbalLanguage.getLanguage();

    // set service token
    const authorizationBasic = btoa("token:" + useServiceToken());

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

// notificationError is used to show notification error
function notificationError(message?: string | undefined) {
  toast.error(message ?? "internal error");
}
