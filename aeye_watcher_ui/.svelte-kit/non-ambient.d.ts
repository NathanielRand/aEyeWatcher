
// this file is generated — do not edit it


declare module "svelte/elements" {
	export interface HTMLAttributes<T> {
		'data-sveltekit-keepfocus'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-noscroll'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-preload-code'?:
			| true
			| ''
			| 'eager'
			| 'viewport'
			| 'hover'
			| 'tap'
			| 'off'
			| undefined
			| null;
		'data-sveltekit-preload-data'?: true | '' | 'hover' | 'tap' | 'off' | undefined | null;
		'data-sveltekit-reload'?: true | '' | 'off' | undefined | null;
		'data-sveltekit-replacestate'?: true | '' | 'off' | undefined | null;
	}
}

export {};


declare module "$app/types" {
	type MatcherParam<M> = M extends (param : string) => param is (infer U extends string) ? U : string;

	export interface AppTypes {
		RouteId(): "/" | "/about" | "/acp" | "/agents" | "/alerts" | "/config" | "/docs" | "/export" | "/gateways" | "/history" | "/mcp" | "/models" | "/plugins" | "/remote" | "/scan" | "/schedule" | "/services";
		RouteParams(): {
			
		};
		LayoutParams(): {
			"/": Record<string, never>;
			"/about": Record<string, never>;
			"/acp": Record<string, never>;
			"/agents": Record<string, never>;
			"/alerts": Record<string, never>;
			"/config": Record<string, never>;
			"/docs": Record<string, never>;
			"/export": Record<string, never>;
			"/gateways": Record<string, never>;
			"/history": Record<string, never>;
			"/mcp": Record<string, never>;
			"/models": Record<string, never>;
			"/plugins": Record<string, never>;
			"/remote": Record<string, never>;
			"/scan": Record<string, never>;
			"/schedule": Record<string, never>;
			"/services": Record<string, never>
		};
		Pathname(): "/" | "/about" | "/acp" | "/agents" | "/alerts" | "/config" | "/docs" | "/export" | "/gateways" | "/history" | "/mcp" | "/models" | "/plugins" | "/remote" | "/scan" | "/schedule" | "/services";
		ResolvedPathname(): `${"" | `/${string}`}${ReturnType<AppTypes['Pathname']>}`;
		Asset(): string & {};
	}
}