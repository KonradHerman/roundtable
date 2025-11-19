<script lang="ts">
	import { cn } from '$lib/utils';

	let { class: className = undefined, children, ...restProps } = $props<{
		class?: string;
		children?: any;
		[key: string]: any;
	}>();
</script>

<div
	class={cn('rounded-xl border bg-card text-card-foreground shadow-lg', className)}
	{...restProps}
>
	{@render children?.()}
</div>
