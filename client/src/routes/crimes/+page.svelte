<script lang="ts">
	import DataTable from './data-table.svelte';
	import { columns } from './columns.js';
	import CrimeFilters from './crime-filters.svelte';
	import type { CrimeData } from './columns.js';

	let { data } = $props();

	let filteredData = $state<CrimeData[]>(data.crimes);
	let isLoading = $state(false);

	interface FilterState {
		years: string[];
		crimeTypes: string[];
		cities: string[];
		neighborhoods: string[];
		sources: string[];
		dateRange: {
			start: string;
			end: string;
		};
		specificDate: string;
		caseNumber: string;
		zipCode: string;
		street: string;
		limit: number;
	}

	async function handleFiltersChanged(filters: FilterState) {
		isLoading = true;

		try {
			const params = new URLSearchParams();

			if (filters.years && filters.years.length > 0) {
				filters.years.forEach((year: string) => params.append('year', year));
			}
			if (filters.crimeTypes && filters.crimeTypes.length > 0) {
				params.append('crimeType', filters.crimeTypes.join(','));
			}
			if (filters.cities && filters.cities.length > 0) {
				params.append('city', filters.cities.join(','));
			}
			if (filters.neighborhoods && filters.neighborhoods.length > 0) {
				params.append('neighborhood', filters.neighborhoods.join(','));
			}
			if (filters.dateRange?.start) {
				params.append('startDate', filters.dateRange.start);
			}
			if (filters.dateRange?.end) {
				params.append('endDate', filters.dateRange.end);
			}
			if (filters.specificDate) {
				params.append('date', filters.specificDate);
			}
			if (filters.caseNumber) {
				params.append('caseNumber', filters.caseNumber);
			}
			if (filters.street) {
				params.append('street', filters.street);
			}
			if (filters.zipCode) {
				params.append('zipCode', filters.zipCode);
			}

			if (filters.sources && filters.sources.length > 0) {
				params.append('source', filters.sources.join(','));
			}

			if (filters.limit) {
				params.append('limit', filters.limit.toString());
			}

			const response = await fetch(`/api/public/crimes/details?${params.toString()}`);
			if (!response.ok) {
				throw new Error('Failed to fetch filtered data');
			}

			const result = await response.json();
			filteredData = result.crimes || [];
		} catch (error) {
			console.error('Error fetching filtered data:', error);
			// You might want to show a toast or error message here
		} finally {
			isLoading = false;
		}
	}
</script>

<div class="h-180 space-y-4 p-10">
	<CrimeFilters filtersChanged={handleFiltersChanged} />

	{#if isLoading}
		<div class="flex items-center justify-center py-12">
			<div class="h-8 w-8 animate-spin rounded-full border-b-2 border-pink-600"></div>
			<span class="ml-3 text-lg text-gray-600">Screaming at children...</span>
		</div>
	{:else}
		<DataTable data={filteredData} {columns} />
	{/if}
</div>
