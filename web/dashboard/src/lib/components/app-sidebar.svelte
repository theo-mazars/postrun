<script lang="ts" module>
	import LayoutDashboard from "@lucide/svelte/icons/layout-dashboard";
	import Globe from "@lucide/svelte/icons/globe";
	import Key from "@lucide/svelte/icons/key";
	import FileText from "@lucide/svelte/icons/file-text";
	import Settings from "@lucide/svelte/icons/settings";

	const data = {
		user: {
			name: "shadcn",
			email: "m@example.com",
		},
		navMain: [
  		{
  			title: "Overview",
  			url: "/dashboard",
  			icon: LayoutDashboard,
  			isActive: true,
  		},
      {
				title: "Domains",
				url: "/dashboard/domains",
				icon: Globe,
  		},
  		{
        title: "API Keys",
        url: "/dashboard/keys",
        icon: Key,
   			},
      {
				title: "Logs",
				url: "/dashboard/emails",
				icon: FileText,
  		},
		],
		navSecondary: [
			{
				title: "Settings",
				url: "#",
				icon: Settings,
			},
		],
	};
</script>

<script lang="ts">
	import NavMain from "./nav-main.svelte";
	import NavSecondary from "./nav-secondary.svelte";
	import NavUser from "./nav-user.svelte";
	import * as Sidebar from "$lib/components/ui/sidebar/index.js";
	import CommandIcon from "@lucide/svelte/icons/command";
	import type { ComponentProps } from "svelte";

	let { ref = $bindable(null), ...restProps }: ComponentProps<typeof Sidebar.Root> = $props();
</script>

<Sidebar.Root bind:ref variant="inset" {...restProps}>
	<Sidebar.Header>
		<Sidebar.Menu>
			<Sidebar.MenuItem>
				<Sidebar.MenuButton size="lg">
					{#snippet child({ props })}
						<a href="##" {...props}>
							<div
								class="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg"
							>
								<CommandIcon class="size-4" />
							</div>
							<div class="grid flex-1 text-start text-sm leading-tight">
								<span class="truncate font-medium">Acme Inc</span>
								<span class="truncate text-xs">Enterprise</span>
							</div>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
		</Sidebar.Menu>
	</Sidebar.Header>
	<Sidebar.Content>
		<NavMain items={data.navMain} />
		<NavSecondary items={data.navSecondary} class="mt-auto" />
	</Sidebar.Content>
	<Sidebar.Footer>
		<NavUser user={data.user} />
	</Sidebar.Footer>
</Sidebar.Root>
