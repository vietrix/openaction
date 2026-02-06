<script lang="ts">
  import { cn } from '@/lib/utils';
  import {
    CheckCircle2,
    XCircle,
    Loader2,
    Clock,
    Package,
    TestTube,
    Rocket,
    GitBranch,
  } from '@lucide/svelte';

  export type NodeStatus = 'success' | 'error' | 'running' | 'pending';

  export let id: string;
  export let name: string;
  export let status: NodeStatus;
  export let duration: string | undefined = undefined;
  export let step: string | undefined = undefined;
  export let className = '';

  const statusConfig: Record<
    NodeStatus,
    { icon: typeof CheckCircle2; borderClass: string; bgClass: string }
  > = {
    success: {
      icon: CheckCircle2,
      borderClass: 'border-status-success/50',
      bgClass: 'bg-status-success-bg',
    },
    error: {
      icon: XCircle,
      borderClass: 'border-status-error/50',
      bgClass: 'bg-status-error-bg',
    },
    running: {
      icon: Loader2,
      borderClass: 'border-status-running/50 shadow-glow-sm',
      bgClass: 'bg-status-running-bg',
    },
    pending: {
      icon: Clock,
      borderClass: 'border-border',
      bgClass: 'bg-muted/50',
    },
  };

  const stepIcons: Record<string, typeof Package> = {
    checkout: GitBranch,
    build: Package,
    test: TestTube,
    deploy: Rocket,
  };

  $: config = statusConfig[status];
  $: StatusIcon = config.icon;
  $: StepIcon = step ? stepIcons[step.toLowerCase()] || Package : Package;
</script>

<div
  data-node-id={id}
  class={cn(
    'relative flex min-w-[160px] flex-col gap-2 rounded-lg border bg-card p-3 transition-all duration-150',
    config.borderClass,
    status === 'running' && 'node-pulse',
    className
  )}
>
  <div class={cn('absolute inset-0 rounded-lg opacity-30', config.bgClass)}></div>

  <div class="relative z-10">
    <div class="flex items-center justify-between gap-2">
      <div class="flex items-center gap-2">
        <svelte:component this={StepIcon} class="h-4 w-4 text-muted-foreground" />
        <span class="text-sm font-medium text-foreground">{name}</span>
      </div>
      <svelte:component
        this={StatusIcon}
        class={cn(
          'h-4 w-4',
          status === 'success' && 'text-status-success',
          status === 'error' && 'text-status-error',
          status === 'running' && 'text-status-running animate-spin',
          status === 'pending' && 'text-muted-foreground'
        )}
      />
    </div>

    <div class="mt-1.5 flex items-center gap-3 text-xs text-muted-foreground">
      {#if step}
        <span class="capitalize">{step}</span>
      {/if}
      {#if duration}
        <span class="font-mono">{duration}</span>
      {/if}
    </div>
  </div>
</div>
