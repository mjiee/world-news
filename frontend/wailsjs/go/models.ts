export namespace dto {
	
	export class CrawlingNewsRequest {
	    startTime?: string;
	    sources?: string[];
	    topics?: string[];
	
	    static createFrom(source: any = {}) {
	        return new CrawlingNewsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.startTime = source["startTime"];
	        this.sources = source["sources"];
	        this.topics = source["topics"];
	    }
	}
	export class CreateAudioRequest {
	    stageId: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateAudioRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stageId = source["stageId"];
	    }
	}
	export class CreateScriptRequest {
	    stageId: number;
	    voiceIds: string[];
	
	    static createFrom(source: any = {}) {
	        return new CreateScriptRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stageId = source["stageId"];
	        this.voiceIds = source["voiceIds"];
	    }
	}
	export class NewsDetail {
	    id: number;
	    title: string;
	    source: string;
	    topic?: string;
	    link?: string;
	    contents?: string[];
	    images?: string[];
	    publishedAt?: string;
	    favorited?: boolean;
	
	    static createFrom(source: any = {}) {
	        return new NewsDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.source = source["source"];
	        this.topic = source["topic"];
	        this.link = source["link"];
	        this.contents = source["contents"];
	        this.images = source["images"];
	        this.publishedAt = source["publishedAt"];
	        this.favorited = source["favorited"];
	    }
	}
	export class CreateTaskRequest {
	    language: string;
	    news?: NewsDetail;
	    voiceIds?: string[];
	
	    static createFrom(source: any = {}) {
	        return new CreateTaskRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.language = source["language"];
	        this.news = this.convertValues(source["news"], NewsDetail);
	        this.voiceIds = source["voiceIds"];
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
	export class CritiqueNewsRequest {
	    contents: string[];
	
	    static createFrom(source: any = {}) {
	        return new CritiqueNewsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.contents = source["contents"];
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
	export class DeleteTaskRequest {
	    batchNo: string;
	
	    static createFrom(source: any = {}) {
	        return new DeleteTaskRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.batchNo = source["batchNo"];
	    }
	}
	export class DownloadAudioRequest {
	    stageId: number;
	    fileName?: string;
	
	    static createFrom(source: any = {}) {
	        return new DownloadAudioRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stageId = source["stageId"];
	        this.fileName = source["fileName"];
	    }
	}
	export class EditScriptRequest {
	    stageId: number;
	    scripts: ttsai.TtsScript[];
	
	    static createFrom(source: any = {}) {
	        return new EditScriptRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stageId = source["stageId"];
	        this.scripts = this.convertValues(source["scripts"], ttsai.TtsScript);
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
	export class GetCrawlingRecordRequest {
	    id: number;
	
	    static createFrom(source: any = {}) {
	        return new GetCrawlingRecordRequest(source);
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
	export class GetTaskRequest {
	    batchNo: string;
	
	    static createFrom(source: any = {}) {
	        return new GetTaskRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.batchNo = source["batchNo"];
	    }
	}
	export class MergeArticleRequest {
	    language: string;
	    title: string;
	    stageIds: number[];
	    voiceIds?: string[];
	
	    static createFrom(source: any = {}) {
	        return new MergeArticleRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.language = source["language"];
	        this.title = source["title"];
	        this.stageIds = source["stageIds"];
	        this.voiceIds = source["voiceIds"];
	    }
	}
	
	export class NewsHasTaskRequest {
	    newsId: number;
	
	    static createFrom(source: any = {}) {
	        return new NewsHasTaskRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.newsId = source["newsId"];
	    }
	}
	export class QueryCrawlingRecordsRequest {
	    recordType?: string;
	    status?: string;
	    pagination?: httpx.Pagination;
	
	    static createFrom(source: any = {}) {
	        return new QueryCrawlingRecordsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.recordType = source["recordType"];
	        this.status = source["status"];
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
	export class QueryNewsRequest {
	    recordId?: number;
	    source?: string;
	    topic?: string;
	    publishDate?: string;
	    favorited?: boolean;
	    pagination?: httpx.Pagination;
	
	    static createFrom(source: any = {}) {
	        return new QueryNewsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.recordId = source["recordId"];
	        this.source = source["source"];
	        this.topic = source["topic"];
	        this.publishDate = source["publishDate"];
	        this.favorited = source["favorited"];
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
	export class QueryTaskRequest {
	    startDate?: string;
	    endDate?: string;
	    pagination?: httpx.Pagination;
	
	    static createFrom(source: any = {}) {
	        return new QueryTaskRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.startDate = source["startDate"];
	        this.endDate = source["endDate"];
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
	export class RestyleArticleRequest {
	    stageId: number;
	    prompt: string;
	
	    static createFrom(source: any = {}) {
	        return new RestyleArticleRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stageId = source["stageId"];
	        this.prompt = source["prompt"];
	    }
	}
	export class SaveNewsFavoriteRequest {
	    id: number;
	    favorited: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SaveNewsFavoriteRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.favorited = source["favorited"];
	    }
	}
	export class SaveWebsiteWeightRequest {
	    website: string;
	    step: number;
	
	    static createFrom(source: any = {}) {
	        return new SaveWebsiteWeightRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.website = source["website"];
	        this.step = source["step"];
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
	export class TranslateNewsRequest {
	    contents: string[];
	    toLang: string;
	
	    static createFrom(source: any = {}) {
	        return new TranslateNewsRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.contents = source["contents"];
	        this.toLang = source["toLang"];
	    }
	}
	export class UpdateCrawlingRecordStatusRequest {
	    id: number;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateCrawlingRecordStatusRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.status = source["status"];
	    }
	}
	export class UpdateTaskOutputRequest {
	    stageId: number;
	    output: string;
	
	    static createFrom(source: any = {}) {
	        return new UpdateTaskOutputRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stageId = source["stageId"];
	        this.output = source["output"];
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

export namespace ttsai {
	
	export class TtsScript {
	    content: string;
	    speaker: string;
	    emotion: string;
	    speechRate: number;
	    volume: number;
	
	    static createFrom(source: any = {}) {
	        return new TtsScript(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.content = source["content"];
	        this.speaker = source["speaker"];
	        this.emotion = source["emotion"];
	        this.speechRate = source["speechRate"];
	        this.volume = source["volume"];
	    }
	}

}

