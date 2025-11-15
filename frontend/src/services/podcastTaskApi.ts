import { call } from "@/utils/http";
import {
  CreateAudio,
  CreateScript,
  CreateTask,
  DeleteTask,
  EditScript,
  GetTask,
  MergeArticle,
  QueryTasks,
  RestyleArticle,
} from "wailsjs/go/adapter/App";
import { dto, httpx } from "wailsjs/go/models";
import { NewsDetail } from "./newsApi";
import { AudioVoice } from "./systemConfigApi";

export interface MergeArticleRequest {
  title: string;
  stageIds?: number[];
  voiceIds?: string[];
}

export interface CreateTaskResult {
  batchNo: string;
}

export interface QueryTaskRequest {
  startDate?: string;
  endDate?: string;
  pagination?: httpx.Pagination;
}

export interface QueryTaskResult {
  data: PodcastTask[];
  total: number;
}

export interface PodcastTask {
  batchNo: string;
  title: string;
  news?: NewsDetail;
  language: string;
  result: PodcastTaskResult;
  stages?: TaskStage[];
  createdAt: string;
}

export enum PodcastTaskResult {
  Completed = "completed",
  Failed = "failed",
}

export interface TaskStage {
  id: number;
  batchNo: string;
  stage: TaskStageName;
  status: TaskStageStatus;
  prompt: string;
  output: string;
  reason: string;
  audio?: PodcastAudio;
  taskAi?: TaskAi;
  createdAt: string;
  updatedAt: string;
}

export enum TaskStageStatus {
  Processing = "processing",
  Completed = "completed",
  Failed = "failed",
}

export enum TaskStageName {
  Approval = "approval", // audit the article
  Rewrite = "rewrite", // rewrite it as an colloquial podcast article
  Merge = "merge", // merge the article
  Classify = "classify", // classify the topic
  Stylize = "stylize", // stylize the article
  Scripted = "scripted", // generate a podcast script
  TTS = "tts", // text to speech
}

export interface PodcastAudio {
  voices: AudioVoice[];
  type: string;
  url: string;
  data: string;
  duration: number;
  scripts?: PodcastScript[];
}

export interface TaskAi {
  platform: string;
  model: string;
}

export interface PodcastScript {
  content: string;
  speaker: string;
  emotion: string;
  speechRate: number;
  volume: number;
}

// createTask to create podcast task
export async function createTask(language: string, news: NewsDetail, voiceIds?: string[]) {
  const request = new dto.CreateTaskRequest({ language, news, voiceIds });

  return await call<CreateTaskResult>(CreateTask(request));
}

// deleteTask to delete podcast task
export async function deleteTask(batchNo: string) {
  return await call(DeleteTask({ batchNo }));
}

// queryTask to query podcast task
export async function queryTask(data: QueryTaskRequest) {
  const request = new dto.QueryTaskRequest(data);
  return await call<QueryTaskResult>(QueryTasks(request));
}

// getTask to get podcast task
export async function getTask(data: dto.GetTaskRequest) {
  return await call<PodcastTask>(GetTask(data));
}

// restyleArticle to restyle podcast
export async function restyleArticle(stageId: number, prompt: string) {
  return await call(RestyleArticle({ stageId, prompt }));
}

// mergeArticle to merge podcast
export async function mergeArticle(language: string, data: MergeArticleRequest) {
  const request = new dto.MergeArticleRequest({ language, ...data });
  return await call<CreateTaskResult>(MergeArticle(request));
}

// createScript to create podcast script
export async function createScript(stageId: number, voiceIds: string[]) {
  return await call(CreateScript({ stageId, voiceIds }));
}

export async function editScript(stageId: number, scripts: PodcastScript[]) {
  const request = new dto.EditScriptRequest({ stageId, scripts });
  return await call(EditScript(request));
}

// createAudio to create podcast audio
export async function createAudio(stageId: number) {
  return await call(CreateAudio({ stageId }));
}

// downloadAudio to download podcast audio
export async function downloadAudio(stageId: number) {
  return await call(CreateAudio({ stageId }));
}
