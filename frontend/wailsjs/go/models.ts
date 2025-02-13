export namespace dto {
	
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
	export class GetSystemConfigRequest {
	    key: string;
	
	    static createFrom(source: any = {}) {
	        return new GetSystemConfigRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	    }
	}
	export class QueryCrawlingRecordsRequest {
	    page?: number;
	    limit?: number;
	
	    static createFrom(source: any = {}) {
	        return new QueryCrawlingRecordsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.page = source["page"];
	        this.limit = source["limit"];
	    }
	}
	export class QueryNewsRequest {
	    recordId: number;
	    pagination?: httpx.Pagination;
	
	    static createFrom(source: any = {}) {
	        return new QueryNewsRequest(source);
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
	export class SystemConfig {
	    key: string;
	    value: any;
	
	    static createFrom(source: any = {}) {
	        return new SystemConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = source["key"];
	        this.value = source["value"];
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
	    message?: string;
	    result?: any;
	
	    static createFrom(source: any = {}) {
	        return new Response(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	        this.result = source["result"];
	    }
	}

}

