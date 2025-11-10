export namespace core {
	
	export class Clip {
	    name: string;
	    startTime: string;
	    endTime: string;
	
	    static createFrom(source: any = {}) {
	        return new Clip(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.startTime = source["startTime"];
	        this.endTime = source["endTime"];
	    }
	}

}

