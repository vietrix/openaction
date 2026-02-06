<script lang="ts">
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
  import { link } from '@/lib/router';

  const pipelineData = {
    id: 'build-1247',
    project: 'forge-api',
    branch: 'main',
    commitHash: 'abc1234',
    commitMessage: 'feat: Add rate limiting to API endpoints',
    status: 'running' as const,
    totalDuration: '2m 34s',
    triggeredBy: 'john.doe',
    startedAt: '2024-01-15 14:32:01',
  };

  const pipelineNodes = [
    { id: 'checkout', name: 'Checkout', status: 'success' as const, duration: '3s', step: 'checkout', dependencies: [] },
    { id: 'install', name: 'Install', status: 'success' as const, duration: '45s', step: 'build', dependencies: ['checkout'] },
    { id: 'lint', name: 'Lint', status: 'success' as const, duration: '12s', step: 'test', dependencies: ['install'] },
    { id: 'test', name: 'Test', status: 'running' as const, duration: '2m 15s', step: 'test', dependencies: ['install'] },
    { id: 'build', name: 'Build', status: 'pending' as const, step: 'build', dependencies: ['lint', 'test'] },
    { id: 'deploy', name: 'Deploy', status: 'pending' as const, step: 'deploy', dependencies: ['build'] },
  ];

  type LogLevel = 'info' | 'warn' | 'error' | 'success' | 'debug';

  interface LogLine {
    id: number;
    timestamp: string;
    level: LogLevel;
    message: string;
  }

  const nodeLogs: Record<string, LogLine[]> = {
    checkout: [
      { id: 1, timestamp: '14:32:01', level: 'info', message: 'Cloning repository...' },
      { id: 2, timestamp: '14:32:02', level: 'info', message: 'Checking out branch: main' },
      { id: 3, timestamp: '14:32:03', level: 'info', message: 'HEAD is now at abc1234' },
      { id: 4, timestamp: '14:32:04', level: 'success', message: 'Checkout completed successfully' },
    ],
    install: [
      { id: 1, timestamp: '14:32:05', level: 'info', message: 'Running npm ci...' },
      { id: 2, timestamp: '14:32:10', level: 'info', message: 'Installing dependencies from package-lock.json' },
      { id: 3, timestamp: '14:32:45', level: 'info', message: 'added 847 packages in 40s' },
      { id: 4, timestamp: '14:32:50', level: 'success', message: 'Dependencies installed successfully' },
    ],
    lint: [
      { id: 1, timestamp: '14:32:51', level: 'info', message: 'Running eslint...' },
      { id: 2, timestamp: '14:33:00', level: 'info', message: 'Checking 127 files...' },
      { id: 3, timestamp: '14:33:03', level: 'success', message: 'No linting errors found' },
    ],
    test: [
      { id: 1, timestamp: '14:33:04', level: 'info', message: 'Running vitest...' },
      { id: 2, timestamp: '14:33:05', level: 'info', message: 'Starting test runner in parallel mode' },
      { id: 3, timestamp: '14:33:10', level: 'info', message: ' PASS  src/auth/auth.test.ts (12 tests)' },
      { id: 4, timestamp: '14:33:15', level: 'info', message: ' PASS  src/api/users.test.ts (8 tests)' },
      { id: 5, timestamp: '14:33:20', level: 'warn', message: ' SLOW  src/db/queries.test.ts took 1.2s' },
      { id: 6, timestamp: '14:33:25', level: 'info', message: ' PASS  src/db/queries.test.ts (15 tests)' },
      { id: 7, timestamp: '14:33:30', level: 'info', message: 'Running integration tests...' },
      { id: 8, timestamp: '14:33:35', level: 'info', message: ' RUN   src/integration/api.test.ts' },
    ],
    build: [],
    deploy: [],
  };

  const defaultLogs: LogLine[] = [
    { id: 1, timestamp: '14:32:01', level: 'info', message: 'Select a pipeline step to view logs' },
  ];

  let sidebarCollapsed = false;
  let selectedNode = 'test';
  let isLogOpen = false;

  $: currentLogs = nodeLogs[selectedNode]?.length ? nodeLogs[selectedNode] : defaultLogs;
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

      <div class="rounded-lg border border-border bg-card p-5 mb-6">
        <div class="flex flex-col lg:flex-row lg:items-start lg:justify-between gap-4">
          <div class="space-y-3">
            <div class="flex items-center gap-3">
              <h1 class="text-xl font-semibold text-foreground">
                {pipelineData.project}
                <span class="text-muted-foreground font-normal"> #{pipelineData.id.split('-')[1]}</span>
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
            isLogOpen = true;
          }}
        />
      </div>
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
              {selectedNode.charAt(0).toUpperCase() + selectedNode.slice(1)} - Logs
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
            logs={currentLogs}
            title={`${selectedNode.charAt(0).toUpperCase() + selectedNode.slice(1)} - Logs`}
            className="h-[60vh] rounded-lg"
          />
        </div>
        </div>
      </div>
    </div>
  {/if}
</div>
