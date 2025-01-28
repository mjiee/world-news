export namespace dto {
	
	export class AddNewsKeywordRequest {
	    keyword: string;
	
	    static createFrom(source: any = {}) {
	        return new AddNewsKeywordRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.keyword = source["keyword"];
	    }
	}
	export class AddNewsWebsiteRequest {
	    websiteType: string;
	    url: string;
	    config: string[];
	
	    static createFrom(source: any = {}) {
	        return new AddNewsWebsiteRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.websiteType = source["websiteType"];
	        this.url = source["url"];
	        this.config = source["config"];
	    }
	}
	export class CrawlingNewsRequest {
	    startTime: string;
	
	    static createFrom(source: any = {}) {
	        return new CrawlingNewsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.startTime = source["startTime"];
	    }
	}
	export class DeleteCrawlingRecordRequest {
	    id: number;
	
	    static createFrom(source: any = {}) {
	        return new DeleteCrawlingRecordRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	    }
	}
	export class DeleteNewsKeywordRequest {
	    keyword: number;
	
	    static createFrom(source: any = {}) {
	        return new DeleteNewsKeywordRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.keyword = source["keyword"];
	    }
	}
	export class DeleteNewsRequest {
	    id: number;
	
	    static createFrom(source: any = {}) {
	        return new DeleteNewsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	    }
	}
	export class DeleteNewsWebsiteRequest {
	    id: number;
	
	    static createFrom(source: any = {}) {
	        return new DeleteNewsWebsiteRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	    }
	}
	export class CrawlingRecord {
	    id: number;
	    date: string;
	    quantity: number;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new CrawlingRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.date = source["date"];
	        this.quantity = source["quantity"];
	        this.status = source["status"];
	    }
	}
	export class GetCrawlingRecordResult {
	    data: CrawlingRecord[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new GetCrawlingRecordResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], CrawlingRecord);
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GetCrawlingRecordsRequest {
	    page?: number;
	    limit?: number;
	
	    static createFrom(source: any = {}) {
	        return new GetCrawlingRecordsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.page = source["page"];
	        this.limit = source["limit"];
	    }
	}
	export class GetCrawlingRecordsResponse {
	    Status: string;
	    StatusCode: number;
	    Proto: string;
	    ProtoMajor: number;
	    ProtoMinor: number;
	    Header: {[key: string]: string[]};
	    Body: any;
	    ContentLength: number;
	    TransferEncoding: string[];
	    Close: boolean;
	    Uncompressed: boolean;
	    Trailer: {[key: string]: string[]};
	    // Go type: http
	    Request?: any;
	    // Go type: tls
	    TLS?: any;
	    result?: GetCrawlingRecordResult;
	
	    static createFrom(source: any = {}) {
	        return new GetCrawlingRecordsResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Status = source["Status"];
	        this.StatusCode = source["StatusCode"];
	        this.Proto = source["Proto"];
	        this.ProtoMajor = source["ProtoMajor"];
	        this.ProtoMinor = source["ProtoMinor"];
	        this.Header = source["Header"];
	        this.Body = source["Body"];
	        this.ContentLength = source["ContentLength"];
	        this.TransferEncoding = source["TransferEncoding"];
	        this.Close = source["Close"];
	        this.Uncompressed = source["Uncompressed"];
	        this.Trailer = source["Trailer"];
	        this.Request = this.convertValues(source["Request"], null);
	        this.TLS = this.convertValues(source["TLS"], null);
	        this.result = this.convertValues(source["result"], GetCrawlingRecordResult);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GetNewsDetailRequest {
	    id: number;
	
	    static createFrom(source: any = {}) {
	        return new GetNewsDetailRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	    }
	}
	export class NewsDetail {
	    id: number;
	    title: string;
	    link: string;
	    contents: string[];
	    images: string[];
	    publishedAt: string;
	
	    static createFrom(source: any = {}) {
	        return new NewsDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.link = source["link"];
	        this.contents = source["contents"];
	        this.images = source["images"];
	        this.publishedAt = source["publishedAt"];
	    }
	}
	export class GetNewsDetailResponse {
	    code: number;
	    message: string;
	    result?: NewsDetail;
	
	    static createFrom(source: any = {}) {
	        return new GetNewsDetailResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	        this.result = this.convertValues(source["result"], NewsDetail);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GetNewsKeywordsResponse {
	    code: number;
	    message: string;
	    result: string[];
	
	    static createFrom(source: any = {}) {
	        return new GetNewsKeywordsResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	        this.result = source["result"];
	    }
	}
	export class GetNewsListRequest {
	    recordId: number;
	    pagination?: httpx.Pagination;
	
	    static createFrom(source: any = {}) {
	        return new GetNewsListRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.recordId = source["recordId"];
	        this.pagination = this.convertValues(source["pagination"], httpx.Pagination);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GetNewsListResult {
	    data: NewsDetail[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new GetNewsListResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], NewsDetail);
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GetNewsListResponse {
	    code: number;
	    message: string;
	    result?: GetNewsListResult;
	
	    static createFrom(source: any = {}) {
	        return new GetNewsListResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	        this.result = this.convertValues(source["result"], GetNewsListResult);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class GetNewsWebsitesRequest {
	    websiteType: string;
	    pagination?: httpx.Pagination;
	
	    static createFrom(source: any = {}) {
	        return new GetNewsWebsitesRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.websiteType = source["websiteType"];
	        this.pagination = this.convertValues(source["pagination"], httpx.Pagination);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class NewsWebsite {
	    id: number;
	    websiteType: string;
	    url: string;
	    config: string[];
	
	    static createFrom(source: any = {}) {
	        return new NewsWebsite(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.websiteType = source["websiteType"];
	        this.url = source["url"];
	        this.config = source["config"];
	    }
	}
	export class GetNewsWebsitesResult {
	    data: NewsWebsite[];
	    total: number;
	
	    static createFrom(source: any = {}) {
	        return new GetNewsWebsitesResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.data = this.convertValues(source["data"], NewsWebsite);
	        this.total = source["total"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class GetNewsWebsitesResponse {
	    code: number;
	    message: string;
	    result?: GetNewsWebsitesResult;
	
	    static createFrom(source: any = {}) {
	        return new GetNewsWebsitesResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	        this.result = this.convertValues(source["result"], GetNewsWebsitesResult);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

export namespace httpx {
	
	export class Pagination {
	    cursor?: number;
	    limit?: number;
	    page?: number;
	    total?: number;
	
	    static createFrom(source: any = {}) {
	        return new Pagination(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cursor = source["cursor"];
	        this.limit = source["limit"];
	        this.page = source["page"];
	        this.total = source["total"];
	    }
	}
	export class Response {
	    code: number;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	    }
	}

}

