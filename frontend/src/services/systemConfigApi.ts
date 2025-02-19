import { call, post, useRemoteService } from "@/utils/http";
import { GetSystemConfig, SaveSystemConfig } from "wailsjs/go/adapter/App";
import { dto } from "wailsjs/go/models";

interface SystemConfig<T> {
  key: string;
  value?: T;
}

interface GetSystemConfigRequest {
  key: string;
}

// getSystemConfig to get system config
export async function getSystemConfig<T>(request: GetSystemConfigRequest): Promise<SystemConfig<T> | undefined> {
  if (useRemoteService()) return await post<GetSystemConfigRequest, SystemConfig<T>>("/api/system/config", request);

  return await call<SystemConfig<T>>(GetSystemConfig(request));
}

// saveSystemConfig to save system config
export async function saveSystemConfig<T>(request: SystemConfig<T>) {
  if (useRemoteService()) return await post<SystemConfig<T>, any>("/api/system/config/save", request);

  return await call(SaveSystemConfig(new dto.SystemConfig(request)));
}

// saveRemoteService is used to save the remote service
export const saveRemoteService = async <T>(data: SystemConfig<T>) => {
  return await call(SaveSystemConfig(new dto.SystemConfig(data)));
};

// getRemoteService is used to get the remote service
export const getRemoteService = async <T>(request: GetSystemConfigRequest): Promise<SystemConfig<T> | undefined> => {
  return await call<SystemConfig<T>>(GetSystemConfig(request));
};
