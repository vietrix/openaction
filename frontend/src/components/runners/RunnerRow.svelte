<script lang="ts">
  import StatusDot from '@/components/ui/StatusDot.svelte';
  import DropdownMenu from '@/components/ui/DropdownMenu.svelte';
  import ResourceBar from './ResourceBar.svelte';
  import { cn } from '@/lib/utils';
  import { link } from '@/lib/router';
  import {
    MoreHorizontal,
    Loader2,
    ExternalLink,
    Eye,
    Pause,
    Trash2,
  } from '@lucide/svelte';

  type RunnerStatus = 'online' | 'offline' | 'busy';

  interface Runner {
    id: string;
    name: string;
    hostname: string;
    status: RunnerStatus;
    tags: string[];
    currentJob?: {
      id: string;
      project: string;
      step: string;
    };
    cpu: number;
    memory: number;
    lastSeen: string;
    version: string;
  }

  export let runner: Runner;

  const statusConfig: Record<RunnerStatus, { color: 'success' | 'error' | 'running' | 'pending'; label: string }> =
    {
      online: { color: 'success', label: 'Idle' },
      offline: { color: 'pending', label: 'Offline' },
      busy: { color: 'running', label: 'Busy' },
    };

  $: config = statusConfig[runner.status];
</script>

<div class="group flex items-center gap-4 px-4 py-4 border-b border-border/50 last:border-b-0 hover:bg-muted/20 transition-colors">
  <StatusDot status={config.color} pulse={runner.status === 'busy'} />

  <div class="min-w-0 flex-1">
    <div class="flex items-center gap-2">
      <span class="font-medium text-sm text-foreground font-mono">{runner.name}</span>
      <span
        class={cn(
          'text-[10px] px-1.5 py-0.5 rounded',
          runner.status === 'online' && 'bg-status-success-bg text-status-success',
          runner.status === 'offline' && 'bg-muted text-muted-foreground',
          runner.status === 'busy' && 'bg-status-running-bg text-status-running'
        )}
      >
        {config.label}
      </span>
    </div>
    <p class="text-xs text-muted-foreground font-mono mt-0.5 truncate">{runner.hostname}</p>
  </div>

  <div class="hidden lg:flex items-center gap-1.5 flex-wrap max-w-[200px]">
    {#each runner.tags.slice(0, 3) as tag (tag)}
      <span class="text-[10px] px-1.5 py-0.5 rounded bg-muted text-muted-foreground">
        {tag}
      </span>
    {/each}
    {#if runner.tags.length > 3}
      <span class="text-[10px] text-muted-foreground">+{runner.tags.length - 3}</span>
    {/if}
  </div>

  <div class="hidden md:block min-w-[150px]">
    {#if runner.currentJob}
      <a
        href={`/pipelines/${runner.currentJob.project}/${runner.currentJob.id}`}
        use:link
        class="flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors group/link"
      >
        <Loader2 class="h-3 w-3 animate-spin text-status-running" />
        <span class="truncate">{runner.currentJob.project}</span>
        <span class="text-muted-foreground/50">·</span>
        <span>{runner.currentJob.step}</span>
        <ExternalLink class="h-3 w-3 opacity-0 group-hover/link:opacity-100 transition-opacity" />
      </a>
    {:else}
      <span class="text-xs text-muted-foreground/50">—</span>
    {/if}
  </div>

  <div class="hidden xl:flex flex-col gap-1 min-w-[140px]">
    <ResourceBar value={runner.cpu} label="CPU" />
    <ResourceBar value={runner.memory} label="RAM" />
  </div>

  <div class="hidden sm:block text-right min-w-[80px]">
    <p class="text-xs text-muted-foreground">{runner.lastSeen}</p>
    <p class="text-[10px] text-muted-foreground/50 font-mono">{runner.version}</p>
  </div>

  <DropdownMenu align="end" className="w-48">
    <button
      slot="trigger"
      class="p-1.5 rounded-md hover:bg-muted text-muted-foreground hover:text-foreground transition-colors opacity-0 group-hover:opacity-100"
      aria-label="Open runner actions"
    >
      <MoreHorizontal class="h-4 w-4" />
    </button>
    <div slot="content">
      <button class="relative flex w-full cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground">
        <Eye class="h-4 w-4 mr-2" />
        View Details
      </button>
      <button class="relative flex w-full cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors hover:bg-accent hover:text-accent-foreground">
        <Pause class="h-4 w-4 mr-2" />
        Pause Runner
      </button>
      <div class="-mx-1 my-1 h-px bg-muted"></div>
      <button class="relative flex w-full cursor-default select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors text-status-error hover:bg-accent hover:text-status-error">
        <Trash2 class="h-4 w-4 mr-2" />
        Delete Runner
      </button>
    </div>
  </DropdownMenu>
</div>
