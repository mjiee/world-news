import { call, post, useRemoteService } from "@/utils/http";
import { isWeb } from "@/utils/platform";
import { GetSystemConfig, SaveSystemConfig } from "wailsjs/go/adapter/App";
import { dto } from "wailsjs/go/models";

export enum SystemConfigKey {
  NewsTopics = "newsTopics",
  NewsWebsiteCollections = "newsWebsiteCollections",
  NewsWebsites = "newsWebsites",
  Language = "language",
  RemoteService = "remoteService",
}

interface SystemConfig<T> {
  key: string;
  value?: T;
}

interface GetSystemConfigRequest {
  key: string;
}

export interface NewsWebsiteValue {
  url: string;
  selector?: NewsSelector;
}

export interface NewsSelector {
  website?: string; // news website selector

  topic?: string; // news topic selector
  link?: string; // news link selector

  title?: string; // news title selector
  time?: string; // publish time selector
  content?: string; // news content selector
  image?: string; // news image selector
  author?: string; // news author selector

  child?: NewsSelector;
}

// getSystemConfig to get system config
export async function getSystemConfig<T>(
  request: GetSystemConfigRequest,
  forceLocal = false,
): Promise<SystemConfig<T> | undefined> {
  if (useRemoteService(forceLocal))
    return await post<GetSystemConfigRequest, SystemConfig<T>>("/api/system/config", request);

  return await call<SystemConfig<T>>(GetSystemConfig(request));
}

// saveSystemConfig to save system config
export async function saveSystemConfig<T>(request: SystemConfig<T>, forceLocal = false) {
  if (useRemoteService(forceLocal)) return await post<SystemConfig<T>, any>("/api/system/config/save", request);

  return await call(SaveSystemConfig(new dto.SystemConfig(request)));
}
