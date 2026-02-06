<script lang="ts">
  import { Activity } from '@lucide/svelte';
  import { cn } from '@/lib/utils';
  import MiniSparkline from './MiniSparkline.svelte';

  export let stats: { success: number; failed: number; running: number };
  export let durationTrend: number[] = [];
  export let failureRateTrend: number[] = [];
  export let className = '';

  $: total = stats.success + stats.failed + stats.running;
  $: successPercent = total > 0 ? (stats.success / total) * 100 : 0;
  $: failedPercent = total > 0 ? (stats.failed / total) * 100 : 0;
  $: runningPercent = total > 0 ? (stats.running / total) * 100 : 0;
</script>

<div class={cn('rounded-lg border border-border bg-card p-4', className)}>
  <div class="flex items-center gap-2 mb-3">
    <Activity class="h-4 w-4 text-muted-foreground" />
    <h3 class="text-sm font-medium text-foreground">System Snapshot</h3>
    <span class="text-xs text-muted-foreground ml-auto">Today</span>
  </div>

  <div class="flex flex-col lg:flex-row gap-4">
    <div class="flex-1">
      <div class="flex items-center justify-between mb-2">
        <span class="text-xs text-muted-foreground">Pipeline Health</span>
        <span class="text-xs font-mono text-muted-foreground">{total} total</span>
      </div>
      <div class="h-3 rounded-full bg-muted/30 overflow-hidden flex">
        {#if successPercent > 0}
          <div
            class="h-full bg-status-success transition-all duration-500"
            style={`width: ${successPercent}%`}
          ></div>
        {/if}
        {#if failedPercent > 0}
          <div
            class="h-full bg-status-error transition-all duration-500"
            style={`width: ${failedPercent}%`}
          ></div>
        {/if}
        {#if runningPercent > 0}
          <div
            class="h-full bg-status-running animate-pulse transition-all duration-500"
            style={`width: ${runningPercent}%`}
          ></div>
        {/if}
      </div>
      <div class="flex items-center gap-4 mt-2">
        <span class="flex items-center gap-1.5 text-xs">
          <span class="h-2 w-2 rounded-full bg-status-success"></span>
          <span class="text-muted-foreground">{stats.success} passed</span>
        </span>
        <span class="flex items-center gap-1.5 text-xs">
          <span class="h-2 w-2 rounded-full bg-status-error"></span>
          <span class="text-muted-foreground">{stats.failed} failed</span>
        </span>
        <span class="flex items-center gap-1.5 text-xs">
          <span class="h-2 w-2 rounded-full bg-status-running animate-pulse"></span>
          <span class="text-muted-foreground">{stats.running} running</span>
        </span>
      </div>
    </div>

    <div class="flex gap-2 shrink-0">
      <MiniSparkline data={durationTrend} color="hsl(240, 5%, 65%)" label="Avg Duration" suffix="m" />
      <MiniSparkline data={failureRateTrend} color="hsl(0, 65%, 55%)" label="Failure Rate" suffix="%" />
    </div>
  </div>
</div>
