import { call, useRemoteService } from "@/utils/http";
import { CrawlingNews, CrawlingWebsite, QueryCrawlingRecords, DeleteCrawlingRecord } from "wailsjs/go/adapter/App";

interface CrawlingNewsRequest {
  startTime: string;
}

interface QueryCrawlingRecordsRequest {
  page: number;
  limit: number;
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
  if (useRemoteService()) return;

  return await call<QueryCrawlingRecordResult>(QueryCrawlingRecords(data));
}

// deleteCrawlingRecord to delete crawling record
export async function deleteCrawlingRecord(data: DeleteCrawlingRecordRequest) {
  if (useRemoteService()) return;

  return await call(DeleteCrawlingRecord(data));
}
