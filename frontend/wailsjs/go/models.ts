export namespace main {
	
	export class Collection {
	    user_id: string;
	    collection_id: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Collection(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.user_id = source["user_id"];
	        this.collection_id = source["collection_id"];
	        this.name = source["name"];
	    }
	}
	export class FileData {
	    fieldName: string;
	    fileName: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new FileData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.fieldName = source["fieldName"];
	        this.fileName = source["fileName"];
	        this.content = source["content"];
	    }
	}
	export class RequestParams {
	    method: string;
	    url: string;
	    headers: {[key: string]: string};
	    body: string;
	    formData: {[key: string]: string};
	    files: FileData[];
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
	        this.files = this.convertValues(source["files"], FileData);
	        this.bodyType = source["bodyType"];
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

