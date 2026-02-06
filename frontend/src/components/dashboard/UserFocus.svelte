<script lang="ts">
  import { AlertCircle, Loader2, ExternalLink } from '@lucide/svelte';
  import { cn } from '@/lib/utils';
  import { link } from '@/lib/router';

  interface UserBuild {
    id: string;
    project: string;
    branch: string;
    status: 'running' | 'error';
    duration: string;
    startedAt: string;
  }

  export let builds: UserBuild[] = [];
  export let className = '';

  $: runningBuilds = builds.filter((build) => build.status === 'running');
  $: failedBuilds = builds.filter((build) => build.status === 'error');
</script>

{#if builds.length === 0}
  <div class={cn('rounded-lg border border-border bg-card p-4', className)}>
    <h3 class="text-sm font-medium text-foreground mb-3">Your Focus</h3>
    <div class="flex flex-col items-center justify-center py-8 text-center">
      <div class="h-10 w-10 rounded-full bg-muted/50 flex items-center justify-center mb-2">
        <AlertCircle class="h-5 w-5 text-muted-foreground" />
      </div>
      <p class="text-sm text-muted-foreground">All clear!</p>
      <p class="text-xs text-muted-foreground/70 mt-0.5">No running or failed builds</p>
    </div>
  </div>
{:else}
  <div class={cn('rounded-lg border border-border bg-card overflow-hidden', className)}>
    <div class="px-4 py-3 border-b border-border bg-muted/20">
      <h3 class="text-sm font-medium text-foreground">Your Focus</h3>
      <p class="text-xs text-muted-foreground mt-0.5">
        {runningBuilds.length} running · {failedBuilds.length} failed
      </p>
    </div>

    <div class="divide-y divide-border/50">
      {#each builds as build (build.id)}
        <a
          href={`/pipelines/${build.project}/${build.id}`}
          use:link
          class="flex items-center gap-3 px-4 py-3 hover:bg-muted/30 transition-colors group"
        >
          {#if build.status === 'running'}
            <Loader2 class="h-4 w-4 text-status-running animate-spin shrink-0" />
          {:else}
            <AlertCircle class="h-4 w-4 text-status-error shrink-0" />
          {/if}

          <div class="min-w-0 flex-1">
            <p class="text-sm font-medium text-foreground truncate">{build.project}</p>
            <p class="text-xs text-muted-foreground truncate">{build.branch}</p>
          </div>

          <div class="text-right shrink-0">
            <p class="text-xs font-mono text-muted-foreground">{build.duration}</p>
            <p class="text-[10px] text-muted-foreground/70">{build.startedAt}</p>
          </div>

          <ExternalLink class="h-3.5 w-3.5 text-muted-foreground opacity-0 group-hover:opacity-100 transition-opacity shrink-0" />
        </a>
      {/each}
    </div>
  </div>
{/if}
