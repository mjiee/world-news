import { call, useRemoteService } from "@/utils/http";
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
  if (useRemoteService()) return undefined;

  return await call<SystemConfig<T>>(GetSystemConfig(request));
}

// saveSystemConfig to save system config
export async function saveSystemConfig<T>(request: SystemConfig<T>) {
  if (useRemoteService()) return;

  return await call(SaveSystemConfig(new dto.SystemConfig(request)));
}
