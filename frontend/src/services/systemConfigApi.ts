import { isWeb } from "@/utils/platform";
import { call } from "@/utils/http";
import { GetSystemConfig } from "wailsjs/go/adapter/App";
import { dto } from "wailsjs/go/models";

export interface SystemConfig {
  key: string;
  value?: any;
}

interface GetSystemConfigRequest {
  key: string;
}

// getSystemConfig to get system config
export async function getSystemConfig(request: GetSystemConfigRequest): Promise<SystemConfig | undefined> {
  if (isWeb()) {
    return undefined;
  }

  return await call<SystemConfig>(GetSystemConfig(request));
}
