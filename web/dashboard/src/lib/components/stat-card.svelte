<script lang="ts">
  import { cn } from "$lib/utils";
  import { ArrowUpRight, ArrowDownRight } from "@lucide/svelte/icons";
  import type { Component } from "svelte";

  interface Props {
    title: string;
    value: string;
    change?: string;
    changeType?: "positive" | "negative" | "neutral";
    icon: Component;
    class?: string;
  }

  let {
    title,
    value,
    change,
    changeType,
    icon: Icon,
    class: className,
  }: Props = $props();
</script>

<div class={cn("bg-card flex flex-col gap-6 shadow-sm rounded-xl p-5", className)}>
    <div class={cn("flex items-center justify-between")}>
        <div class={cn("text-muted-foreground")}>{title}</div>
        <Icon class={cn("text-muted-foreground size-5")} />
    </div>
    <div class="flex items-end justify-between">
        <span class="text-2xl font-semibold tracking-tight">{value}</span>
        {#if change}
          <div
            class={cn(
              "flex items-center gap-0.5 text-xs font-medium",
              changeType === "positive" && "text-success",
              changeType === "negative" && "text-destructive",
              changeType === "neutral" && "text-muted-foreground"
            )}
          >
            {#if changeType === "positive"}
              <ArrowUpRight class="h-3 w-3" />
            {:else if changeType === "negative"}
              <ArrowDownRight class="h-3 w-3" />
            {/if}
            <span>{change}</span>
          </div>
        {/if}
    </div>
</div>