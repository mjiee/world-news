// Cynhyrchwyd y ffeil hon yn awtomatig. PEIDIWCH Â MODIWL
// This file is automatically generated. DO NOT EDIT
import {dto} from '../models';
import {httpx} from '../models';

export function CrawlingNews(arg1:dto.CrawlingNewsRequest):Promise<httpx.Response>;

export function DeleteCrawlingRecord(arg1:dto.DeleteCrawlingRecordRequest):Promise<httpx.Response>;

export function DeleteNewsDetail(arg1:dto.DeleteNewsRequest):Promise<httpx.Response>;

export function GetNewsDetail(arg1:dto.GetNewsDetailRequest):Promise<dto.GetNewsDetailResponse>;

export function GetSystemConfig(arg1:dto.GetSystemConfigRequest):Promise<dto.GetSystemConfigResponse>;

export function QueryCrawlingRecords(arg1:dto.QueryCrawlingRecordsRequest):Promise<dto.QueryCrawlingRecordsResponse>;

export function QueryNews(arg1:dto.QueryNewsRequest):Promise<dto.QueryNewsResponse>;

export function SaveSystemConfig(arg1:dto.SystemConfig):Promise<httpx.Response>;
