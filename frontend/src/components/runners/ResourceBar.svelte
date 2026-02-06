<script lang="ts">
  import { cn } from '@/lib/utils';

  export let value: number;
  export let label: string;

  const getColor = (v: number) => {
    if (v > 80) return 'bg-status-error';
    if (v > 60) return 'bg-status-warning';
    return 'bg-status-success';
  };
</script>

<div class="flex items-center gap-2 min-w-[100px]">
  <span class="text-[10px] text-muted-foreground uppercase w-8">{label}</span>
  <div class="flex-1 h-1.5 rounded-full bg-muted/30 overflow-hidden">
    <div
      class={cn('h-full rounded-full transition-all', getColor(value))}
      style={`width: ${value}%`}
    ></div>
  </div>
  <span class="text-xs font-mono text-muted-foreground w-8 text-right">{value}%</span>
</div>
