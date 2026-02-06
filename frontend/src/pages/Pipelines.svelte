<script lang="ts">
  import { link } from '@/lib/router';
  import {
    Workflow,
    Play,
    CheckCircle2,
    XCircle,
    Loader2,
    GitBranch,
    Clock3,
    GitCommit,
  } from '@lucide/svelte';
  import StatusBadge from '@/components/ui/StatusBadge.svelte';
  import PageShell from '@/components/layout/PageShell.svelte';

  const pipelines = [
    {
      id: '1251',
      project: 'forge-api',
      branch: 'main',
      commitHash: 'd4f6a2c',
      commitMessage: 'feat: add OpenAPI validation middleware',
      status: 'running' as const,
      triggeredBy: 'minh',
      duration: '1m 22s',
      startedAt: '2 min ago',
    },
    {
      id: '1250',
      project: 'forge-web',
      branch: 'feature/billing',
      commitHash: '8ab13ef',
      commitMessage: 'fix: prevent duplicate checkout session',
      status: 'success' as const,
      triggeredBy: 'linh',
      duration: '4m 03s',
      startedAt: '11 min ago',
    },
    {
      id: '1249',
      project: 'forge-worker',
      branch: 'main',
      commitHash: '6e20da1',
      commitMessage: 'chore: bump redis client to v5',
      status: 'error' as const,
      triggeredBy: 'bot',
      duration: '2m 48s',
      startedAt: '22 min ago',
    },
    {
      id: '1248',
      project: 'forge-infra',
      branch: 'main',
      commitHash: 'a78c09b',
      commitMessage: 'infra: optimize autoscaling thresholds',
      status: 'success' as const,
      triggeredBy: 'ops-bot',
      duration: '7m 16s',
      startedAt: '48 min ago',
    },
  ];

  $: total = pipelines.length;
  $: running = pipelines.filter((item) => item.status === 'running').length;
  $: success = pipelines.filter((item) => item.status === 'success').length;
  $: failed = pipelines.filter((item) => item.status === 'error').length;
</script>

<PageShell title="Pipelines" description="Quản lý, theo dõi và điều phối toàn bộ pipeline CI/CD">
  <div class="grid grid-cols-2 md:grid-cols-4 gap-3 mb-6">
    <div class="rounded-lg border border-border bg-card p-4">
      <div class="flex items-center gap-2 text-muted-foreground">
        <Workflow class="h-4 w-4" />
        <span class="text-xs uppercase tracking-wide">Total</span>
      </div>
      <p class="text-2xl font-semibold mt-1">{total}</p>
    </div>
    <div class="rounded-lg border border-border bg-card p-4">
      <div class="flex items-center gap-2 text-muted-foreground">
        <Loader2 class="h-4 w-4 text-status-running animate-spin" />
        <span class="text-xs uppercase tracking-wide">Running</span>
      </div>
      <p class="text-2xl font-semibold mt-1">{running}</p>
    </div>
    <div class="rounded-lg border border-border bg-card p-4">
      <div class="flex items-center gap-2 text-muted-foreground">
        <CheckCircle2 class="h-4 w-4 text-status-success" />
        <span class="text-xs uppercase tracking-wide">Success</span>
      </div>
      <p class="text-2xl font-semibold mt-1">{success}</p>
    </div>
    <div class="rounded-lg border border-border bg-card p-4">
      <div class="flex items-center gap-2 text-muted-foreground">
        <XCircle class="h-4 w-4 text-status-error" />
        <span class="text-xs uppercase tracking-wide">Failed</span>
      </div>
      <p class="text-2xl font-semibold mt-1">{failed}</p>
    </div>
  </div>

  <div class="rounded-lg border border-border bg-card overflow-hidden">
    <div class="px-4 py-3 border-b border-border bg-muted/20 flex items-center justify-between">
      <h2 class="text-sm font-medium">Recent Pipeline Runs</h2>
      <button class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs bg-primary text-primary-foreground hover:bg-primary/90 transition-colors">
        <Play class="h-3.5 w-3.5" />
        Run Pipeline
      </button>
    </div>

    <div class="divide-y divide-border/50">
      {#each pipelines as item (item.id)}
        <a
          href={`/pipelines/${item.project}/${item.id}`}
          use:link
          class="grid grid-cols-1 lg:grid-cols-[1.2fr_1fr_1fr_1fr_auto] gap-3 px-4 py-4 hover:bg-muted/20 transition-colors"
        >
          <div class="min-w-0">
            <p class="text-sm font-medium truncate">{item.project} #{item.id}</p>
            <p class="text-xs text-muted-foreground truncate mt-0.5">{item.commitMessage}</p>
          </div>

          <div class="text-xs text-muted-foreground flex items-center gap-1.5">
            <GitBranch class="h-3.5 w-3.5" />
            <span class="truncate">{item.branch}</span>
          </div>

          <div class="text-xs text-muted-foreground flex items-center gap-1.5 font-mono">
            <GitCommit class="h-3.5 w-3.5" />
            {item.commitHash}
          </div>

          <div class="text-xs text-muted-foreground flex items-center gap-1.5">
            <Clock3 class="h-3.5 w-3.5" />
            {item.duration} • {item.startedAt}
          </div>

          <StatusBadge status={item.status} />
        </a>
      {/each}
    </div>
  </div>
</PageShell>
