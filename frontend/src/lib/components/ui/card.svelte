<script lang="ts">
	import { cn } from '$lib/utils';

	let className: string | undefined = undefined;
	export { className as class };
</script>

<div
	class={cn('rounded-xl border bg-card text-card-foreground shadow-lg', className)}
	{...$$restProps}
>
	<slot />
</div>
