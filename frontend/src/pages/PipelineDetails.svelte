<script lang="ts">
  import { onMount } from 'svelte';
  import Sidebar from '@/components/layout/Sidebar.svelte';
  import TopBar from '@/components/layout/TopBar.svelte';
  import PipelineVisualizer from '@/components/pipeline/PipelineVisualizer.svelte';
  import LiveLog from '@/components/terminal/LiveLog.svelte';
  import StatusBadge from '@/components/ui/StatusBadge.svelte';
  import {
    ArrowLeft,
    RotateCcw,
    XCircle,
    Clock,
    GitCommit,
    GitBranch,
    User,
    Calendar,
    X,
    FileText,
  } from '@lucide/svelte';
  import { cn } from '@/lib/utils';
  import { link, location } from '@/lib/router';
  import { api, buildLogStreamUrl } from '@/lib/api';

  type LogLevel = 'info' | 'warn' | 'error' | 'success' | 'debug';

  interface LogLine {
    id: number;
    timestamp: string;
    level: LogLevel;
    message: string;
  }

  type Pipeline = {
    id: string;
    project_id: string;
    status: 'success' | 'error' | 'running' | 'pending';
    commit_hash: string;
    branch: string;
    triggered_by: string;
    started_at: number;
    finished_at: number;
  };

  type PipelineStep = {
    id: string;
    name: string;
    status: 'success' | 'error' | 'running' | 'pending';
    started_at: number;
    finished_at: number;
    log_path: string;
  };

  let sidebarCollapsed = false;
  let pipelineId = '';
  let projectId = '';
  let pipeline: Pipeline | null = null;
  let pipelineError = '';
  let steps: PipelineStep[] = [];
  let pipelineNodes: Array<{
    id: string;
    name: string;
    status: PipelineStep['status'];
    duration?: string;
    step?: string;
    dependencies: string[];
  }> = [];

  let selectedNode = '';
  let isLogOpen = false;
  let logLines: LogLine[] = [];
  let logLoading = false;
  let logError = '';
  let streamUrl = '';
  let lastPath = '';

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

  const formatDateTime = (timestamp?: number) => {
    if (!timestamp) return '—';
    return new Date(timestamp * 1000).toLocaleString('vi-VN');
  };

  const parseLogLevel = (line: string): LogLevel => {
    if (/ERROR|FAILED|FAIL/i.test(line)) return 'error';
    if (/WARN|WARNING/i.test(line)) return 'warn';
    if (/SUCCESS|PASSED|OK/i.test(line)) return 'success';
    if (/DEBUG/i.test(line)) return 'debug';
    return 'info';
  };

  const buildLogLines = (raw: string) => {
    const base = new Date().toLocaleTimeString('en-US', { hour12: false }).slice(0, 8);
    const lines = raw
      .split('\n')
      .map((line) => line.trim())
      .filter(Boolean);
    if (lines.length === 0) {
      return [
        {
          id: 1,
          timestamp: base,
          level: 'info' as const,
          message: 'Chưa có log cho bước này.',
        },
      ];
    }
    return lines.map((line, index) => ({
      id: index + 1,
      timestamp: base,
      level: parseLogLevel(line),
      message: line,
    }));
  };

  const loadPipeline = async () => {
    if (!pipelineId) return;
    pipelineError = '';
    try {
      pipeline = await api.getPipeline(pipelineId);
      steps = await api.getPipelineSteps(pipelineId);
      const ordered = [...steps].sort((a, b) => {
        const aTime = a.started_at || 0;
        const bTime = b.started_at || 0;
        if (aTime === bTime) return a.name.localeCompare(b.name);
        return aTime - bTime;
      });
      pipelineNodes = ordered.map((step, index) => ({
        id: step.id,
        name: step.name,
        status: step.status,
        duration: formatDuration(step.started_at, step.finished_at),
        step: step.name,
        dependencies: index === 0 ? [] : [ordered[index - 1].id],
      }));
      if (!selectedNode && ordered.length > 0) {
        selectedNode = ordered[0].id;
      }
      if (selectedNode) {
        await loadLogs(selectedNode);
      }
    } catch (err) {
      pipelineError = err instanceof Error ? err.message : 'Không thể tải pipeline';
    }
  };

  const loadLogs = async (stepId: string) => {
    const step = steps.find((item) => item.id === stepId);
    if (!step || !pipelineId) return;
    logLoading = true;
    logError = '';
    streamUrl = step.log_path ? buildLogStreamUrl(pipelineId, step.log_path) : '';
    try {
      const raw = step.log_path ? await api.getLogSnapshot(pipelineId, step.id) : '';
      logLines = buildLogLines(raw);
    } catch (err) {
      logError = err instanceof Error ? err.message : 'Không thể tải log';
      logLines = [
        {
          id: 1,
          timestamp: new Date().toLocaleTimeString('en-US', { hour12: false }).slice(0, 8),
          level: 'error',
          message: logError,
        },
      ];
    } finally {
      logLoading = false;
    }
  };

  onMount(() => {
    const match = $location.pathname.match(/^\/pipelines\/([^/]+)\/([^/]+)$/);
    if (match) {
      projectId = match[1];
      pipelineId = match[2];
      lastPath = $location.pathname;
      loadPipeline();
    }
  });

  $: if ($location.pathname !== lastPath) {
    const match = $location.pathname.match(/^\/pipelines\/([^/]+)\/([^/]+)$/);
    if (match) {
      projectId = match[1];
      pipelineId = match[2];
      lastPath = $location.pathname;
      selectedNode = '';
      loadPipeline();
    }
  }

  $: selectedStep = steps.find((item) => item.id === selectedNode);
  $: selectedStepName = selectedStep?.name ?? 'Logs';

  $: pipelineData = pipeline
    ? {
        id: pipeline.id,
        project: projectId,
        branch: pipeline.branch,
        commitHash: pipeline.commit_hash,
        commitMessage: pipeline.commit_hash ? `Commit ${pipeline.commit_hash}` : 'No commit message',
        status: pipeline.status,
        totalDuration: formatDuration(pipeline.started_at, pipeline.finished_at),
        triggeredBy: pipeline.triggered_by,
        startedAt: formatDateTime(pipeline.started_at),
      }
    : null;
