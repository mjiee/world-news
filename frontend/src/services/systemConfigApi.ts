import { useRemoteService } from "@/stores";
import { call, post } from "@/utils/http";
import { GetSystemConfig, SaveSystemConfig, SaveWebsiteWeight } from "wailsjs/go/adapter/App";
import { dto } from "wailsjs/go/models";

export enum SystemConfigKey {
  NewsTopics = "newsTopics",
  NewsWebsiteCollections = "newsWebsiteCollections",
  NewsWebsites = "newsWebsites",
  Language = "language",
  RemoteService = "remoteService",
  Translater = "translater",
  TextAi = "textAI",
  TextToSpeechAi = "textToSpeechAI",
  NewsCritiquePromptKey = "newsCritiquePrompt",
  PodcastScriptPromptKey = "podcastScriptPrompt",
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
  weight: number;
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

export interface TextAIConfig {
  platform: string;
  apiKey: string;
  apiUrl: string;
  model: string;
}

export interface TextToSpeechAIConfig {
  platform: string;
  appId?: string;
  apiKey: string;
  model: string;
  voices?: AudioVoice[];
  autoTask?: boolean;
}

export interface AudioVoice {
  id: string;
  name: string;
  description?: string;
  model?: string;
}

export interface NewsCritiquePrompt {
  systemPrompt: string;
}

export interface PodcastScriptPrompt {
  systemPrompt: string;
  approvalPrompt?: string;
  rewritePrompt?: string;
  mergePrompt?: string;
  classifyPrompt?: string;
  stylizePrompts?: StylePrompt[];
}

export interface StylePrompt {
  style: string;
  prompt: string;
}

export interface TranslaterConfig {
  platform: string;
  appId: string;
  appSecret?: string;
}

export interface SaveWebsiteWeightRequest {
  website: string;
  step: number;
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

// saveWebsiteWeight to save website weight
export async function saveWebsiteWeight(data: SaveWebsiteWeightRequest) {
  if (useRemoteService()) return await post<SaveWebsiteWeightRequest, any>("/api/system/website/weight", data);

  return await call(SaveWebsiteWeight(data));
}
