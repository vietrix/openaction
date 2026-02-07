<script lang="ts">
  import { onMount } from 'svelte';
  import { link } from '@/lib/router';
  import { api } from '@/lib/api';
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

  type PipelineItem = {
    id: string;
    projectId: string;
    projectName: string;
    branch: string;
    commitHash: string;
    commitMessage: string;
    status: 'running' | 'success' | 'error' | 'pending';
    triggeredBy: string;
    duration: string;
    startedAt: string;
  };

  let pipelines: PipelineItem[] = [];
  let loading = true;
  let error = '';

  const formatDuration = (start?: number, end?: number) => {
    if (!start) return '—';
    const now = Math.floor(Date.now() / 1000);
    const finish = end && end > 0 ? end : now;
    const total = Math.max(0, finish - start);
    const mins = Math.floor(total / 60);
    const secs = total % 60;
    if (mins <= 0) return `${secs}s`;
    return `${mins}m ${secs}s`;
  };

  const formatRelativeTime = (timestamp?: number) => {
    if (!timestamp) return '—';
    const now = Math.floor(Date.now() / 1000);
    const diff = Math.max(0, now - timestamp);
    if (diff < 60) return `${diff}s ago`;
    if (diff < 3600) return `${Math.floor(diff / 60)}m ago`;
    if (diff < 86400) return `${Math.floor(diff / 3600)}h ago`;
    return `${Math.floor(diff / 86400)}d ago`;
  };

  const loadPipelines = async () => {
    loading = true;
    error = '';
    try {
      const projects = await api.getProjects();
      const lists = await Promise.all(
        projects.map(async (project) => ({
          project,
          pipelines: await api.getProjectPipelines(project.id),
        }))
      );
      pipelines = lists.flatMap(({ project, pipelines }) =>
        pipelines.map((item) => ({
          id: item.id,
          projectId: project.id,
          projectName: project.name,
          branch: item.branch,
          commitHash: item.commit_hash,
          commitMessage: item.commit_hash ? `Commit ${item.commit_hash}` : 'No commit message',
          status: item.status as PipelineItem['status'],
          triggeredBy: item.triggered_by,
          duration: formatDuration(item.started_at, item.finished_at),
          startedAt: formatRelativeTime(item.started_at),
        }))
      );
    } catch (err) {
      error = err instanceof Error ? err.message : 'Không thể tải pipeline';
    } finally {
      loading = false;
    }
  };

  onMount(loadPipelines);

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
      {#if loading}
        <div class="px-4 py-6 text-sm text-muted-foreground">Đang tải pipeline...</div>
      {:else if error}
        <div class="px-4 py-6 text-sm text-status-error">{error}</div>
      {:else if pipelines.length === 0}
        <div class="px-4 py-6 text-sm text-muted-foreground">Chưa có pipeline nào.</div>
      {:else}
        {#each pipelines as item (item.id)}
          <a
            href={`/pipelines/${item.projectId}/${item.id}`}
            use:link
            class="grid grid-cols-1 lg:grid-cols-[1.2fr_1fr_1fr_1fr_auto] gap-3 px-4 py-4 hover:bg-muted/20 transition-colors"
          >
            <div class="min-w-0">
              <p class="text-sm font-medium truncate">{item.projectName} #{item.id}</p>
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
      {/if}
    </div>
  </div>
</PageShell>
