<script lang="ts">
	import * as Popover from '$lib/components/ui/popover/index.js';
	import * as Select from '$lib/components/ui/select/index.js';
	import { buttonVariants } from '$lib/components/ui/button/index.js';
	import { Input } from '$lib/components/ui/input/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import { Checkbox } from '$lib/components/ui/checkbox/index.js';
	import { FilterIcon, X } from '@lucide/svelte';
	import { onMount } from 'svelte';

	interface Props {
		filtersChanged: (filters: FilterState) => void;
	}
	let { filtersChanged }: Props = $props();

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

	let availableYears = $state<string[]>([]);
	let availableCrimeTypes = $state<string[]>([]);
	let availableCities = $state<string[]>([]);
	let availableNeighborhoods = $state<string[]>([]);
	let availableSources = $state<string[]>([]);
	let isLoadingOptions = $state(false);

	const limitOptions = [100, 500, 1000, 2000, 5000];

	let filters = $state<FilterState>({
		years: [],
		crimeTypes: [],
		cities: [],
		neighborhoods: [],
		sources: [],
		dateRange: {
			start: '',
			end: ''
		},
		specificDate: '',
		caseNumber: '',
		zipCode: '',
		street: '',
		limit: 1000
	});

	async function loadFilterOptions() {
		isLoadingOptions = true;
		try {
			const response = await fetch('/api/public/options');

			if (response.ok) {
				const data = await response.json();
				availableYears = data.years || [];
				availableCrimeTypes = data.crimeTypes || [];
				availableCities = data.cities || [];
				availableNeighborhoods = data.neighborhoods || [];
				availableSources = data.sources || [];
			} else {
				throw new Error('Failed to fetch filter options');
			}
		} catch (error) {
			console.error('Error loading filter options:', error);
			availableYears = [];
			availableCrimeTypes = [];
			availableCities = [];
			availableNeighborhoods = [];
			availableSources = [];
		} finally {
			isLoadingOptions = false;
		}
	}

	function toggleSelectAllYears(checked: boolean) {
		if (checked) {
			filters.years = [...availableYears];
		} else {
			filters.years = [];
		}
	}

	function applyFilters() {
		console.log('Applying filters:', filters);
		filtersChanged(filters);
	}

	function clearFilters() {
		filters = {
			years: [],
			crimeTypes: [],
			cities: [],
			neighborhoods: [],
			sources: [],
			dateRange: { start: '', end: '' },
			specificDate: '',
			caseNumber: '',
			zipCode: '',
			street: '',
			limit: 1000
		};
		applyFilters();
	}

	onMount(async () => {
		console.log('CrimeFilters component mounted');
		try {
			await loadFilterOptions();
			console.log('Filter options loaded:', {
				years: availableYears.length,
				crimeTypes: availableCrimeTypes.length,
				cities: availableCities.length,
				neighborhoods: availableNeighborhoods.length,
				sources: availableSources.length
			});

			const currentYear = new Date().getFullYear().toString();
			if (availableYears.includes(currentYear)) {
				filters.years = [currentYear];
				console.log('Set default year:', currentYear);
				applyFilters();
			}
		} catch (error) {
			console.error('Error in onMount:', error);
		}
	});

	$effect(() => {});
</script>

