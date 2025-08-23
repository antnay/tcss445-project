<script lang="ts">
	import * as Select from '$lib/components/ui/select/index.js';
	import { onMount } from 'svelte';

	let { filter = $bindable() } = $props();

	let availableYears = $state<string[]>([]);
	let isLoadingOptions = $state(false);
	let isInitialized = $state(false);

	async function loadFilterOptions() {
		isLoadingOptions = true;
		try {
			const response = await fetch('/api/public/options');
			if (response.ok) {
				const data = await response.json();
				availableYears = data.years || [];
			} else {
				throw new Error('Failed to fetch filter options');
			}
		} catch (error) {
			console.error('Error loading filter options:', error);
			availableYears = [];
		} finally {
			isLoadingOptions = false;
		}
	}

	function handleYearChange(value: string) {
		if (value !== filter.year) {
			filter = { ...filter, year: value };
		}
	}

	onMount(async () => {
		console.log('CrimeFilters component mounted');
		try {
			await loadFilterOptions();
			console.log('Filter options loaded:', {
				years: availableYears.length
			});
			isInitialized = true;
		} catch (error) {
			console.error('Error in onMount:', error);
		}
	});

	$effect(() => {});
</script>

<div class="flex flex-col">
	<div class="flex items-center gap-2">
		<Select.Root
			type="single"
			bind:value={filter.year as string}
			onValueChange={handleYearChange}
			name={'Year'}
		>
			<Select.Trigger>
				{filter.year ? filter.year : 'Select year'}
			</Select.Trigger>
			<Select.Content>
				{#each availableYears as year (year)}
					<Select.Item value={year} label={year} />
				{/each}
			</Select.Content>
		</Select.Root>
	</div>
</div>
