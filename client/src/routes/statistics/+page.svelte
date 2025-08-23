<script lang="ts">
	import { onMount } from 'svelte';
	import CrimeFilters from './crime-filters.svelte';
	import Statcard from '$lib/components/statcard/stat-card.svelte';
	import TimeStatCard from '$lib/components/statcard/time-stat-card.svelte';

	let data = $state<CrimesResponse>();
	let isLoading = $state(false);
	let filter = $state<FilterState>({ year: '2025' });

	interface FilterState {
		year: string;
	}

	interface OrderedPair {
		key: string;
		value: number;
	}

	interface CrimesResponse {
		total_crimes: number;
		crimes_by_type: OrderedPair[];
		crimes_by_date: OrderedPair[];
		crimes_by_hour: OrderedPair[];
		most_dangerous_areas: string[];
		safest_areas: string[];
	}

	async function fetchCrimes(): Promise<void> {
		console.log('fetching the crime!');
		isLoading = true;
		try {
			const response = await fetch(`api/public/crimes/stats?year=${filter.year}`);
			if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

			data = await response.json();
			console.log(data);
		} catch (err) {
			console.log(err);
		} finally {
			isLoading = false;
		}
	}

	onMount(async () => {
		await fetchCrimes();
	});

	$effect(() => {
		fetchCrimes();
	});
</script>

<div class="min-h-screen space-y-6 p-6">
	{#if isLoading}
		<div class="flex items-center justify-center py-12">
			<div class="h-8 w-8 animate-spin rounded-full border-b-2 border-sky-600"></div>
			<span class="ml-3 text-lg text-gray-600">Engaging Granular networks...</span>
		</div>
	{:else}
		<div class="mb-4 flex items-center justify-between">
			<h2 class="text-xl font-bold">Crime Statistics</h2>
			<div class="">
				<CrimeFilters bind:filter />
			</div>
		</div>
		<div class="flex justify-center gap-3">
			<Statcard array={data?.crimes_by_type} title={'Top Crimes'} />
			<TimeStatCard array={data?.crimes_by_hour} title={'Crimes by Time'} />
			<Statcard array={data?.most_dangerous_areas} title={'Most Dangerous Neighborhoods'} />
			<Statcard array={data?.safest_areas} title={'Safest Neighborhoods'} />
		</div>
	{/if}
</div>

<style>
</style>