<div class="flex flex-col">
	<div class="flex items-center gap-2">
		<Popover.Root>
			<Popover.Trigger class={buttonVariants({ variant: 'secondary' })}>
				<FilterIcon class="mr-2 h-4 w-4" />
				Filters
			</Popover.Trigger>
			<Popover.Content class="max-h-120 w-108 ml-4 overflow-y-auto pl-4">
				{#if isLoadingOptions}
					<div class="flex items-center justify-center py-8">
						<div class="h-6 w-6 animate-spin rounded-full border-b-2 border-gray-900"></div>
						<span class="ml-2 text-sm">Loading filter options...</span>
					</div>
				{:else}
					<div class="space-y-4">
						<div class="space-y-2">
							<div class="flex items-center justify-between">
								<span class="text-sm font-medium">Years</span>
								<div class="flex items-center space-x-2">
									<Checkbox
										id="select-all-years"
										checked={filters.years.length === availableYears.length &&
											availableYears.length > 0}
										onCheckedChange={toggleSelectAllYears}
									/>
									<label for="select-all-years" class="text-xs text-gray-600">Select All</label>
								</div>
							</div>
							<div class="grid grid-cols-3 gap-2">
								{#each availableYears as year (year)}
									<div class="flex items-center space-x-2">
										<Checkbox
											id="year-{year}"
											checked={filters.years.includes(year)}
											onCheckedChange={(checked) => {
												filters.years = checked
													? [...filters.years, year]
													: filters.years.filter((y) => y !== year);
											}}
										/>
										<label for="year-{year}" class="text-sm">{year}</label>
									</div>
								{/each}
							</div>
						</div>

						<div class="space-y-2">
							<span class="text-sm font-medium">Crime Categories</span>
							<div class="max-h-40 overflow-y-auto rounded-md border p-2">
								<div class="grid grid-cols-1 gap-2">
									{#each availableCrimeTypes as crimeType (crimeType)}
										<div class="flex items-center space-x-2">
											<Checkbox
												id="crime-{crimeType.replace(/\s+/g, '-').toLowerCase()}"
												checked={filters.crimeTypes.includes(crimeType)}
												onCheckedChange={(checked) => {
													filters.crimeTypes = checked
														? [...filters.crimeTypes, crimeType]
														: filters.crimeTypes.filter((c) => c !== crimeType);
												}}
											/>
											<label
												for="crime-{crimeType.replace(/\s+/g, '-').toLowerCase()}"
												class="cursor-pointer text-xs leading-tight hover:text-blue-600"
												title={crimeType}
											>
												{crimeType}
											</label>
										</div>
									{/each}
								</div>
							</div>
							{#if filters.crimeTypes.length > 0}
								<div class="text-xs text-gray-500">
									{filters.crimeTypes.length} crime type{filters.crimeTypes.length !== 1 ? 's' : ''}
									selected
								</div>
							{/if}
						</div>

						<div class="space-y-2">
							<span class="text-sm font-medium">Date Range</span>
							<div class="grid grid-cols-2 gap-2">
								<div>
									<label for="start-date" class="text-xs text-gray-500">Start Date</label>
									<Input
										id="start-date"
										type="date"
										bind:value={filters.dateRange.start}
										class="w-full"
									/>
								</div>
								<div>
									<label for="end-date" class="text-xs text-gray-500">End Date</label>
									<Input
										id="end-date"
										type="date"
										bind:value={filters.dateRange.end}
										class="w-full"
									/>
								</div>
							</div>
						</div>

						<div class="space-y-2">
							<label for="specific-date" class="text-sm font-medium">Specific Date</label>
							<Input
								id="specific-date"
								type="date"
								bind:value={filters.specificDate}
								placeholder="Select specific date"
								class="w-full"
							/>
						</div>

						<!-- Location Filters -->
						<div class="space-y-2">
							<span class="text-sm font-medium">Cities</span>
							{#if availableCities.length === 1}
								<div class="rounded bg-gray-50 p-2 text-sm text-gray-600">
									Showing data for: <strong>{availableCities[0]}</strong>
								</div>
							{:else}
								<div class="grid grid-cols-2 gap-2">
									{#each availableCities as city (city)}
										<div class="flex items-center space-x-2">
											<Checkbox
												id="city-{city.replace(/\s+/g, '-').toLowerCase()}"
												checked={filters.cities.includes(city)}
												onCheckedChange={(checked) => {
													filters.cities = checked
														? [...filters.cities, city]
														: filters.cities.filter((c) => c !== city);
												}}
											/>
											<label for="city-{city.replace(/\s+/g, '-').toLowerCase()}" class="text-sm"
												>{city}</label
											>
										</div>
									{/each}
								</div>
							{/if}
						</div>

						<div class="space-y-2">
							<span class="text-sm font-medium"
								>Neighborhoods ({availableNeighborhoods.length} available)</span
							>
							<div class="max-h-32 overflow-y-auto rounded-md border p-2">
								<div class="grid grid-cols-2 gap-2">
									{#each availableNeighborhoods as neighborhood (neighborhood)}
										<div class="flex items-center space-x-2">
											<Checkbox
												id="neighborhood-{neighborhood.replace(/\s+/g, '-').toLowerCase()}"
												checked={filters.neighborhoods.includes(neighborhood)}
												onCheckedChange={(checked) => {
													filters.neighborhoods = checked
														? [...filters.neighborhoods, neighborhood]
														: filters.neighborhoods.filter((n) => n !== neighborhood);
												}}
											/>
											<label
												for="neighborhood-{neighborhood.replace(/\s+/g, '-').toLowerCase()}"
												class="text-sm">{neighborhood}</label
											>
										</div>
									{/each}
								</div>
							</div>
							{#if filters.neighborhoods.length > 0}
								<div class="text-xs text-gray-500">
									{filters.neighborhoods.length} neighborhood{filters.neighborhoods.length !== 1
										? 's'
										: ''} selected
								</div>
							{/if}
						</div>

						<div class="space-y-2">
							<label for="case-number" class="text-sm font-medium">Case Number</label>
							<Input
								id="case-number"
								bind:value={filters.caseNumber}
								placeholder="Enter case number..."
							/>
						</div>

						<div class="space-y-2">
							<label for="street-address" class="text-sm font-medium">Street Address</label>
							<Input
								id="street-address"
								bind:value={filters.street}
								placeholder="Enter street name..."
							/>
						</div>

						<div class="space-y-2">
							<label for="zip-code" class="text-sm font-medium">Zip Code</label>
							<Input id="zip-code" bind:value={filters.zipCode} placeholder="Enter zip code..." />
						</div>

						<div class="space-y-2">
							<span class="text-sm font-medium"
								>Data Sources ({availableSources.length} available)</span
							>
							<div class="max-h-32 overflow-y-auto rounded-md border p-2">
								<div class="grid grid-cols-1 gap-2">
									{#each availableSources as source (source)}
										<div class="flex items-center space-x-2">
											<Checkbox
												id="source-{source.replace(/\s+/g, '-').toLowerCase()}"
												checked={filters.sources.includes(source)}
												onCheckedChange={(checked) => {
													filters.sources = checked
														? [...filters.sources, source]
														: filters.sources.filter((s) => s !== source);
												}}
											/>
											<label
												for="source-{source.replace(/\s+/g, '-').toLowerCase()}"
												class="cursor-pointer text-xs leading-tight hover:text-blue-600"
												title={source}
											>
												{source}
											</label>
										</div>
									{/each}
								</div>
							</div>
							{#if filters.sources.length > 0}
								<div class="text-xs text-gray-500">
									{filters.sources.length} source{filters.sources.length !== 1 ? 's' : ''} selected
								</div>
							{/if}
						</div>

						<div class="space-y-2">
							<span class="text-sm font-medium">Result Limit</span>
							<Select.Root
								type="single"
								bind:value={filters.limit as unknown as string}
								name={'Limit'}
							>
								<Select.Trigger>
									{filters.limit ? filters.limit : 'Select limit'}
								</Select.Trigger>
								<Select.Content>
									{#each limitOptions as limit (limit)}
										<Select.Item
											value={limit as unknown as string}
											label={limit as unknown as string}
										/>
									{/each}
								</Select.Content>
							</Select.Root>
						</div>

						<Button onclick={applyFilters} class="w-full">Apply Filters</Button>
					</div>
				{/if}
			</Popover.Content>
		</Popover.Root>
	</div>
</div>
