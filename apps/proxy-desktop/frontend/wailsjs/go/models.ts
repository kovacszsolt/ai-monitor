export namespace config {
	
	export class StorageConfig {
	    DatabasePath: string;
	    JournalMode: string;
	    BusyTimeoutMS: number;
	    Synchronous: string;
	
	    static createFrom(source: any = {}) {
	        return new StorageConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.DatabasePath = source["DatabasePath"];
	        this.JournalMode = source["JournalMode"];
	        this.BusyTimeoutMS = source["BusyTimeoutMS"];
	        this.Synchronous = source["Synchronous"];
	    }
	}
	export class Config {
	    Storage: StorageConfig;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Storage = this.convertValues(source["Storage"], StorageConfig);
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
	export class Bridge {
	    Desktop: Config;
	    ProxyPath: string;
	    ProxyFile: Config;
	    ProxyMerged: Config;
	    BridgeMessage: string;
	
	    static createFrom(source: any = {}) {
	        return new Bridge(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Desktop = this.convertValues(source["Desktop"], Config);
	        this.ProxyPath = source["ProxyPath"];
	        this.ProxyFile = this.convertValues(source["ProxyFile"], Config);
	        this.ProxyMerged = this.convertValues(source["ProxyMerged"], Config);
	        this.BridgeMessage = source["BridgeMessage"];
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

export namespace service {
	
	export class AIRequestCaptureErrorDTO {
	    stage: string;
	    error_code: string;
	    error_message: string;
	
	    static createFrom(source: any = {}) {
	        return new AIRequestCaptureErrorDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.stage = source["stage"];
	        this.error_code = source["error_code"];
	        this.error_message = source["error_message"];
	    }
	}
	export class AIRequestContentDTO {
	    request_id: string;
	    request_state?: string;
	    proxy_error?: string;
	    capture_enabled: boolean;
	    has_content: boolean;
	    last_user_text?: string;
	    messages_json?: string;
	    prompt_truncated: boolean;
	    assistant_text?: string;
	    response_truncated: boolean;
	    tool_calls_json?: string;
	    tool_calls_truncated: boolean;
	    raw_body?: string;
	    raw_body_truncated: boolean;
	    capture_errors?: AIRequestCaptureErrorDTO[];
	
	    static createFrom(source: any = {}) {
	        return new AIRequestContentDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.request_id = source["request_id"];
	        this.request_state = source["request_state"];
	        this.proxy_error = source["proxy_error"];
	        this.capture_enabled = source["capture_enabled"];
	        this.has_content = source["has_content"];
	        this.last_user_text = source["last_user_text"];
	        this.messages_json = source["messages_json"];
	        this.prompt_truncated = source["prompt_truncated"];
	        this.assistant_text = source["assistant_text"];
	        this.response_truncated = source["response_truncated"];
	        this.tool_calls_json = source["tool_calls_json"];
	        this.tool_calls_truncated = source["tool_calls_truncated"];
	        this.raw_body = source["raw_body"];
	        this.raw_body_truncated = source["raw_body_truncated"];
	        this.capture_errors = this.convertValues(source["capture_errors"], AIRequestCaptureErrorDTO);
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
	export class AIRequestDTO {
	    id: string;
	    project_id: string;
	    project_name: string;
	    provider_id: string;
	    provider_name: string;
	    requested_model: string;
	    actual_model?: string;
	    endpoint: string;
	    protocol: string;
	    started_at: string;
	    completed_at?: string;
	    duration_ms: number;
	    has_duration: boolean;
	    ttft_ms: number;
	    has_ttft: boolean;
	    status_code: number;
	    has_status_code: boolean;
	    request_state: string;
	    streaming: boolean;
	    input_tokens: number;
	    has_input_tokens: boolean;
	    output_tokens: number;
	    has_output_tokens: boolean;
	    cached_input_tokens: number;
	    has_cached_input_tokens: boolean;
	    uncached_input_tokens: number;
	    has_uncached_input_tokens: boolean;
	    prompt_tps: number;
	    has_prompt_tps: boolean;
	    prompt_speed_estimated: boolean;
	    generation_tps: number;
	    has_generation_tps: boolean;
	    generation_speed_estimated: boolean;
	    effective_tps: number;
	    has_effective_tps: boolean;
	    cache_hit_rate: number;
	    has_cache_hit_rate: boolean;
	    usage_status: string;
	    request_bytes: number;
	    response_bytes: number;
	    error_type?: string;
	    tag?: string;
	
	    static createFrom(source: any = {}) {
	        return new AIRequestDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.project_id = source["project_id"];
	        this.project_name = source["project_name"];
	        this.provider_id = source["provider_id"];
	        this.provider_name = source["provider_name"];
	        this.requested_model = source["requested_model"];
	        this.actual_model = source["actual_model"];
	        this.endpoint = source["endpoint"];
	        this.protocol = source["protocol"];
	        this.started_at = source["started_at"];
	        this.completed_at = source["completed_at"];
	        this.duration_ms = source["duration_ms"];
	        this.has_duration = source["has_duration"];
	        this.ttft_ms = source["ttft_ms"];
	        this.has_ttft = source["has_ttft"];
	        this.status_code = source["status_code"];
	        this.has_status_code = source["has_status_code"];
	        this.request_state = source["request_state"];
	        this.streaming = source["streaming"];
	        this.input_tokens = source["input_tokens"];
	        this.has_input_tokens = source["has_input_tokens"];
	        this.output_tokens = source["output_tokens"];
	        this.has_output_tokens = source["has_output_tokens"];
	        this.cached_input_tokens = source["cached_input_tokens"];
	        this.has_cached_input_tokens = source["has_cached_input_tokens"];
	        this.uncached_input_tokens = source["uncached_input_tokens"];
	        this.has_uncached_input_tokens = source["has_uncached_input_tokens"];
	        this.prompt_tps = source["prompt_tps"];
	        this.has_prompt_tps = source["has_prompt_tps"];
	        this.prompt_speed_estimated = source["prompt_speed_estimated"];
	        this.generation_tps = source["generation_tps"];
	        this.has_generation_tps = source["has_generation_tps"];
	        this.generation_speed_estimated = source["generation_speed_estimated"];
	        this.effective_tps = source["effective_tps"];
	        this.has_effective_tps = source["has_effective_tps"];
	        this.cache_hit_rate = source["cache_hit_rate"];
	        this.has_cache_hit_rate = source["has_cache_hit_rate"];
	        this.usage_status = source["usage_status"];
	        this.request_bytes = source["request_bytes"];
	        this.response_bytes = source["response_bytes"];
	        this.error_type = source["error_type"];
	        this.tag = source["tag"];
	    }
	}
	export class AIRequestListDTO {
	    items: AIRequestDTO[];
	    total: number;
	    page: number;
	    page_size: number;
	
	    static createFrom(source: any = {}) {
	        return new AIRequestListDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.items = this.convertValues(source["items"], AIRequestDTO);
	        this.total = source["total"];
	        this.page = source["page"];
	        this.page_size = source["page_size"];
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
	export class AIRequestTagOptionsDTO {
	    enabled: boolean;
	    tags: string[];
	    has_untagged: boolean;
	    no_tag_label: string;
	
	    static createFrom(source: any = {}) {
	        return new AIRequestTagOptionsDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.tags = source["tags"];
	        this.has_untagged = source["has_untagged"];
	        this.no_tag_label = source["no_tag_label"];
	    }
	}
	export class APIKeyDTO {
	    id: string;
	    project_id: string;
	    name: string;
	    api_key: string;
	    allow_all_models: boolean;
	    enabled: boolean;
	    created_at: string;
	    updated_at: string;
	    expires_at?: string;
	    last_used_at?: string;
	
	    static createFrom(source: any = {}) {
	        return new APIKeyDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.project_id = source["project_id"];
	        this.name = source["name"];
	        this.api_key = source["api_key"];
	        this.allow_all_models = source["allow_all_models"];
	        this.enabled = source["enabled"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	        this.expires_at = source["expires_at"];
	        this.last_used_at = source["last_used_at"];
	    }
	}
	export class AboutDTO {
	    app_name: string;
	    version: string;
	    build_time: string;
	    author: string;
	    license: string;
	    config_file_path: string;
	    database_path: string;
	    proxy_config_path: string;
	    schema_ready: boolean;
	    migration_version: number;
	    migration_dirty: boolean;
	    config_bridge_ok: boolean;
	    bridge_message?: string;
	
	    static createFrom(source: any = {}) {
	        return new AboutDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.app_name = source["app_name"];
	        this.version = source["version"];
	        this.build_time = source["build_time"];
	        this.author = source["author"];
	        this.license = source["license"];
	        this.config_file_path = source["config_file_path"];
	        this.database_path = source["database_path"];
	        this.proxy_config_path = source["proxy_config_path"];
	        this.schema_ready = source["schema_ready"];
	        this.migration_version = source["migration_version"];
	        this.migration_dirty = source["migration_dirty"];
	        this.config_bridge_ok = source["config_bridge_ok"];
	        this.bridge_message = source["bridge_message"];
	    }
	}
	export class BulkDeleteAIRequestResultDTO {
	    deleted: number;
	
	    static createFrom(source: any = {}) {
	        return new BulkDeleteAIRequestResultDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.deleted = source["deleted"];
	    }
	}
	export class BulkSetAIRequestTagResultDTO {
	    updated: number;
	
	    static createFrom(source: any = {}) {
	        return new BulkSetAIRequestTagResultDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.updated = source["updated"];
	    }
	}
	export class ChatMessageDTO {
	    role: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new ChatMessageDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.role = source["role"];
	        this.content = source["content"];
	    }
	}
	export class ChatCompletionInputDTO {
	    api_key_id: string;
	    model_route_id: string;
	    messages: ChatMessageDTO[];
	    stream: boolean;
	    stream_client_id: string;
	
	    static createFrom(source: any = {}) {
	        return new ChatCompletionInputDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.api_key_id = source["api_key_id"];
	        this.model_route_id = source["model_route_id"];
	        this.messages = this.convertValues(source["messages"], ChatMessageDTO);
	        this.stream = source["stream"];
	        this.stream_client_id = source["stream_client_id"];
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
	export class ChatCompletionResultDTO {
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new ChatCompletionResultDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.content = source["content"];
	    }
	}
	
	export class ChatPinnedPromptDTO {
	    id: string;
	    title: string;
	    body: string;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new ChatPinnedPromptDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.body = source["body"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}
	export class ConnectionInfoDTO {
	    database_path: string;
	    connected: boolean;
	    schema_ready: boolean;
	    message?: string;
	    build_time: string;
	    proxy_config_path?: string;
	    config_bridge_ok: boolean;
	    bridge_message?: string;
	    migration_version: number;
	    migration_dirty: boolean;
	    proxy_running: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ConnectionInfoDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.database_path = source["database_path"];
	        this.connected = source["connected"];
	        this.schema_ready = source["schema_ready"];
	        this.message = source["message"];
	        this.build_time = source["build_time"];
	        this.proxy_config_path = source["proxy_config_path"];
	        this.config_bridge_ok = source["config_bridge_ok"];
	        this.bridge_message = source["bridge_message"];
	        this.migration_version = source["migration_version"];
	        this.migration_dirty = source["migration_dirty"];
	        this.proxy_running = source["proxy_running"];
	    }
	}
	export class CreateAPIKeyResultDTO {
	    key: APIKeyDTO;
	    plaintext: string;
	
	    static createFrom(source: any = {}) {
	        return new CreateAPIKeyResultDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.key = this.convertValues(source["key"], APIKeyDTO);
	        this.plaintext = source["plaintext"];
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
	export class UsageGroupStatsDTO {
	    name: string;
	    request_count: number;
	    input_tokens: number;
	    output_tokens: number;
	    cached_input_tokens: number;
	
	    static createFrom(source: any = {}) {
	        return new UsageGroupStatsDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.request_count = source["request_count"];
	        this.input_tokens = source["input_tokens"];
	        this.output_tokens = source["output_tokens"];
	        this.cached_input_tokens = source["cached_input_tokens"];
	    }
	}
	export class NamedCountDTO {
	    name: string;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new NamedCountDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.count = source["count"];
	    }
	}
	export class StateCountDTO {
	    state: string;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new StateCountDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.state = source["state"];
	        this.count = source["count"];
	    }
	}
	export class DashboardStatsDTO {
	    total_requests: number;
	    requests_last_24h: number;
	    sum_input_tokens: number;
	    sum_output_tokens: number;
	    sum_cached_input_tokens: number;
	    avg_duration_ms: number;
	    streaming_count: number;
	    by_state: StateCountDTO[];
	    top_models: NamedCountDTO[];
	    top_projects: NamedCountDTO[];
	    usage_by_tag: UsageGroupStatsDTO[];
	
	    static createFrom(source: any = {}) {
	        return new DashboardStatsDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_requests = source["total_requests"];
	        this.requests_last_24h = source["requests_last_24h"];
	        this.sum_input_tokens = source["sum_input_tokens"];
	        this.sum_output_tokens = source["sum_output_tokens"];
	        this.sum_cached_input_tokens = source["sum_cached_input_tokens"];
	        this.avg_duration_ms = source["avg_duration_ms"];
	        this.streaming_count = source["streaming_count"];
	        this.by_state = this.convertValues(source["by_state"], StateCountDTO);
	        this.top_models = this.convertValues(source["top_models"], NamedCountDTO);
	        this.top_projects = this.convertValues(source["top_projects"], NamedCountDTO);
	        this.usage_by_tag = this.convertValues(source["usage_by_tag"], UsageGroupStatsDTO);
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
	export class ErrorDTO {
	    code: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new ErrorDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.code = source["code"];
	        this.message = source["message"];
	    }
	}
	export class ExportPayloadDTO {
	    filename: string;
	    content_base64: string;
	    byte_count: number;
	
	    static createFrom(source: any = {}) {
	        return new ExportPayloadDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filename = source["filename"];
	        this.content_base64 = source["content_base64"];
	        this.byte_count = source["byte_count"];
	    }
	}
	export class HourlyModelRequestCountDTO {
	    bucket_start: string;
	    model: string;
	    request_count: number;
	    input_tokens: number;
	    cached_input_tokens: number;
	    output_tokens: number;
	    total_cost_usd: number;
	    has_total_cost: boolean;
	    generation_tps: number;
	    has_generation_tps: boolean;
	
	    static createFrom(source: any = {}) {
	        return new HourlyModelRequestCountDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bucket_start = source["bucket_start"];
	        this.model = source["model"];
	        this.request_count = source["request_count"];
	        this.input_tokens = source["input_tokens"];
	        this.cached_input_tokens = source["cached_input_tokens"];
	        this.output_tokens = source["output_tokens"];
	        this.total_cost_usd = source["total_cost_usd"];
	        this.has_total_cost = source["has_total_cost"];
	        this.generation_tps = source["generation_tps"];
	        this.has_generation_tps = source["has_generation_tps"];
	    }
	}
	export class HourlyModelRequestSeriesDTO {
	    bucket_unit: string;
	    days: number;
	    range_start: string;
	    range_end: string;
	    points: HourlyModelRequestCountDTO[];
	
	    static createFrom(source: any = {}) {
	        return new HourlyModelRequestSeriesDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bucket_unit = source["bucket_unit"];
	        this.days = source["days"];
	        this.range_start = source["range_start"];
	        this.range_end = source["range_end"];
	        this.points = this.convertValues(source["points"], HourlyModelRequestCountDTO);
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
	export class ModelRouteDTO {
	    id: string;
	    public_model_name: string;
	    provider_id: string;
	    provider_name?: string;
	    upstream_model_name: string;
	    protocol: string;
	    enabled: boolean;
	    input_token_usd_per_1m?: number;
	    cached_input_token_usd_per_1m?: number;
	    output_token_usd_per_1m?: number;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new ModelRouteDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.public_model_name = source["public_model_name"];
	        this.provider_id = source["provider_id"];
	        this.provider_name = source["provider_name"];
	        this.upstream_model_name = source["upstream_model_name"];
	        this.protocol = source["protocol"];
	        this.enabled = source["enabled"];
	        this.input_token_usd_per_1m = source["input_token_usd_per_1m"];
	        this.cached_input_token_usd_per_1m = source["cached_input_token_usd_per_1m"];
	        this.output_token_usd_per_1m = source["output_token_usd_per_1m"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}
	export class ModelTokenTimePointDTO {
	    bucket_start: string;
	    model: string;
	    input_tokens: number;
	    output_tokens: number;
	
	    static createFrom(source: any = {}) {
	        return new ModelTokenTimePointDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bucket_start = source["bucket_start"];
	        this.model = source["model"];
	        this.input_tokens = source["input_tokens"];
	        this.output_tokens = source["output_tokens"];
	    }
	}
	export class ModelTokenTimeSeriesDTO {
	    bucket_unit: string;
	    days: number;
	    points: ModelTokenTimePointDTO[];
	
	    static createFrom(source: any = {}) {
	        return new ModelTokenTimeSeriesDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bucket_unit = source["bucket_unit"];
	        this.days = source["days"];
	        this.points = this.convertValues(source["points"], ModelTokenTimePointDTO);
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
	export class ModelUsageStatsDTO {
	    name: string;
	    request_count: number;
	    input_tokens: number;
	    output_tokens: number;
	    cached_input_tokens: number;
	    total_cost_usd: number;
	    has_total_cost: boolean;
	    effective_tps: number;
	    generation_tps: number;
	    has_generation_tps: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ModelUsageStatsDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.request_count = source["request_count"];
	        this.input_tokens = source["input_tokens"];
	        this.output_tokens = source["output_tokens"];
	        this.cached_input_tokens = source["cached_input_tokens"];
	        this.total_cost_usd = source["total_cost_usd"];
	        this.has_total_cost = source["has_total_cost"];
	        this.effective_tps = source["effective_tps"];
	        this.generation_tps = source["generation_tps"];
	        this.has_generation_tps = source["has_generation_tps"];
	    }
	}
	
	export class ProjectDTO {
	    id: string;
	    name: string;
	    description: string;
	    enabled: boolean;
	    created_at: string;
	    updated_at: string;
	    request_count: number;
	    last_request_at?: string;
	
	    static createFrom(source: any = {}) {
	        return new ProjectDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.description = source["description"];
	        this.enabled = source["enabled"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	        this.request_count = source["request_count"];
	        this.last_request_at = source["last_request_at"];
	    }
	}
	export class ProviderDTO {
	    id: string;
	    name: string;
	    provider_type: string;
	    base_url: string;
	    has_credential: boolean;
	    enabled: boolean;
	    connection_timeout_ms: number;
	    request_timeout_ms: number;
	    created_at: string;
	    updated_at: string;
	
	    static createFrom(source: any = {}) {
	        return new ProviderDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.provider_type = source["provider_type"];
	        this.base_url = source["base_url"];
	        this.has_credential = source["has_credential"];
	        this.enabled = source["enabled"];
	        this.connection_timeout_ms = source["connection_timeout_ms"];
	        this.request_timeout_ms = source["request_timeout_ms"];
	        this.created_at = source["created_at"];
	        this.updated_at = source["updated_at"];
	    }
	}
	export class UnhealthyProviderDTO {
	    id: string;
	    name: string;
	    message: string;
	
	    static createFrom(source: any = {}) {
	        return new UnhealthyProviderDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.message = source["message"];
	    }
	}
	export class ProviderHealthCheckDTO {
	    unhealthy_providers: UnhealthyProviderDTO[];
	
	    static createFrom(source: any = {}) {
	        return new ProviderHealthCheckDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.unhealthy_providers = this.convertValues(source["unhealthy_providers"], UnhealthyProviderDTO);
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
	export class ProxyConfigSummaryDTO {
	    listen_addr: string;
	    database_path: string;
	    capture_mode: string;
	    store_content: boolean;
	    proxy_config_path: string;
	    bridge_ok: boolean;
	    bridge_message?: string;
	
	    static createFrom(source: any = {}) {
	        return new ProxyConfigSummaryDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.listen_addr = source["listen_addr"];
	        this.database_path = source["database_path"];
	        this.capture_mode = source["capture_mode"];
	        this.store_content = source["store_content"];
	        this.proxy_config_path = source["proxy_config_path"];
	        this.bridge_ok = source["bridge_ok"];
	        this.bridge_message = source["bridge_message"];
	    }
	}
	export class ProxyStatusDTO {
	    running: boolean;
	    listen_addr?: string;
	    error?: string;
	    config_path?: string;
	
	    static createFrom(source: any = {}) {
	        return new ProxyStatusDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.running = source["running"];
	        this.listen_addr = source["listen_addr"];
	        this.error = source["error"];
	        this.config_path = source["config_path"];
	    }
	}
	export class ProxyVariablesDTO {
	    server_host: string;
	    server_port: number;
	    server_max_request_bytes: number;
	    logging_level: string;
	    logging_format: string;
	    capture_mode: string;
	    capture_store_content: boolean;
	    capture_max_prompt_bytes: number;
	    capture_max_response_bytes: number;
	
	    static createFrom(source: any = {}) {
	        return new ProxyVariablesDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.server_host = source["server_host"];
	        this.server_port = source["server_port"];
	        this.server_max_request_bytes = source["server_max_request_bytes"];
	        this.logging_level = source["logging_level"];
	        this.logging_format = source["logging_format"];
	        this.capture_mode = source["capture_mode"];
	        this.capture_store_content = source["capture_store_content"];
	        this.capture_max_prompt_bytes = source["capture_max_prompt_bytes"];
	        this.capture_max_response_bytes = source["capture_max_response_bytes"];
	    }
	}
	export class Result___github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_APIKeyDTO_ {
	    ok: boolean;
	    data: APIKeyDTO[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result___github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_APIKeyDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], APIKeyDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result___github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ChatPinnedPromptDTO_ {
	    ok: boolean;
	    data: ChatPinnedPromptDTO[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result___github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ChatPinnedPromptDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], ChatPinnedPromptDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result___github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ModelRouteDTO_ {
	    ok: boolean;
	    data: ModelRouteDTO[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result___github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ModelRouteDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], ModelRouteDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result___github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ModelUsageStatsDTO_ {
	    ok: boolean;
	    data: ModelUsageStatsDTO[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result___github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ModelUsageStatsDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], ModelUsageStatsDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result___github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ProjectDTO_ {
	    ok: boolean;
	    data: ProjectDTO[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result___github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ProjectDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], ProjectDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result___github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ProviderDTO_ {
	    ok: boolean;
	    data: ProviderDTO[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result___github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ProviderDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], ProviderDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result___string_ {
	    ok: boolean;
	    data: string[];
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result___string_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = source["data"];
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_bool_ {
	    ok: boolean;
	    data: boolean;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_bool_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = source["data"];
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_AIRequestContentDTO_ {
	    ok: boolean;
	    data: AIRequestContentDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_AIRequestContentDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], AIRequestContentDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_AIRequestListDTO_ {
	    ok: boolean;
	    data: AIRequestListDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_AIRequestListDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], AIRequestListDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_AIRequestTagOptionsDTO_ {
	    ok: boolean;
	    data: AIRequestTagOptionsDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_AIRequestTagOptionsDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], AIRequestTagOptionsDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_BulkDeleteAIRequestResultDTO_ {
	    ok: boolean;
	    data: BulkDeleteAIRequestResultDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_BulkDeleteAIRequestResultDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], BulkDeleteAIRequestResultDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_BulkSetAIRequestTagResultDTO_ {
	    ok: boolean;
	    data: BulkSetAIRequestTagResultDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_BulkSetAIRequestTagResultDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], BulkSetAIRequestTagResultDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ChatCompletionResultDTO_ {
	    ok: boolean;
	    data: ChatCompletionResultDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ChatCompletionResultDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], ChatCompletionResultDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ChatPinnedPromptDTO_ {
	    ok: boolean;
	    data: ChatPinnedPromptDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ChatPinnedPromptDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], ChatPinnedPromptDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ConnectionInfoDTO_ {
	    ok: boolean;
	    data: ConnectionInfoDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ConnectionInfoDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], ConnectionInfoDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_CreateAPIKeyResultDTO_ {
	    ok: boolean;
	    data: CreateAPIKeyResultDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_CreateAPIKeyResultDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], CreateAPIKeyResultDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_DashboardStatsDTO_ {
	    ok: boolean;
	    data: DashboardStatsDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_DashboardStatsDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], DashboardStatsDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ExportPayloadDTO_ {
	    ok: boolean;
	    data: ExportPayloadDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ExportPayloadDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], ExportPayloadDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_HourlyModelRequestSeriesDTO_ {
	    ok: boolean;
	    data: HourlyModelRequestSeriesDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_HourlyModelRequestSeriesDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], HourlyModelRequestSeriesDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ModelRouteDTO_ {
	    ok: boolean;
	    data: ModelRouteDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ModelRouteDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], ModelRouteDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ModelTokenTimeSeriesDTO_ {
	    ok: boolean;
	    data: ModelTokenTimeSeriesDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ModelTokenTimeSeriesDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], ModelTokenTimeSeriesDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ProjectDTO_ {
	    ok: boolean;
	    data: ProjectDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ProjectDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], ProjectDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ProviderDTO_ {
	    ok: boolean;
	    data: ProviderDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ProviderDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], ProviderDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ProviderHealthCheckDTO_ {
	    ok: boolean;
	    data: ProviderHealthCheckDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ProviderHealthCheckDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], ProviderHealthCheckDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ProxyConfigSummaryDTO_ {
	    ok: boolean;
	    data: ProxyConfigSummaryDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ProxyConfigSummaryDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], ProxyConfigSummaryDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ProxyStatusDTO_ {
	    ok: boolean;
	    data: ProxyStatusDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ProxyStatusDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], ProxyStatusDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ProxyVariablesDTO_ {
	    ok: boolean;
	    data: ProxyVariablesDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_ProxyVariablesDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], ProxyVariablesDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class SettingsDTO {
	    config_file_path: string;
	    database_path: string;
	    effective_database_path: string;
	    auto_resolve: boolean;
	    proxy_config_path: string;
	    bridge_message?: string;
	    config_bridge_ok: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SettingsDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.config_file_path = source["config_file_path"];
	        this.database_path = source["database_path"];
	        this.effective_database_path = source["effective_database_path"];
	        this.auto_resolve = source["auto_resolve"];
	        this.proxy_config_path = source["proxy_config_path"];
	        this.bridge_message = source["bridge_message"];
	        this.config_bridge_ok = source["config_bridge_ok"];
	    }
	}
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_SettingsDTO_ {
	    ok: boolean;
	    data: SettingsDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_SettingsDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], SettingsDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class StartProxyResultDTO {
	    started: boolean;
	    unhealthy_providers: UnhealthyProviderDTO[];
	
	    static createFrom(source: any = {}) {
	        return new StartProxyResultDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.started = source["started"];
	        this.unhealthy_providers = this.convertValues(source["unhealthy_providers"], UnhealthyProviderDTO);
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
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_StartProxyResultDTO_ {
	    ok: boolean;
	    data: StartProxyResultDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_StartProxyResultDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], StartProxyResultDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	export class StorageStatsDTO {
	    request_rows: number;
	    inserted: number;
	    dropped: number;
	    queue_len: number;
	    proxy_running: boolean;
	
	    static createFrom(source: any = {}) {
	        return new StorageStatsDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.request_rows = source["request_rows"];
	        this.inserted = source["inserted"];
	        this.dropped = source["dropped"];
	        this.queue_len = source["queue_len"];
	        this.proxy_running = source["proxy_running"];
	    }
	}
	export class Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_StorageStatsDTO_ {
	    ok: boolean;
	    data: StorageStatsDTO;
	    error?: ErrorDTO;
	
	    static createFrom(source: any = {}) {
	        return new Result_github_com_kovacszso_ai_monitor_apps_proxy_desktop_internal_service_StorageStatsDTO_(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ok = source["ok"];
	        this.data = this.convertValues(source["data"], StorageStatsDTO);
	        this.error = this.convertValues(source["error"], ErrorDTO);
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
	
	
	
	
	
	export class UpdateProxyVariablesDTO {
	    server_host: string;
	    server_port: number;
	    server_max_request_bytes: number;
	    logging_level: string;
	    logging_format: string;
	    capture_store_content: boolean;
	    capture_max_prompt_bytes: number;
	    capture_max_response_bytes: number;
	
	    static createFrom(source: any = {}) {
	        return new UpdateProxyVariablesDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.server_host = source["server_host"];
	        this.server_port = source["server_port"];
	        this.server_max_request_bytes = source["server_max_request_bytes"];
	        this.logging_level = source["logging_level"];
	        this.logging_format = source["logging_format"];
	        this.capture_store_content = source["capture_store_content"];
	        this.capture_max_prompt_bytes = source["capture_max_prompt_bytes"];
	        this.capture_max_response_bytes = source["capture_max_response_bytes"];
	    }
	}

}

