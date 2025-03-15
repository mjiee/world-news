import { useRemoteService } from "@/stores";
import { call, post } from "@/utils/http";
import { QueryNews, GetNewsDetail, DeleteNews, CritiqueNews, TranslateNews } from "wailsjs/go/adapter/App";
import { dto, httpx } from "wailsjs/go/models";

interface QueryNewsRequest {
  recordId: number;
  source: string;
  topic: string;
  pagination: httpx.Pagination;
}

interface QueryNewsResult {
  data: NewsDetail[];
  total: number;
}

export interface NewsDetail {
  id: number;
  title: string;
  source: string;
  topic: string;
  link: string;
  contents: string[];
  images: string[];
  publishedAt: string;
}

interface GetNewsDetailRequest {
  id: number;
}

interface DeleteNewsRequest {
  id: number;
}

interface CritiqueNewsRequest {
  id: number;
}

interface TranslateNewsRequest {
  id: number;
  texts: string[];
  toLang: string;
}

// queryNews to query news
export async function queryNews(data: QueryNewsRequest) {
  const request = new dto.QueryNewsRequest(data);

  if (useRemoteService()) return await post<dto.QueryNewsRequest, QueryNewsResult>("/api/news/query", request);

  return await call<QueryNewsResult>(QueryNews(request));
}

// getNewsDetail to get news detail
export async function getNewsDetail(data: GetNewsDetailRequest) {
  if (useRemoteService()) return await post<GetNewsDetailRequest, NewsDetail>("/api/news/detail", data);

  return await call<NewsDetail>(GetNewsDetail(data));
}

// deleteNewsDetail to delete news detail
export async function deleteNews(data: DeleteNewsRequest) {
  if (useRemoteService()) return await post<DeleteNewsRequest, any>("/api/news/delete", data);

  return await call(DeleteNews(data));
}

// critiqueNews to critique news
export async function critiqueNews(data: CritiqueNewsRequest) {
  if (useRemoteService()) return await post<CritiqueNewsRequest, string[]>("/api/news/critique", data);

  return await call<string[]>(CritiqueNews(data));
}

// translateNews to translate news
export async function translateNews(data: TranslateNewsRequest) {
  if (useRemoteService()) return await post<TranslateNewsRequest, string[]>("/api/news/translate", data);

  return await call<string[]>(TranslateNews(data));
}
