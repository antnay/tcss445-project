<script lang="ts">
	import * as Card from '$lib/components/ui/card/index.js';

	interface OrderedPair {
		key: string;
		value: number;
	}

	let {
		title,
		array
	}: {
		title: string;
		array?: OrderedPair[] | string[];
	} = $props();
</script>

<Card.Root class="w-full min-w-[300px] max-w-xs ring-sky-300">
	<Card.Header>
		<Card.Title class="self-center">{title}</Card.Title>
	</Card.Header>
	<Card.Content class="self-center">
		{#if array}
			{#each array as item, count}
				{#if typeof item === 'string'}
					<p>{count + 1}. {item}</p>
				{:else}
					<p>{item.key} - {item.value} crimes</p>
				{/if}
			{/each}
		{/if}
	</Card.Content>
</Card.Root>
