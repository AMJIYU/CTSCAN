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
	export class DiskInfo {
	    mount_point: string;
	    total_size: number;
	    used_size: number;
	    free_size: number;
	    usage: number;
	
	    static createFrom(source: any = {}) {
	        return new DiskInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mount_point = source["mount_point"];
	        this.total_size = source["total_size"];
	        this.used_size = source["used_size"];
	        this.free_size = source["free_size"];
	        this.usage = source["usage"];
	    }
	}
	export class EVTXEvent {
	    time: string;
	    event_id: number;
	    provider: string;
	    level: string;
	    channel: string;
	    computer: string;
	    user_id: string;
	    description: string;
	    data: Record<string, any>;
	    event_record_id: number;
	    version: number;
	    qualifiers: number;
	    task: number;
	    opcode: number;
	    keywords: string;
	    process_id: number;
	    thread_id: number;
	    message: string;
	    system_info: Record<string, any>;
	    event_data: Record<string, any>;
	    user_data: Record<string, any>;
	
	    static createFrom(source: any = {}) {
	        return new EVTXEvent(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time = source["time"];
	        this.event_id = source["event_id"];
	        this.provider = source["provider"];
	        this.level = source["level"];
	        this.channel = source["channel"];
	        this.computer = source["computer"];
	        this.user_id = source["user_id"];
	        this.description = source["description"];
	        this.data = source["data"];
	        this.event_record_id = source["event_record_id"];
	        this.version = source["version"];
	        this.qualifiers = source["qualifiers"];
	        this.task = source["task"];
	        this.opcode = source["opcode"];
	        this.keywords = source["keywords"];
	        this.process_id = source["process_id"];
	        this.thread_id = source["thread_id"];
	        this.message = source["message"];
	        this.system_info = source["system_info"];
	        this.event_data = source["event_data"];
	        this.user_data = source["user_data"];
	    }
	}
	export class FileInfo {
	    path: string;
	    exists: boolean;
	    size: number;
	    mode: string;
	    // Go type: time
	    mod_time: any;
	    // Go type: time
	    create_time: any;
	    // Go type: time
	    access_time: any;
	    // Go type: time
	    change_time: any;
	    is_dir: boolean;
	    is_symlink: boolean;
	    owner: string;
	    group: string;
	    permissions: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.exists = source["exists"];
	        this.size = source["size"];
	        this.mode = source["mode"];
	        this.mod_time = this.convertValues(source["mod_time"], null);
	        this.create_time = this.convertValues(source["create_time"], null);
	        this.access_time = this.convertValues(source["access_time"], null);
	        this.change_time = this.convertValues(source["change_time"], null);
	        this.is_dir = source["is_dir"];
	        this.is_symlink = source["is_symlink"];
	        this.owner = source["owner"];
	        this.group = source["group"];
	        this.permissions = source["permissions"];
	        this.description = source["description"];
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
	export class InterfaceStats {
	    name: string;
	    bytes_sent: number;
	    bytes_recv: number;
	    packets_sent: number;
	    packets_recv: number;
	
	    static createFrom(source: any = {}) {
	        return new InterfaceStats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.bytes_sent = source["bytes_sent"];
	        this.bytes_recv = source["bytes_recv"];
	        this.packets_sent = source["packets_sent"];
	        this.packets_recv = source["packets_recv"];
	    }
	}
	export class LoginFailed {
	    time: string;
	    event_id: string;
	    event_type: string;
	    source: string;
	    username: string;
	    ip_address: string;
	    reason: string;
	
	    static createFrom(source: any = {}) {
	        return new LoginFailed(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time = source["time"];
	        this.event_id = source["event_id"];
	        this.event_type = source["event_type"];
	        this.source = source["source"];
	        this.username = source["username"];
	        this.ip_address = source["ip_address"];
	        this.reason = source["reason"];
	    }
	}
	export class LoginSuccess {
	    time: string;
	    event_id: string;
	    event_type: string;
	    source: string;
	    username: string;
	    ip_address: string;
	
	    static createFrom(source: any = {}) {
	        return new LoginSuccess(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time = source["time"];
	        this.event_id = source["event_id"];
	        this.event_type = source["event_type"];
	        this.source = source["source"];
	        this.username = source["username"];
	        this.ip_address = source["ip_address"];
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
	    interface_stats: InterfaceStats[];
	    gateway: string;
	
	    static createFrom(source: any = {}) {
	        return new NetworkInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hostname = source["hostname"];
	        this.ips = source["ips"];
	        this.macs = source["macs"];
	        this.interfaces = source["interfaces"];
	        this.interface_stats = this.convertValues(source["interface_stats"], InterfaceStats);
	        this.gateway = source["gateway"];
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
	    cpu_percent: number;
	    mem_percent: number;
	
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
	        this.cpu_percent = source["cpu_percent"];
	        this.mem_percent = source["mem_percent"];
	    }
	}
	export class RDPLoginInfo {
	    time: string;
	    username: string;
	    ip: string;
	    status: string;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new RDPLoginInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time = source["time"];
	        this.username = source["username"];
	        this.ip = source["ip"];
	        this.status = source["status"];
	        this.description = source["description"];
	    }
	}
	export class ShellHistory {
	    time: string;
	    command: string;
	    user: string;
	    shell: string;
	
	    static createFrom(source: any = {}) {
	        return new ShellHistory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time = source["time"];
	        this.command = source["command"];
	        this.user = source["user"];
	        this.shell = source["shell"];
	    }
	}
	export class StartupItem {
	    name: string;
	    path: string;
	    type: string;
	    enabled: boolean;
	    // Go type: time
	    lastModTime: any;
	    size: number;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new StartupItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.type = source["type"];
	        this.enabled = source["enabled"];
	        this.lastModTime = this.convertValues(source["lastModTime"], null);
	        this.size = source["size"];
	        this.description = source["description"];
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
	export class SystemInfo {
	    hostname: string;
	    os: string;
	    arch: string;
	    cpu_cores: number;
	    kernel_version: string;
	    cpu_usage: number;
	    total_memory: number;
	    memory_usage: number;
	    disks: DiskInfo[];
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.hostname = source["hostname"];
	        this.os = source["os"];
	        this.arch = source["arch"];
	        this.cpu_cores = source["cpu_cores"];
	        this.kernel_version = source["kernel_version"];
	        this.cpu_usage = source["cpu_usage"];
	        this.total_memory = source["total_memory"];
	        this.memory_usage = source["memory_usage"];
	        this.disks = this.convertValues(source["disks"], DiskInfo);
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

