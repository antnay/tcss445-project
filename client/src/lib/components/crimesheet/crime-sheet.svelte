<script lang="ts">
	let { visibility = $bindable(), close = $bindable(), crimes = $bindable() } = $props();

	function handleBackdropClick(event: MouseEvent) {
		if (event.target === event.currentTarget) {
			close();
		}
	}
</script>

<div
	class="fixed inset-0 z-40 transition-opacity duration-300"
	class:opacity-50={visibility}
	class:opacity-0={!visibility}
	class:pointer-events-none={!visibility}
	class:bg-black={visibility}
	onclick={handleBackdropClick}
	role="button"
	tabindex="0"
	onkeydown={(e) => e.key === 'Escape' && close()}
></div>

<div
	class="fixed right-0 top-0 z-50 h-full w-80 transform bg-white shadow-2xl transition-transform duration-300 ease-in-out"
	class:translate-x-0={visibility}
	class:translate-x-full={!visibility}
>
	<div class="flex items-center justify-between border-b border-gray-200 p-4">
		<h3 class="text-lg font-semibold text-gray-900">Crime Details</h3>
		<button
			onclick={close()}
			class="rounded-full p-1 transition-colors duration-200 hover:bg-gray-100"
			aria-label="Close crime sheet"
		>
			<svg class="h-5 w-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path
					stroke-linecap="round"
					stroke-linejoin="round"
					stroke-width="2"
					d="M6 18L18 6M6 6l12 12"
				></path>
			</svg>
		</button>
	</div>
	<div class="h-full overflow-y-auto p-4 pb-20">
		{#if crimes && crimes.length > 0}
			<div class="space-y-4">
				{#each crimes as crime, index}
					<div
						class="rounded-lg border border-gray-200 bg-gray-50 p-3 transition-colors duration-200 hover:bg-gray-100"
					>
						<p class="mb-1 font-medium text-gray-900">{crime.crime_type}</p>
						<p class="mb-1 flex items-center text-sm text-gray-600">
							<svg class="mr-1 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
								></path>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
								></path>
							</svg>
							{crime.address}
						</p>
						{#if crime.date}
							<p class="flex items-center text-sm text-gray-500">
								<svg class="mr-1 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
									></path>
								</svg>
								{crime.date}
							</p>
						{/if}
					</div>
				{/each}
			</div>
		{:else}
			<div class="py-8 text-center">
				<svg
					class="mx-auto mb-4 h-12 w-12 text-gray-400"
					fill="none"
					stroke="currentColor"
					viewBox="0 0 24 24"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
					></path>
				</svg>
				<p class="text-gray-500">No crime data available for this location</p>
			</div>
		{/if}
	</div>
</div>
