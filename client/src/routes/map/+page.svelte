<script lang="ts">
	import { onMount } from 'svelte';
	import { PUBLIC_GOOGLE_MAPS_API_KEY, PUBLIC_GOOGLE_MAPS_ID } from '$env/static/public';
	import type { Crime, CrimesResponse } from '$lib/types';

	const INITIAL_ZOOM = 12;

	let mElement: HTMLDivElement;
	let crimes: Crime[] = [];
	let totalCount: number = 0;
	let loading: boolean = false;
	let error: string | null = null;
	let map: any = null;
	let markers: any[] = [];
	let currentZoom: number = INITIAL_ZOOM;
	let updateTimeout: number | null = null;

	// Simple debounce
	function debounce(func: Function, delay: number) {
		return (...args: any[]) => {
			clearTimeout(updateTimeout!);
			updateTimeout = setTimeout(() => func(...args), delay) as any;
		};
	}

	async function fetchCrimes(): Promise<void> {
		loading = true;
		error = null;
		try {
			const response = await fetch('/api/public/crimes?year=2025&limit=500');
			if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
			
			const data: CrimesResponse = await response.json();
			crimes = data.crimes;
			totalCount = data.count;
			
			console.log(`Fetched ${crimes.length} crimes, total count: ${totalCount}`);
			
			if (map) updateMarkers();
		} catch (err) {
			error = err instanceof Error ? err.message : 'Unknown error';
		} finally {
			loading = false;
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

			// Single debounced update handler
			const debouncedUpdate = debounce(() => {
				const newZoom = map.getZoom();
				console.log(`Map event fired - Zoom: ${newZoom}`);
				currentZoom = newZoom;
				updateMarkers();
			}, 200);

			// Also add idle event as backup
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

	// Group crimes by exact location
	function groupByLocation(crimes: Crime[]): Map<string, Crime[]> {
		const groups = new Map<string, Crime[]>();
		crimes.forEach(crime => {
			const key = `${crime.latitude},${crime.longitude}`;
			if (!groups.has(key)) groups.set(key, []);
			groups.get(key)!.push(crime);
		});
		return groups;
	}

	// Get max markers based on zoom (more generous limits)
	function getMaxMarkers(zoom: number): number {
		if (zoom < 9) return 0;
		if (zoom < 11) return 50;
		if (zoom < 13) return 100;
		if (zoom < 15) return 200;
		return 300;
	}

	// Main update function
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

		// Filter crimes in viewport
		const visibleCrimes = crimes.filter(crime => {
			const pos = new google.maps.LatLng(crime.latitude, crime.longitude);
			return bounds.contains(pos);
		});

		console.log(`Found ${visibleCrimes.length} crimes in viewport out of ${crimes.length} total`);

		// Group by location and limit
		const locationGroups = groupByLocation(visibleCrimes);
		const maxMarkers = getMaxMarkers(currentZoom);
		const limitedGroups = Array.from(locationGroups.entries()).slice(0, maxMarkers);

		console.log(`Zoom: ${currentZoom}, Max markers: ${maxMarkers}, Showing: ${limitedGroups.length} groups`);

		// Clear and create new markers
		clearMarkers();
		createMarkers(limitedGroups);
	}

	// Create markers from location groups with smooth animation
	async function createMarkers(locationGroups: [string, Crime[]][]) {
		if (!map || locationGroups.length === 0) {
			console.log('No map or location groups to create markers');
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

			// Create markers with staggered animation
			locationGroups.forEach(([locationKey, crimesAtLocation], index) => {
				setTimeout(() => {
					const crime = crimesAtLocation[0];
					const count = crimesAtLocation.length;
					
					const markerElement = createMarkerElement(crimesAtLocation);
					
					// Start invisible
					markerElement.style.opacity = '0';
					markerElement.style.transform = 'scale(0.5)';
					
					const marker = new AdvancedMarkerElement({
						map: map,
						position: { lat: crime.latitude, lng: crime.longitude },
						title: count === 1 ? 
							`${crime.crime_type}\n${crime.address}\n${crime.date}` :
							`${count} crimes at ${crime.address}`,
						content: markerElement
					});

					marker.addListener('click', () => {
						if (count === 1) {
							alert(`${crime.crime_type}\n${crime.address}\n${crime.date}`);
						} else {
							const summary = crimesAtLocation.slice(0, 5)
								.map(c => `• ${c.crime_type} (${c.date})`)
								.join('\n') + (count > 5 ? `\n...and ${count - 5} more` : '');
							alert(`${count} crimes at ${crime.address}:\n\n${summary}`);
						}
					});

					markers.push(marker);

					// Animate in
					setTimeout(() => {
						markerElement.style.transition = 'opacity 0.3s ease, transform 0.3s ease';
						markerElement.style.opacity = '1';
						markerElement.style.transform = 'scale(1)';
					}, 50);

				}, index * 20); // 20ms delay between each marker
			});

			console.log(`Creating ${locationGroups.length} markers with smooth animation`);
		} catch (e) {
			console.error('Error creating markers:', e);
		}
	}

	// Create marker element
	function createMarkerElement(crimesAtLocation: Crime[]): HTMLElement {
		const count = crimesAtLocation.length;
		const crime = crimesAtLocation[0];

		if (count === 1) {
			// Single crime marker
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
			// Cluster marker
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
			'Robbery': '#cc0044'
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
		markers.forEach(marker => {
			if (marker.map) marker.map = null;
		});
		markers = [];
	}

	onMount(async () => {
		await Promise.all([fetchCrimes(), setupMap()]);
	});
</script>

<div class="p-4">
	{#if loading}
		<p class="mb-4">Loading crimes...</p>
	{:else if error}
		<p class="mb-4 text-red-500">Error: {error}</p>
	{:else}
		<div class="mb-4">
			<h2 class="text-xl font-bold">Tacoma Crime Map</h2>
			<p class="text-gray-600">{totalCount} total crimes</p>
			<div class="flex gap-4 text-sm text-gray-500">
				<span>Zoom: {currentZoom}</span>
				<span>Markers: {markers.length}</span>
			</div>
		</div>
	{/if}
	
	<div bind:this={mElement} class="h-170 w-full rounded bg-gray-200 shadow-lg"></div>
	
	{#if !loading && !error && totalCount > 0}
		<div class="mt-4 text-sm text-gray-600">
			<p>Markers are grouped by location and limited by zoom level for optimal performance.</p>
		</div>
	{/if}
</div>