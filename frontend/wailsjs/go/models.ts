export namespace main {
	
	export class GreetRequest {
	    msg: string;
	
	    static createFrom(source: any = {}) {
	        return new GreetRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.msg = source["msg"];
	    }
	}
	export class GreetResponse {
	    msg: string;
	
	    static createFrom(source: any = {}) {
	        return new GreetResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.msg = source["msg"];
	    }
	}

}

