<script lang="ts">
	import DataTable from './data-table.svelte';
	import { columns } from './columns.js';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	let { data } = $props();
	let isLoading = $state(false);

	let searchForm = $state({
		latitude: data.initialSearch?.latitude?.toString() || '47.2529',
		longitude: data.initialSearch?.longitude?.toString() || '-122.4443',
		radius: data.initialSearch?.radius?.toString() || '1',
		year: $page.url.searchParams.get('year') || new Date().getFullYear().toString()
	});

	const availableYears = Array.from({ length: 10 }, (_, i) =>
		(new Date().getFullYear() - i).toString()
	);

	async function handleSearch() {
		if (!searchForm.latitude || !searchForm.longitude || !searchForm.radius) {
			alert('Please fill in latitude, longitude, and radius fields');
			return;
		}

		isLoading = true;

		const params = new URLSearchParams({
			lat: searchForm.latitude,
			lng: searchForm.longitude,
			radius: searchForm.radius,
			year: searchForm.year
		});

		await goto(`?${params.toString()}`);

		isLoading = false;
	}

	function clearForm() {
		searchForm = {
			latitude: '47.2529',
			longitude: '-122.4443',
			radius: '1',
			year: new Date().getFullYear().toString()
		};
	}
</script>

<div class="min-h-screen space-y-6 p-6">
	<div class="rounded-lgp-6 shadow-md">
		<h2 class="text-primary mb-6 text-2xl font-bold">Crime Radius Search</h2>

		<!-- svelte-ignore event_directive_deprecated -->
		<form on:submit|preventDefault={handleSearch} class="space-y-4">
			<div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
				<div>
					<label for="latitude" class="text-primary mb-1 block text-sm font-medium">
						Latitude
					</label>
					<input
						id="latitude"
						type="number"
						step="any"
						bind:value={searchForm.latitude}
						placeholder="47.2529"
						class="border-ring focus:ring-pink-30 w-full rounded-md border px-3 py-2 shadow-sm focus:outline-none focus:ring-1"
						required
					/>
				</div>

				<div>
					<label for="longitude" class="text-primary mb-1 block text-sm font-medium">
						Longitude
					</label>
					<input
						id="longitude"
						type="number"
						step="any"
						bind:value={searchForm.longitude}
						placeholder="-122.4443"
						class="border-ring focus:ring-pink-30 w-full rounded-md border px-3 py-2 shadow-sm focus:outline-none focus:ring-1"
						required
					/>
				</div>

				<div>
					<label for="radius" class="text-primary mb-1 block text-sm font-medium">
						Radius (miles)
					</label>
					<input
						id="radius"
						type="number"
						step="0.01"
						min="0.01"
						bind:value={searchForm.radius}
						placeholder="1"
						class="border-ring w-full rounded-md border px-3 py-2 shadow-sm focus:outline-none focus:ring-1 focus:ring-pink-300"
						required
					/>
				</div>

				<div>
					<label for="year" class="text-primary mb-1 block text-sm font-medium"> Year </label>
					<select
						id="year"
						bind:value={searchForm.year}
						class="border-ring w-full rounded-md border px-3 py-2 shadow-sm focus:outline-none focus:ring-1 focus:ring-pink-300"
					>
						{#each availableYears as year}
							<option value={year}>{year}</option>
						{/each}
					</select>
				</div>
			</div>

			<div class="flex flex-wrap gap-3 pt-4">
				<button
					type="submit"
					disabled={isLoading}
					class="text-primary rounded-md bg-rose-600 px-6 py-2 hover:bg-rose-700"
				>
					{isLoading ? 'Searching...' : 'Search Crimes'}
				</button>
				<button
					type="button"
					on:click={clearForm}
					class="text-primary rounded-md bg-gray-600 px-6 py-2 hover:bg-gray-700 focus:ring-gray-500"
				>
					Clear
				</button>
			</div>
		</form>
	</div>

	<div class="">
		{#if isLoading}
			<div class="flex items-center justify-center py-12">
				<div class="h-8 w-8 animate-spin rounded-full border-b-2 border-pink-600"></div>
				<span class="ml-3 text-lg text-gray-600">Loading Dynamic Datacenters...</span>
			</div>
		{:else if data.initialSearch}
			<DataTable data={data.initialSearch.result?.crimes || []} {columns} />
		{/if}
	</div>
</div>

<style>
</style>
