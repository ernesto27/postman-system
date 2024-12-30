export namespace main {
	
	export class RequestParams {
	    method: string;
	    url: string;
	    headers: {[key: string]: string};
	    body: string;
	    formData: {[key: string]: string};
	    bodyType: string;
	
	    static createFrom(source: any = {}) {
	        return new RequestParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.method = source["method"];
	        this.url = source["url"];
	        this.headers = source["headers"];
	        this.body = source["body"];
	        this.formData = source["formData"];
	        this.bodyType = source["bodyType"];
	    }
	}

}

