import { call, post, useRemoteService } from "@/utils/http";
import { QueryNews, GetNewsDetail, DeleteNews } from "wailsjs/go/adapter/App";
import { dto, httpx } from "wailsjs/go/models";

export interface QueryNewsRequest {
  recordId: number;
  pagination: httpx.Pagination;
}

export interface QueryNewsResult {
  data: NewsDetail[];
  total: number;
}

export interface NewsDetail {
  id: number;
  title: string;
  link: string;
  contents: string[];
  images: string[];
  publishedAt: string;
}

export interface GetNewsDetailRequest {
  id: number;
}

export interface DeleteNewsRequest {
  id: number;
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
