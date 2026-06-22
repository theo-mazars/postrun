<script lang="ts">
	import { ScrollArea } from '$lib/components/ui/scroll-area/index.js';
	import * as Item from '$lib/components/ui/item/index.js';
	import * as Card from '$lib/components/ui/card/index.js';
	import { Badge } from '$lib/components/ui/badge';
	import Label from '$lib/components/ui/label/label.svelte';
	import * as InputGroup from '$lib/components/ui/input-group/index.js';
	import {
		BadgeAlertIcon,
		BadgeCheckIcon,
		ChevronRightIcon,
		CopyIcon,
		PlusIcon,
		RefreshCwIcon,
		TrashIcon
	} from '@lucide/svelte';
	import Button from '$lib/components/ui/button/button.svelte';

	type DNSType = {
		type: 'TXT' | 'SPF' | 'DKIM' | 'DMARC';
		verified: boolean;
		name: string;
		value: string;
	};

	type DomainsType = {
		[name: string]: {
			sent: number;
			verified: boolean;
			dns: DNSType[];
		};
	};

	const domains: DomainsType = {
		'mail.acme.com': {
			sent: 15420,
			verified: true,
			dns: [
				{
					type: 'TXT',
					verified: true,
					name: 'mail.acme.com',
					value: 'v=spf1 include:postrun.dev ~all'
				},
				{
					type: 'TXT',
					verified: true,
					name: 'pr._domainkey.mail.acme.com',
					value: 'v=DKIM1; k=rsa; p=MIGfMA0GCS...'
				},
				{
					type: 'TXT',
					verified: true,
					name: '_dmarc.mail.acme.com',
					value: 'v=DMARC1; p=quarantine; rua=mailto:dmarc@acme.com'
				}
			]
		},
		'notifications.acme.com': {
			sent: 8234,
			verified: true,
			dns: [
				{
					type: 'TXT',
					verified: true,
					name: 'notifications.acme.com',
					value: 'v=spf1 include:postrun.dev ~all'
				},
				{
					type: 'TXT',
					verified: true,
					name: 'pr._domainkey.notifications.acme.com',
					value: 'v=DKIM1; k=rsa; p=MIGfMA0GCS...'
				},
				{
					type: 'TXT',
					verified: true,
					name: '_dmarc.notifications.acme.com',
					value: 'v=DMARC1; p=quarantine; rua=mailto:dmarc@acme.com'
				}
			]
		},
		'billing.acme.com': {
			sent: 0,
			verified: false,
			dns: [
				{
					type: 'TXT',
					verified: true,
					name: 'billing.acme.com',
					value: 'v=spf1 include:postrun.dev ~all'
				},
				{
					type: 'TXT',
					verified: true,
					name: 'pr._domainkey.billing.acme.com',
					value: 'v=DKIM1; k=rsa; p=MIGfMA0GCS...'
				},
				{
					type: 'TXT',
					verified: false,
					name: '_dmarc.billing.acme.com',
					value: 'v=DMARC1; p=quarantine; rua=mailto:dmarc@acme.com'
				}
			]
		}
	};
	let active: string = $state(Object.keys(domains)[0]);
</script>

