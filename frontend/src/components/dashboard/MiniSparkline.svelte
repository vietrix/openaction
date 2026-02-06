<script lang="ts">
  import { TrendingUp, TrendingDown, Minus } from '@lucide/svelte';
  import { cn } from '@/lib/utils';

  export let data: number[] = [];
  export let color = 'hsl(240, 5%, 65%)';
  export let label = '';
  export let suffix = '';
  export let showTrend = true;

  $: max = Math.max(...data);
  $: min = Math.min(...data);
  $: range = max - min || 1;
  $: current = data[data.length - 1] ?? 0;
  $: previous = data[data.length - 2] ?? current;
  $: trend = current - previous;

  $: points = data
    .map((value, i) => {
      const x = (i / (data.length - 1)) * 60;
      const y = 20 - ((value - min) / range) * 16;
      return `${x},${y}`;
    })
    .join(' ');
</script>

<div class="flex items-center gap-3 px-4 py-2 rounded-md border border-border bg-card/50">
  <div class="flex flex-col gap-0.5">
    <span class="text-[10px] uppercase tracking-wide text-muted-foreground">{label}</span>
    <div class="flex items-center gap-1.5">
      <span class="text-sm font-medium text-foreground font-mono">
        {current.toFixed(1)}{suffix}
      </span>
      {#if showTrend}
        <span
          class={cn(
            'flex items-center gap-0.5 text-[10px]',
            trend > 0
              ? 'text-status-error'
              : trend < 0
                ? 'text-status-success'
                : 'text-muted-foreground'
          )}
        >
          {#if trend > 0}
            <TrendingUp class="h-3 w-3" />
          {:else if trend < 0}
            <TrendingDown class="h-3 w-3" />
          {:else}
            <Minus class="h-3 w-3" />
          {/if}
        </span>
      {/if}
    </div>
  </div>
  <svg width="60" height="24" class="opacity-60">
    <polyline
      points={points}
      fill="none"
      stroke={color}
      stroke-width="1.5"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
  </svg>
</div>
