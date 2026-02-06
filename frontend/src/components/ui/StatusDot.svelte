<script lang="ts">
  import { cn } from '@/lib/utils';
  import type { Status } from '@/lib/status';

  export let status: Status;
  export let className = '';
  export let pulse = false;

  const dotColors: Record<Status, string> = {
    success: 'bg-status-success',
    error: 'bg-status-error',
    warning: 'bg-status-warning',
    running: 'bg-status-running',
    pending: 'bg-status-pending',
  };
</script>

<span class={cn('status-dot', dotColors[status], pulse && 'animate-pulse-subtle', className)}></span>