<div class="grid grid-cols-12 gap-4">
	<ScrollArea class="col-span-3 flex flex-col gap-y-5 rounded-xl">
		<Button class="mb-2 w-full rounded-xl"><PlusIcon /> Add new domain</Button>
		{#each Object.entries(domains) as [name, domain] (name)}
			<Item.Root
				class={`mb-2 cursor-pointer rounded-xl ${domain.verified ? 'border-border' : 'border-destructive'}`}
				variant={name === active ? 'muted' : 'outline'}
			>
				{#snippet child({ props })}
					<a
						{...props}
						onclick={() => {
							active = name;
						}}
					>
						<Item.Media>
							{#if domain.verified}
								<BadgeCheckIcon class="size-4" />
							{:else}
								<BadgeAlertIcon class="size-4 text-destructive" />
							{/if}
						</Item.Media>
						<Item.Content>
							<Item.Title class="max-w-full break-all">{name}</Item.Title>
							<Item.Description
								>{new Intl.NumberFormat('fr-FR').format(domain.sent)} sent</Item.Description
							>
						</Item.Content>
						<Item.Actions>
							<ChevronRightIcon class="size-4" />
						</Item.Actions>
					</a>
				{/snippet}
			</Item.Root>{/each}
	</ScrollArea>
	<div class="col-span-9 max-w-3xl">
		<Card.Root class="rounded-xl p-4 shadow-none">
			<Card.Header class="flex flex-row items-center justify-between">
				<h3 class="scroll-m-20 text-xl font-semibold tracking-tight">{active}</h3>
				<div class="flex flex-row gap-4">
					<Button variant="outline"><RefreshCwIcon /> Verify</Button>
					<Button variant="destructive"><TrashIcon /> Delete</Button>
				</div>
			</Card.Header>
			<Card.Content>
				<p class="leading-7 not-first:mt-6">DNS Records</p>
				<p class="text-sm text-muted-foreground">Add these records to your DNS provider</p>
				{#each domains[active].dns as dns (dns.name)}
					<Card.Root class="my-4 bg-muted p-4">
						<Card.Header>
							<div>
								<Badge variant="outline">TXT</Badge>
								{#if dns.verified}
									<Badge>Verified</Badge>
								{:else}
									<Badge variant="destructive">Missing</Badge>
								{/if}
							</div>
						</Card.Header>
						<Card.Content>
							<div class="flex w-full flex-col gap-1.5">
								<Label class="uppercase" for="name-{dns.name}">Name</Label>
								<InputGroup.Root>
									<InputGroup.Input value={dns.name} id="name-{dns.name}" />
									<InputGroup.Addon align="inline-end">
										<InputGroup.Button>
											<CopyIcon />
										</InputGroup.Button>
									</InputGroup.Addon>
								</InputGroup.Root>
							</div>
							<div class="mt-4 flex w-full flex-col gap-1.5">
								<Label class="uppercase" for="value-{dns.name}">Value</Label>
								<InputGroup.Root>
									<InputGroup.Input value={dns.value} id="value-{dns.value}" />
									<InputGroup.Addon align="inline-end">
										<InputGroup.Button>
											<CopyIcon />
										</InputGroup.Button>
									</InputGroup.Addon>
								</InputGroup.Root>
							</div>
						</Card.Content>
					</Card.Root>
				{/each}
				<!-- <Card.Root class="my-4 bg-muted p-4">
					<Card.Header>
						<div>
							<Badge variant="outline">TXT</Badge>
							<Badge>Verified</Badge>
						</div>
					</Card.Header>
					<Card.Content>
						<div class="flex w-full flex-col gap-1.5">
							<Label class="uppercase" for="name-1">Name</Label>
							<InputGroup.Root>
								<InputGroup.Input value="mail.acme.com" id="name-1" />
								<InputGroup.Addon align="inline-end">
									<InputGroup.Button>
										<CopyIcon />
									</InputGroup.Button>
								</InputGroup.Addon>
							</InputGroup.Root>
						</div>
						<div class="mt-4 flex w-full flex-col gap-1.5">
							<Label class="uppercase" for="value-1">Value</Label>
							<InputGroup.Root>
								<InputGroup.Input value="v=spf1 include:postrun.dev ~all" id="value-1" />
								<InputGroup.Addon align="inline-end">
									<InputGroup.Button>
										<CopyIcon />
									</InputGroup.Button>
								</InputGroup.Addon>
							</InputGroup.Root>
						</div>
					</Card.Content>
				</Card.Root>
				<Card.Root class="my-4 bg-muted p-4">
					<Card.Header>
						<div>
							<Badge variant="outline">TXT</Badge>
							<Badge>Verified</Badge>
						</div>
					</Card.Header>
					<Card.Content>
						<div class="flex w-full flex-col gap-1.5">
							<Label class="uppercase" for="name-2">Name</Label>
							<InputGroup.Root>
								<InputGroup.Input value="zm._domainkey.mail.acme.com" id="name-2" />
								<InputGroup.Addon align="inline-end">
									<InputGroup.Button>
										<CopyIcon />
									</InputGroup.Button>
								</InputGroup.Addon>
							</InputGroup.Root>
						</div>
						<div class="mt-4 flex w-full flex-col gap-1.5">
							<Label class="uppercase" for="value-2">Value</Label>
							<InputGroup.Root>
								<InputGroup.Input value="v=DKIM1; k=rsa; p=MIGfMA0GCS..." id="value-2" />
								<InputGroup.Addon align="inline-end">
									<InputGroup.Button>
										<CopyIcon />
									</InputGroup.Button>
								</InputGroup.Addon>
							</InputGroup.Root>
						</div>
					</Card.Content>
				</Card.Root>
				<Card.Root class="my-4 bg-muted p-4">
					<Card.Header>
						<div>
							<Badge variant="outline">TXT</Badge>
							<Badge>Verified</Badge>
						</div>
					</Card.Header>
					<Card.Content>
						<div class="flex w-full flex-col gap-1.5">
							<Label class="uppercase" for="name-3">Name</Label>
							<InputGroup.Root>
								<InputGroup.Input value="_dmarc.mail.acme.com" id="name-3" />
								<InputGroup.Addon align="inline-end">
									<InputGroup.Button>
										<CopyIcon />
									</InputGroup.Button>
								</InputGroup.Addon>
							</InputGroup.Root>
						</div>
						<div class="mt-4 flex w-full flex-col gap-1.5">
							<Label class="uppercase" for="value-3">Value</Label>
							<InputGroup.Root>
								<InputGroup.Input
									value="v=DMARC1; p=quarantine; rua=mailto:dmarc@acme.com"
									id="value-3"
								/>
								<InputGroup.Addon align="inline-end">
									<InputGroup.Button>
										<CopyIcon />
									</InputGroup.Button>
								</InputGroup.Addon>
							</InputGroup.Root>
						</div>
					</Card.Content>
				</Card.Root> -->
			</Card.Content>
		</Card.Root>
	</div>
</div>
