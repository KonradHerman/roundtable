<script lang="ts">
	import { cn } from '$lib/utils';

	type Variant = 'default' | 'secondary' | 'destructive' | 'outline';

	let { 
		class: className = undefined, 
		variant = 'default' as Variant,
		children,
		...restProps
	} = $props<{
		class?: string;
		variant?: Variant;
		children?: any;
		[key: string]: any;
	}>();

	const variants: Record<Variant, string> = {
		default: 'bg-primary text-primary-foreground hover:bg-primary/80',
		secondary: 'bg-secondary text-secondary-foreground hover:bg-secondary/80',
		destructive: 'bg-destructive text-destructive-foreground hover:bg-destructive/80',
		outline: 'text-foreground border border-input'
	};
</script>

<div
	class={cn(
		'inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold transition-colors focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2',
		variants[variant as Variant],
		className
	)}
	{...restProps}
>
	{@render children?.()}
</div>
