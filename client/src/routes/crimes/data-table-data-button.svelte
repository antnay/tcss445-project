<script lang="ts">
	import type { ComponentProps } from 'svelte';
	import ArrowUpDownIcon from '@lucide/svelte/icons/arrow-up-down';
	import { Button } from '$lib/components/ui/button/index.js';

	type Props = ComponentProps<typeof Button> & {
		text: string;
	};

	let { text, variant = 'ghost', ...restProps }: Props = $props();
</script>

<Button {variant} {...restProps}>
	{text}
	<ArrowUpDownIcon class="ml-2" />
</Button>
