<script lang="ts">
	import NavigationMenuItem from '../ui/navigation-menu/navigation-menu-item.svelte';
	import NavigationMenuLink from '../ui/navigation-menu/navigation-menu-link.svelte';
	import NavigationMenuList from '../ui/navigation-menu/navigation-menu-list.svelte';
	import NavigationMenu from '../ui/navigation-menu/navigation-menu.svelte';
</script>

<div class="p-4">
	<NavigationMenu>
		<NavigationMenuList>
			<NavigationMenuItem>
				<NavigationMenuLink href="/">home</NavigationMenuLink>
			</NavigationMenuItem>
			<NavigationMenuItem>
				<NavigationMenuLink href="/map">map</NavigationMenuLink>
			</NavigationMenuItem>
			<NavigationMenuItem>
				<NavigationMenuLink href="/statistics">statistics</NavigationMenuLink>
			</NavigationMenuItem>
			<NavigationMenuItem>
				<NavigationMenuLink href="/crimes">crimes</NavigationMenuLink>
			</NavigationMenuItem>
		</NavigationMenuList>
	</NavigationMenu>
</div>
