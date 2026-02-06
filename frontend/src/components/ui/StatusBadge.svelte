<script lang="ts">
  import { cn } from '@/lib/utils';
  import { CheckCircle2, XCircle, AlertCircle, Loader2, Clock } from '@lucide/svelte';
  import type { Status } from '@/lib/status';

  export let status: Status;
  export let label: string | undefined = undefined;
  export let className = '';

  const statusConfig: Record<
    Status,
    { icon: typeof CheckCircle2; label: string; classes: string }
  > = {
    success: {
      icon: CheckCircle2,
      label: 'Success',
      classes: 'bg-status-success-bg text-status-success',
    },
    error: {
      icon: XCircle,
      label: 'Failed',
      classes: 'bg-status-error-bg text-status-error',
    },
    warning: {
      icon: AlertCircle,
      label: 'Warning',
      classes: 'bg-status-warning-bg text-status-warning',
    },
    running: {
      icon: Loader2,
      label: 'Running',
      classes: 'bg-status-running-bg text-status-running',
    },
    pending: {
      icon: Clock,
      label: 'Pending',
      classes: 'bg-status-pending-bg text-status-pending',
    },
  };

  $: config = statusConfig[status];
  $: Icon = config.icon;
</script>

<span class={cn('status-badge', config.classes, className)}>
  <svelte:component
    this={Icon}
    class={cn('h-3 w-3', status === 'running' && 'animate-spin')}
  />
  <span>{label || config.label}</span>
</span>
