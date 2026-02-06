<script lang="ts">
  import { GitBranch, GitCommit, Clock, User, Play, FileText, RotateCcw } from '@lucide/svelte';
  import StatusDot from '@/components/ui/StatusDot.svelte';
  import { cn } from '@/lib/utils';
  import { link } from '@/lib/router';

  interface Build {
    id: string;
    project: string;
    branch: string;
    commitHash: string;
    commitMessage: string;
    status: 'success' | 'error' | 'warning' | 'running' | 'pending';
    triggeredBy: {
      name: string;
      avatar?: string;
    };
    timeAgo: string;
    duration: string;
  }

  export let builds: Build[] = [];
  export let className = '';
</script>

<div class={cn('rounded-lg border border-border bg-card overflow-hidden', className)}>
  <div class="flex items-center justify-between px-4 py-3 border-b border-border bg-muted/20">
    <div class="flex items-center gap-2">
      <Play class="h-4 w-4 text-muted-foreground" />
      <h3 class="text-sm font-medium text-foreground">Recent Activity</h3>
      <span class="text-xs text-muted-foreground bg-muted px-1.5 py-0.5 rounded">
        {builds.length}
      </span>
    </div>
    <button class="text-xs text-muted-foreground hover:text-foreground transition-colors">
      View all →
    </button>
  </div>

  <div class="max-h-[400px] overflow-y-auto">
    {#each builds as build (build.id)}
      <a
        href={`/pipelines/${build.project}/${build.id}`}
        use:link
        class={cn(
          'group flex items-center gap-3 px-4 py-3 border-b border-border/50 last:border-b-0',
          'transition-colors duration-150 hover:bg-muted/30 cursor-pointer relative'
        )}
      >
        <StatusDot status={build.status} pulse={build.status === 'running'} className="shrink-0" />

        <div class="min-w-0 flex-1">
          <div class="flex items-center gap-2">
            <span class="font-medium text-sm text-foreground truncate">
              {build.project}
            </span>
            <span class="flex items-center gap-1 text-xs text-muted-foreground">
              <GitBranch class="h-3 w-3" />
              <span class="truncate max-w-[120px]">{build.branch}</span>
            </span>
          </div>
          <p class="text-xs text-muted-foreground truncate mt-0.5">
            {build.commitMessage}
          </p>
        </div>

        <div class="hidden md:flex items-center gap-1.5 text-xs text-muted-foreground shrink-0">
          <GitCommit class="h-3 w-3" />
          <span class="font-mono">{build.commitHash}</span>
        </div>

        <div class="hidden lg:flex items-center gap-2 text-xs text-muted-foreground shrink-0 min-w-[100px]">
          {#if build.triggeredBy.avatar}
            <img
              src={build.triggeredBy.avatar}
              alt={build.triggeredBy.name}
              class="h-5 w-5 rounded-full"
            />
          {:else}
            <div class="h-5 w-5 rounded-full bg-muted flex items-center justify-center">
              <User class="h-3 w-3" />
            </div>
          {/if}
          <span class="truncate">{build.triggeredBy.name}</span>
        </div>

        <div class="flex items-center gap-4 text-xs text-muted-foreground shrink-0">
          <span class="hidden sm:block">{build.timeAgo}</span>
          <span class="flex items-center gap-1 font-mono">
            <Clock class="h-3 w-3" />
            {build.duration}
          </span>
        </div>

        <div
          class={cn(
            'absolute right-3 top-1/2 -translate-y-1/2 flex items-center gap-1',
            'transition-opacity duration-150 opacity-0 pointer-events-none',
            'group-hover:opacity-100 group-hover:pointer-events-auto'
          )}
        >
          <button
            on:click|preventDefault|stopPropagation
            class="p-1.5 rounded-md bg-muted/80 hover:bg-muted text-muted-foreground hover:text-foreground transition-colors"
            title="View Logs"
          >
            <FileText class="h-3.5 w-3.5" />
          </button>
          <button
            on:click|preventDefault|stopPropagation
            class="p-1.5 rounded-md bg-muted/80 hover:bg-muted text-muted-foreground hover:text-foreground transition-colors"
            title="Re-run"
          >
            <RotateCcw class="h-3.5 w-3.5" />
          </button>
        </div>
      </a>
    {/each}
  </div>
</div>