</script>

<div class="dark min-h-screen bg-background text-foreground">
  <Sidebar collapsed={sidebarCollapsed} onToggle={() => (sidebarCollapsed = !sidebarCollapsed)} />
  <TopBar {sidebarCollapsed} />

  <main
    class={cn(
      'min-h-screen pt-14 transition-all duration-200',
      sidebarCollapsed ? 'pl-14' : 'pl-56'
    )}
  >
    <div class="p-6">
      <a
        href="/"
        use:link
        class="inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors mb-4"
      >
        <ArrowLeft class="h-4 w-4" />
        Back to Dashboard
      </a>

      {#if pipelineError}
        <div class="rounded-lg border border-status-error/40 bg-status-error-bg p-6 text-sm text-status-error">
          {pipelineError}
        </div>
      {:else if !pipelineData}
        <div class="rounded-lg border border-border bg-card p-6 text-sm text-muted-foreground">
          Đang tải pipeline...
        </div>
      {:else}
        <div class="rounded-lg border border-border bg-card p-5 mb-6">
          <div class="flex flex-col lg:flex-row lg:items-start lg:justify-between gap-4">
            <div class="space-y-3">
              <div class="flex items-center gap-3">
                <h1 class="text-xl font-semibold text-foreground">
                  {pipelineData.project}
                  <span class="text-muted-foreground font-normal"> #{pipelineData.id.slice(0, 8)}</span>
                </h1>
                <StatusBadge status={pipelineData.status} />
              </div>

              <p class="text-sm text-foreground">{pipelineData.commitMessage}</p>

              <div class="flex flex-wrap items-center gap-4 text-xs text-muted-foreground">
                <span class="flex items-center gap-1.5">
                  <GitBranch class="h-3.5 w-3.5" />
                  {pipelineData.branch}
                </span>
                <span class="flex items-center gap-1.5 font-mono">
                  <GitCommit class="h-3.5 w-3.5" />
                  {pipelineData.commitHash}
                </span>
                <span class="flex items-center gap-1.5">
                  <User class="h-3.5 w-3.5" />
                  {pipelineData.triggeredBy}
                </span>
                <span class="flex items-center gap-1.5">
                  <Calendar class="h-3.5 w-3.5" />
                  {pipelineData.startedAt}
                </span>
                <span class="flex items-center gap-1.5 font-mono">
                  <Clock class="h-3.5 w-3.5" />
                  {pipelineData.totalDuration}
                </span>
              </div>
            </div>

            <div class="flex items-center gap-2 shrink-0">
              <button class="btn-secondary inline-flex items-center gap-2 px-3 py-2 rounded-md border border-border bg-muted/50 text-sm font-medium text-foreground hover:bg-muted transition-colors">
                <RotateCcw class="h-4 w-4" />
                Re-run
              </button>
              <button class="btn-secondary inline-flex items-center gap-2 px-3 py-2 rounded-md border border-status-error/30 bg-status-error-bg text-sm font-medium text-status-error hover:bg-status-error/20 transition-colors">
                <XCircle class="h-4 w-4" />
                Cancel
              </button>
            </div>
          </div>
        </div>

        <div class="rounded-lg border border-border bg-card p-5">
          <div class="flex items-center justify-between mb-4">
            <div>
              <h2 class="text-sm font-medium text-foreground">Workflow Graph</h2>
              <p class="text-xs text-muted-foreground mt-0.5">Click a step to view logs</p>
            </div>
            <button
              class="inline-flex items-center gap-2 px-3 py-1.5 rounded-md text-xs border border-border text-muted-foreground hover:text-foreground hover:bg-muted transition-colors"
              on:click={() => (isLogOpen = true)}
            >
              <FileText class="h-3.5 w-3.5" />
              View Logs
            </button>
          </div>
          <PipelineVisualizer
            nodes={pipelineNodes}
            selectedNode={selectedNode}
            on:select={(event) => {
              selectedNode = event.detail;
              loadLogs(selectedNode);
              isLogOpen = true;
            }}
          />
        </div>
      {/if}
    </div>
  </main>

  {#if isLogOpen}
    <div class="fixed inset-0 z-50">
      <button
        type="button"
        class="absolute inset-0 bg-black/60 backdrop-blur-sm"
        on:click={() => (isLogOpen = false)}
        aria-label="Close logs"
      ></button>
      <div class="absolute inset-0 flex items-center justify-center p-4">
        <div class="w-[92vw] max-w-5xl max-h-[82vh] rounded-xl border border-border bg-card shadow-xl">
          <div class="flex items-center justify-between px-4 py-3 border-b border-border">
            <div class="flex items-center gap-2">
              <FileText class="h-4 w-4 text-muted-foreground" />
              <div class="text-sm font-medium text-foreground">
                {selectedStepName} - Logs
              </div>
            </div>
            <button
              class="h-8 w-8 inline-flex items-center justify-center rounded-md hover:bg-muted text-muted-foreground hover:text-foreground transition-colors"
              on:click={() => (isLogOpen = false)}
              aria-label="Close logs"
            >
              <X class="h-4 w-4" />
            </button>
          </div>
          <div class="p-4">
            <LiveLog
              logs={logLines}
              title={`${selectedStepName} - Logs`}
              className="h-[60vh] rounded-lg"
              streamUrl={streamUrl}
              loading={logLoading}
            />
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>
