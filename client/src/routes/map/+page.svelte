<script lang="ts">
	import { onMount } from 'svelte';
	import { PUBLIC_GOOGLE_MAPS_API_KEY, PUBLIC_GOOGLE_MAPS_ID } from '$env/static/public';
	import type { Crime, CrimesResponse } from '$lib/types';
	import Crimesheet from '$lib/components/crimesheet/crime-sheet.svelte';
	import CrimeFilters from './crime-filters.svelte';

	const INITIAL_ZOOM = 12;

	let mElement: HTMLDivElement;
	let crimes: Crime[] = [];
	let totalCount: number = 0;
	let error: string | null = null;
	let map: any = null;
	let markers: any[] = [];
	let currentZoom: number = INITIAL_ZOOM;
	let updateTimeout: number | null = null;
	let showCrimeSheet: boolean = $state(false);
	let selectedCrimes: any[] = $state([]);

	let data = $state<CrimesResponse>();
	let filter = $state<FilterState>({ year: '2025' });
	let isLoading = $state(false);

	interface FilterState {
		year: string;
	}

	function debounce(func: Function, delay: number) {
		return (...args: any[]) => {
			clearTimeout(updateTimeout!);
			updateTimeout = setTimeout(() => func(...args), delay) as any;
		};
	}

	async function fetchCrimes(): Promise<void> {
		console.log('fetching the crime!');
		isLoading = true;
		error = null;
		try {
			const response = await fetch(`/api/public/crimes?year=${filter.year}&limit=500`);
			if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);

			data = await response.json();
			crimes = data.crimes;
			totalCount = data.count;

			console.log(`Fetched ${crimes.length} crimes, total count: ${totalCount}`);

			if (map) updateMarkers();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unknown error';
		} finally {
			isLoading = false;
		}
	}

	async function setupMap() {
		try {
			const { Loader } = await import('@googlemaps/js-api-loader');
			const loader = new Loader({
				apiKey: PUBLIC_GOOGLE_MAPS_API_KEY,
				version: 'weekly',
				libraries: ['marker']
			});

			const { Map } = await loader.importLibrary('maps');
			map = new Map(mElement, {
				center: { lat: 47.2454, lng: -122.43854 },
				zoom: INITIAL_ZOOM,
				mapId: PUBLIC_GOOGLE_MAPS_ID
			});

			const debouncedUpdate = debounce(() => {
				const newZoom = map.getZoom();
				console.log(`Map event fired - Zoom: ${newZoom}`);
				currentZoom = newZoom;
				updateMarkers();
			}, 200);

			map.addListener('idle', () => {
				const newZoom = map.getZoom();
				if (newZoom !== currentZoom) {
					console.log(`Idle event - Zoom changed from ${currentZoom} to ${newZoom}`);
					currentZoom = newZoom;
					updateMarkers();
				}
			});

			map.addListener('bounds_changed', debouncedUpdate);
			map.addListener('zoom_changed', debouncedUpdate);

			if (crimes.length > 0) updateMarkers();
		} catch (e) {
			console.error('Error loading Google Maps:', e);
			error = 'Failed to load Google Maps';
		}
	}

	function groupByLocation(crimes: Crime[]): Map<string, Crime[]> {
		const groups = new Map<string, Crime[]>();
		crimes.forEach((crime) => {
			const key = `${crime.latitude},${crime.longitude}`;
			if (!groups.has(key)) groups.set(key, []);
			groups.get(key)!.push(crime);
		});
		return groups;
	}

	function getMaxMarkers(zoom: number): number {
		if (zoom < 9) return 0;
		if (zoom < 11) return 50;
		if (zoom < 13) return 100;
		if (zoom < 15) return 200;
		return 300;
	}

	function updateMarkers() {
		if (!map || crimes.length === 0) {
			console.log('No map or crimes available');
			return;
		}

		const bounds = map.getBounds();
		if (!bounds) {
			console.log('No map bounds available');
			return;
		}

		const visibleCrimes = crimes.filter((crime) => {
			const pos = new google.maps.LatLng(crime.latitude, crime.longitude);
			return bounds.contains(pos);
		});

		console.log(`Found ${visibleCrimes.length} crimes in viewport out of ${crimes.length} total`);

		const locationGroups = groupByLocation(visibleCrimes);
		const maxMarkers = getMaxMarkers(currentZoom);
		const limitedGroups = Array.from(locationGroups.entries()).slice(0, maxMarkers);

		clearMarkers();
		createMarkers(limitedGroups);
	}

	async function createMarkers(locationGroups: [string, Crime[]][]) {
		if (!map) {
			console.log('No map to create markers');
			return;
		} else if (locationGroups.length === 0) {
			console.log('No location groups to create markers');
			return;
		}

		try {
			const { Loader } = await import('@googlemaps/js-api-loader');
			const loader = new Loader({
				apiKey: PUBLIC_GOOGLE_MAPS_API_KEY,
				version: 'weekly',
				libraries: ['marker']
			});
			const { AdvancedMarkerElement } = await loader.importLibrary('marker');

			locationGroups.forEach(([locationKey, crimesAtLocation], index) => {
				setTimeout(() => {
					const crime = crimesAtLocation[0];
					const count = crimesAtLocation.length;

					const markerElement = createMarkerElement(crimesAtLocation);

					markerElement.style.opacity = '0';
					markerElement.style.transform = 'scale(0.5)';

					const marker = new AdvancedMarkerElement({
						map: map,
						position: { lat: crime.latitude, lng: crime.longitude },
						title:
							count === 1
								? `${crime.crime_type}\n${crime.address}\n${crime.date}`
								: `${count} crimes at ${crime.address}`,
						content: markerElement
					});

					marker.addListener('click', () => {
						selectedCrimes = crimesAtLocation;
						showCrimeSheet = true;
					});

					markers.push(marker);

					setTimeout(() => {
						markerElement.style.transition = 'opacity 0.3s ease, transform 0.3s ease';
						markerElement.style.opacity = '1';
						markerElement.style.transform = 'scale(1)';
					}, 50);
				}, index * 20);
			});

			console.log(`Creating ${locationGroups.length} markers with smooth animation`);
		} catch (e) {
			console.error('Error creating markers:', e);
		}
	}

	function createMarkerElement(crimesAtLocation: Crime[]): HTMLElement {
		const count = crimesAtLocation.length;
		const crime = crimesAtLocation[0];

		if (count === 1) {
			const div = document.createElement('div');
			div.style.cssText = `
				width: 12px;
				height: 12px;
				background-color: ${getCrimeColor(crime.crime_type)};
				border: 2px solid white;
				border-radius: 50%;
				cursor: pointer;
				box-shadow: 0 1px 3px rgba(0,0,0,0.3);
			`;
			return div;
		} else {
			const container = document.createElement('div');
			const size = Math.min(16 + count, 28);

			container.style.cssText = `
				width: ${size}px;
				height: ${size}px;
				background-color: ${getClusterColor(count)};
				border: 2px solid white;
				border-radius: 50%;
				display: flex;
				align-items: center;
				justify-content: center;
				cursor: pointer;
				box-shadow: 0 2px 4px rgba(0,0,0,0.3);
				color: white;
				font-size: ${Math.max(10, size * 0.4)}px;
				font-weight: bold;
			`;

			container.textContent = count.toString();
			return container;
		}
	}

	function getCrimeColor(crimeType: string): string {
		const colors: Record<string, string> = {
			'Assault Offenses': '#ff4444',
			'Motor Vehicle Theft': '#ff8800',
			'Burglary/Breaking & Entering': '#cc0000',
			'Larceny/Theft Offenses': '#ffaa00',
			'Destruction/Damage/Vandalism': '#4488ff',
			'Traffic - DUI (Liquor)': '#8800cc',
			'Fraud Offenses': '#00cc88',
			'Drug/Narcotics Violations': '#cc8800',
			Robbery: '#cc0044'
		};
		return colors[crimeType] || '#777777';
	}

	function getClusterColor(count: number): string {
		if (count < 5) return '#34a853';
		if (count < 10) return '#fbbc04';
		if (count < 20) return '#ea4335';
		return '#9c27b0';
	}

	function clearMarkers() {
		markers.forEach((marker) => {
			if (marker.map) marker.map = null;
		});
		markers = [];
	}

	function closeCrimeSheet() {
		showCrimeSheet = false;
	}

	onMount(async () => {
		await Promise.all([fetchCrimes(), setupMap()]);
	});

	$effect(() => {
		fetchCrimes();
	});
</script>

<div class="p-4">
	<div class="mb-4 flex items-center justify-between">
		<h2 class="text-xl font-bold">Crime Map</h2>
		<div class="">
			<CrimeFilters bind:filter />
		</div>
	</div>
	<div bind:this={mElement} class="h-170 w-full rounded bg-gray-200 shadow-lg"></div>
	<Crimesheet crimes={selectedCrimes} visibility={showCrimeSheet} close={closeCrimeSheet} />
</div>
