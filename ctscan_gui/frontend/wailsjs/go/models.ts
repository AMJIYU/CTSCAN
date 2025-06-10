export namespace pkg {
	
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
	export class LoginSuccess {
	    time: string;
	    event_id: string;
	    login_type: string;
	    source_ip: string;
	    username: string;
	    workstation: string;
	    subject_user: string;
	    subject_domain: string;
	    process: string;
	
	    static createFrom(source: any = {}) {
	        return new LoginSuccess(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time = source["time"];
	        this.event_id = source["event_id"];
	        this.login_type = source["login_type"];
	        this.source_ip = source["source_ip"];
	        this.username = source["username"];
	        this.workstation = source["workstation"];
	        this.subject_user = source["subject_user"];
	        this.subject_domain = source["subject_domain"];
	        this.process = source["process"];
	    }
	}
	export class NetworkConn {
	    proto: string;
	    local_addr: string;
	    remote_addr: string;
	    status: string;
	    pid: number;
	
	    static createFrom(source: any = {}) {
	        return new NetworkConn(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.proto = source["proto"];
	        this.local_addr = source["local_addr"];
	        this.remote_addr = source["remote_addr"];
	        this.status = source["status"];
	        this.pid = source["pid"];
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
	export class PatchInfo {
	    name: string;
	    date: string;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new PatchInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.date = source["date"];
	        this.status = source["status"];
	    }
	}
	export class ProcInfo {
	    pid: number;
	    name: string;
	    ppid: number;
	    parent_name: string;
	    create_time: number;
	    exe: string;
	    file_ctime: number;
	    file_mtime: number;
	    md5: string;
	    signature: string;
	
	    static createFrom(source: any = {}) {
	        return new ProcInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.pid = source["pid"];
	        this.name = source["name"];
	        this.ppid = source["ppid"];
	        this.parent_name = source["parent_name"];
	        this.create_time = source["create_time"];
	        this.exe = source["exe"];
	        this.file_ctime = source["file_ctime"];
	        this.file_mtime = source["file_mtime"];
	        this.md5 = source["md5"];
	        this.signature = source["signature"];
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
	export class SystemUser {
	    username: string;
	    uid: string;
	    gid: string;
	    home_dir: string;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemUser(source);
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

