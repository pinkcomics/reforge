export namespace main {
	
	export class ScanResult {
	    cbrCount: number;
	    cbzCount: number;
	    cbrFiles: string[];
	    cbzFiles: string[];
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new ScanResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.cbrCount = source["cbrCount"];
	        this.cbzCount = source["cbzCount"];
	        this.cbrFiles = source["cbrFiles"];
	        this.cbzFiles = source["cbzFiles"];
	        this.error = source["error"];
	    }
	}

}

