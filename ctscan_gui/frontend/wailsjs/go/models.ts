export namespace main {
	
	export class Stats {
	    total: number;
	    network: number;
	    startup: number;
	    tasks: number;
	    process: number;
	
	    static createFrom(source: any = {}) {
	        return new Stats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total = source["total"];
	        this.network = source["network"];
	        this.startup = source["startup"];
	        this.tasks = source["tasks"];
	        this.process = source["process"];
	    }
	}

}

