<script lang="ts">
  import { onMount } from 'svelte';
  import Sidebar from '@/components/layout/Sidebar.svelte';
  import TopBar from '@/components/layout/TopBar.svelte';
  import SystemSnapshot from '@/components/dashboard/SystemSnapshot.svelte';
  import ActivityStream from '@/components/dashboard/ActivityStream.svelte';
  import UserFocus from '@/components/dashboard/UserFocus.svelte';
  import PipelineVisualizer from '@/components/pipeline/PipelineVisualizer.svelte';
  import LiveLog from '@/components/terminal/LiveLog.svelte';
  import { Zap, GitCommit } from '@lucide/svelte';
  import { cn } from '@/lib/utils';

  const pipelineNodes = [
    { id: 'checkout', name: 'Checkout', status: 'success' as const, duration: '3s', step: 'checkout', dependencies: [] },
    { id: 'install', name: 'Install', status: 'success' as const, duration: '45s', step: 'build', dependencies: ['checkout'] },
    { id: 'lint', name: 'Lint', status: 'success' as const, duration: '12s', step: 'test', dependencies: ['install'] },
    { id: 'test', name: 'Test', status: 'running' as const, duration: '2m 15s', step: 'test', dependencies: ['install'] },
    { id: 'build', name: 'Build', status: 'pending' as const, step: 'build', dependencies: ['lint', 'test'] },
    { id: 'deploy', name: 'Deploy', status: 'pending' as const, step: 'deploy', dependencies: ['build'] },
  ];

  const initialLogs = [
    { id: 1, timestamp: '14:32:01', level: 'info' as const, message: '→ Starting pipeline: forge-api #1247' },
    { id: 2, timestamp: '14:32:01', level: 'info' as const, message: '→ Triggered by: push to main (abc1234)' },
    { id: 3, timestamp: '14:32:02', level: 'info' as const, message: '[checkout] Cloning repository...' },
    { id: 4, timestamp: '14:32:05', level: 'success' as const, message: '[checkout] SUCCESS - Completed in 3s' },
    { id: 5, timestamp: '14:32:06', level: 'info' as const, message: '[install] Running npm ci...' },
    { id: 6, timestamp: '14:32:51', level: 'success' as const, message: '[install] SUCCESS - Installed 847 packages in 45s' },
    { id: 7, timestamp: '14:32:52', level: 'info' as const, message: '[lint] Running eslint...' },
    { id: 8, timestamp: '14:33:04', level: 'success' as const, message: '[lint] SUCCESS - No issues found' },
    { id: 9, timestamp: '14:33:05', level: 'info' as const, message: '[test] Running vitest...' },
    { id: 10, timestamp: '14:33:45', level: 'info' as const, message: '[test] Running 127 tests across 24 suites...' },
    { id: 11, timestamp: '14:34:20', level: 'warn' as const, message: '[test] WARNING: 2 tests are slow (>500ms)' },
  ];

  const recentBuilds = [
    { id: '1247', project: 'forge-api', branch: 'main', commitHash: 'abc1234', commitMessage: 'feat: Add rate limiting to API endpoints', status: 'running' as const, triggeredBy: { name: 'john.doe' }, timeAgo: '2 min ago', duration: '2m 34s' },
    { id: '1246', project: 'forge-web', branch: 'feature/auth', commitHash: 'def5678', commitMessage: 'fix: Resolve OAuth callback issue', status: 'success' as const, triggeredBy: { name: 'jane.smith' }, timeAgo: '15 min ago', duration: '4m 12s' },
    { id: '1245', project: 'forge-cli', branch: 'main', commitHash: 'ghi9012', commitMessage: 'chore: Update dependencies', status: 'success' as const, triggeredBy: { name: 'bot' }, timeAgo: '32 min ago', duration: '1m 58s' },
    { id: '1244', project: 'forge-api', branch: 'fix/memory-leak', commitHash: 'jkl3456', commitMessage: 'fix: Memory leak in connection pool', status: 'error' as const, triggeredBy: { name: 'john.doe' }, timeAgo: '1 hr ago', duration: '3m 45s' },
    { id: '1243', project: 'forge-infra', branch: 'main', commitHash: 'mno7890', commitMessage: 'infra: Scale up worker nodes', status: 'success' as const, triggeredBy: { name: 'ops-bot' }, timeAgo: '2 hr ago', duration: '8m 22s' },
    { id: '1242', project: 'forge-web', branch: 'main', commitHash: 'pqr1234', commitMessage: 'feat: Add dark mode toggle', status: 'success' as const, triggeredBy: { name: 'jane.smith' }, timeAgo: '3 hr ago', duration: '5m 15s' },
    { id: '1241', project: 'forge-api', branch: 'main', commitHash: 'stu5678', commitMessage: 'docs: Update API documentation', status: 'success' as const, triggeredBy: { name: 'john.doe' }, timeAgo: '4 hr ago', duration: '2m 08s' },
  ];

  const userBuilds = [
    { id: '1247', project: 'forge-api', branch: 'main', status: 'running' as const, duration: '2m 34s', startedAt: '14:32' },
    { id: '1244', project: 'forge-api', branch: 'fix/memory-leak', status: 'error' as const, duration: '3m 45s', startedAt: '13:45' },
  ];

  const systemStats = { success: 42, failed: 3, running: 2 };
  const durationTrend = [4.2, 4.5, 4.1, 3.8, 4.0, 3.9, 3.7];
  const failureRateTrend = [5.2, 3.1, 8.5, 2.0, 6.8, 4.2, 6.4];

  let sidebarCollapsed = false;
  let logs = [...initialLogs];

  onMount(() => {
    const interval = setInterval(() => {
      const newLogMessages = [
        '[test] Test suite: auth.test.ts - 12 passed',
        '[test] Test suite: api.test.ts - 8 passed',
        '[test] Running integration tests...',
        '[test] Database connection established',
      ];

      logs = (() => {
        if (logs.length > 50) return logs;
        const randomMessage = newLogMessages[Math.floor(Math.random() * newLogMessages.length)];
        const now = new Date();
        const timestamp = now.toLocaleTimeString('en-US', { hour12: false }).slice(0, 8);
        return [
          ...logs,
          {
            id: logs.length + 1,
            timestamp,
            level: 'info' as const,
            message: randomMessage,
          },
        ];
      })();
    }, 3000);

    return () => clearInterval(interval);
  });
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
      <div class="mb-6">
        <h1 class="text-xl font-semibold text-foreground">Dashboard</h1>
        <p class="text-sm text-muted-foreground mt-1">Monitor your CI/CD pipelines in real-time</p>
      </div>

      <SystemSnapshot
        stats={systemStats}
        durationTrend={durationTrend}
        failureRateTrend={failureRateTrend}
        className="mb-6"
      />

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-6">
        <ActivityStream builds={recentBuilds} className="lg:col-span-2" />
        <UserFocus builds={userBuilds} />
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div class="rounded-lg border border-border bg-card p-5">
          <div class="flex items-center justify-between mb-4">
            <div>
              <h2 class="text-sm font-medium text-foreground flex items-center gap-2">
                <Zap class="h-4 w-4 text-status-running" />
                Active Pipeline
              </h2>
              <p class="text-xs text-muted-foreground mt-0.5">
                forge-api #1247 · main · abc1234
              </p>
            </div>
            <div class="flex items-center gap-2 text-xs text-muted-foreground">
              <GitCommit class="h-3.5 w-3.5" />
              <span class="font-mono">abc1234</span>
            </div>
          </div>
          <PipelineVisualizer nodes={pipelineNodes} />
        </div>

        <LiveLog logs={logs} title="forge-api #1247 — Live Logs" />
      </div>
    </div>
  </main>
</div>
