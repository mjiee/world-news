import { call, useRemoteService } from "@/utils/http";
import {
  CrawlingNews,
  CrawlingWebsite,
  QueryCrawlingRecords,
  DeleteCrawlingRecord,
  HasCrawlingTasks,
} from "wailsjs/go/adapter/App";
import { dto, httpx } from "wailsjs/go/models";

interface CrawlingNewsRequest {
  startTime: string;
}

interface QueryCrawlingRecordsRequest {
  recordType?: string;
  status?: string;
  pagination: httpx.Pagination;
}

interface DeleteCrawlingRecordRequest {
  id: number;
}

interface QueryCrawlingRecordResult {
  data: CrawlingRecord[];
  total: number;
}

export interface CrawlingRecord {
  id: number;
  recordType: string;
  date: string;
  quantity: number;
  status: string;
}

// crawlingNews to crawl news
export async function crawlingNews(data: CrawlingNewsRequest) {
  if (useRemoteService()) return;

  return await call(CrawlingNews(data));
}

// crawlingWebsite to crawl website
export async function crawlingWebsite() {
  if (useRemoteService()) return;

  return await call(CrawlingWebsite());
}

// queryCrawlingRecords to query crawling records
export async function queryCrawlingRecords(data: QueryCrawlingRecordsRequest) {
  const request = new dto.QueryCrawlingRecordsRequest(data);

  if (useRemoteService()) return;

  return await call<QueryCrawlingRecordResult>(QueryCrawlingRecords(request));
}

// deleteCrawlingRecord to delete crawling record
export async function deleteCrawlingRecord(data: DeleteCrawlingRecordRequest) {
  if (useRemoteService()) return;

  return await call(DeleteCrawlingRecord(data));
}

// hasCrawlingTask to check if has crawling task
export async function hasCrawlingTask() {
  if (useRemoteService()) return;

  return await call<boolean>(HasCrawlingTasks());
}
