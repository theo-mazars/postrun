<script lang="ts">
    import StatCard from '$lib/components/stat-card.svelte';
    import { Send, CircleCheckBig, CircleX, CircleAlert } from '@lucide/svelte/icons';
</script>

<div class="grid auto-rows-min gap-4 md:grid-cols-4">
    <StatCard title="Total Sent" value="24,847" change="+12.5%" changeType="positive" icon={Send} />
    <StatCard
        title="Delivered"
        value="23,982"
        change="96.5%"
        changeType="positive"
        icon={CircleCheckBig}
    />
    <StatCard title="Rejected" value="598" change="-2.3%" changeType="positive" icon={CircleX} />
    <StatCard title="Bounced" value="267" change="1.1%" changeType="neutral" icon={CircleAlert} />
</div>
<div class="grid auto-rows-min gap-4 md:grid-cols-5">
    <div class="col-span-3 h-96 w-full rounded-xl bg-muted/50"></div>
    <div class="col-span-2 h-96 w-full rounded-xl bg-muted/50"></div>
</div>
<div class="aspect-video rounded-xl bg-muted/50"></div>
