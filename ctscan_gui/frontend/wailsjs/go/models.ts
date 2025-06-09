export namespace main {
	
	export class CronTask {
	    line: string;
	
	    static createFrom(source: any = {}) {
	        return new CronTask(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.line = source["line"];
	    }
	}
	export class NetworkInfo {
	    hostname: string;
	    ips: string[];
	    macs: string[];
	    interfaces: string[];
	
	    static createFrom(source: any = {}) {
	        return new NetworkInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hostname = source["hostname"];
	        this.ips = source["ips"];
	        this.macs = source["macs"];
	        this.interfaces = source["interfaces"];
	    }
	}
	export class ProcInfo {
	    pid: number;
	    name: string;
	    cmdline: string;
	    user: string;
	    exe: string;
	
	    static createFrom(source: any = {}) {
	        return new ProcInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pid = source["pid"];
	        this.name = source["name"];
	        this.cmdline = source["cmdline"];
	        this.user = source["user"];
	        this.exe = source["exe"];
	    }
	}
	export class StartupItem {
	    name: string;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new StartupItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	    }
	}
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
	export class SystemInfo {
	    hostname: string;
	    os: string;
	    arch: string;
	    cpu_cores: number;
	    go_version: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hostname = source["hostname"];
	        this.os = source["os"];
	        this.arch = source["arch"];
	        this.cpu_cores = source["cpu_cores"];
	        this.go_version = source["go_version"];
	    }
	}
	export class UserInfo {
	    username: string;
	    uid: string;
	    gid: string;
	    home_dir: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new UserInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.uid = source["uid"];
	        this.gid = source["gid"];
	        this.home_dir = source["home_dir"];
	        this.name = source["name"];
	    }
	}

}

