<script lang="ts">
  import Sidebar from '@/components/layout/Sidebar.svelte';
  import TopBar from '@/components/layout/TopBar.svelte';
  import RunnerRow from '@/components/runners/RunnerRow.svelte';
  import { Server, Plus, CheckCircle2, Clock, Loader2 } from '@lucide/svelte';
  import { cn } from '@/lib/utils';

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

  const runners: Runner[] = [
    {
      id: 'runner-01',
      name: 'runner-prod-01',
      hostname: 'ip-10-0-1-101.ec2.internal',
      status: 'busy',
      tags: ['linux', 'x64', 'docker', 'production'],
      currentJob: { id: 'build-1247', project: 'forge-api', step: 'test' },
      cpu: 67,
      memory: 45,
      lastSeen: 'now',
      version: 'v2.4.1',
    },
    {
      id: 'runner-02',
      name: 'runner-prod-02',
      hostname: 'ip-10-0-1-102.ec2.internal',
      status: 'online',
      tags: ['linux', 'x64', 'docker', 'production'],
      cpu: 12,
      memory: 28,
      lastSeen: 'now',
      version: 'v2.4.1',
    },
    {
      id: 'runner-03',
      name: 'runner-prod-03',
      hostname: 'ip-10-0-1-103.ec2.internal',
      status: 'busy',
      tags: ['linux', 'x64', 'docker', 'production', 'high-mem'],
      currentJob: { id: 'build-1245', project: 'forge-web', step: 'build' },
      cpu: 89,
      memory: 72,
      lastSeen: 'now',
      version: 'v2.4.1',
    },
    {
      id: 'runner-04',
      name: 'runner-gpu-01',
      hostname: 'gpu-node-01.internal',
      status: 'online',
      tags: ['linux', 'x64', 'docker', 'gpu-enabled', 'cuda'],
      cpu: 5,
      memory: 15,
      lastSeen: 'now',
      version: 'v2.4.1',
    },
    {
      id: 'runner-05',
      name: 'runner-staging-01',
      hostname: 'staging-runner.internal',
      status: 'offline',
      tags: ['linux', 'x64', 'docker', 'staging'],
      cpu: 0,
      memory: 0,
      lastSeen: '15 min ago',
      version: 'v2.4.0',
    },
    {
      id: 'runner-06',
      name: 'runner-macos-01',
      hostname: 'mac-mini-01.local',
      status: 'online',
      tags: ['macos', 'arm64', 'ios-build'],
      cpu: 8,
      memory: 32,
      lastSeen: 'now',
      version: 'v2.4.1',
    },
  ];

  let sidebarCollapsed = false;

  $: onlineCount = runners.filter((runner) => runner.status === 'online' || runner.status === 'busy').length;
  $: busyCount = runners.filter((runner) => runner.status === 'busy').length;
  $: offlineCount = runners.filter((runner) => runner.status === 'offline').length;
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
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
        <div>
          <h1 class="text-xl font-semibold text-foreground flex items-center gap-2">
            <Server class="h-5 w-5" />
            Runners
          </h1>
          <p class="text-sm text-muted-foreground mt-1">
            Manage your CI/CD runner pool
          </p>
        </div>
        <button class="btn-primary inline-flex items-center gap-2 px-4 py-2 rounded-md bg-primary text-primary-foreground text-sm font-medium hover:bg-primary/90 transition-colors self-start">
          <Plus class="h-4 w-4" />
          Add Runner
        </button>
      </div>

      <div class="grid grid-cols-2 md:grid-cols-4 gap-3 mb-6">
        <div class="rounded-lg border border-border bg-card p-4">
          <div class="flex items-center gap-2">
            <Server class="h-4 w-4 text-muted-foreground" />
            <span class="text-xs text-muted-foreground uppercase tracking-wide">Total</span>
          </div>
          <p class="text-2xl font-semibold text-foreground mt-1">{runners.length}</p>
        </div>
        <div class="rounded-lg border border-border bg-card p-4">
          <div class="flex items-center gap-2">
            <CheckCircle2 class="h-4 w-4 text-status-success" />
            <span class="text-xs text-muted-foreground uppercase tracking-wide">Online</span>
          </div>
          <p class="text-2xl font-semibold text-foreground mt-1">{onlineCount}</p>
        </div>
        <div class="rounded-lg border border-border bg-card p-4">
          <div class="flex items-center gap-2">
            <Loader2 class="h-4 w-4 text-status-running animate-spin" />
            <span class="text-xs text-muted-foreground uppercase tracking-wide">Busy</span>
          </div>
          <p class="text-2xl font-semibold text-foreground mt-1">{busyCount}</p>
        </div>
        <div class="rounded-lg border border-border bg-card p-4">
          <div class="flex items-center gap-2">
            <Clock class="h-4 w-4 text-muted-foreground" />
            <span class="text-xs text-muted-foreground uppercase tracking-wide">Offline</span>
          </div>
          <p class="text-2xl font-semibold text-foreground mt-1">{offlineCount}</p>
        </div>
      </div>

      <div class="rounded-lg border border-border bg-card overflow-hidden">
        <div class="px-4 py-3 border-b border-border bg-muted/20">
          <div class="flex items-center justify-between">
            <h2 class="text-sm font-medium text-foreground">Runner Pool</h2>
            <span class="text-xs text-muted-foreground">
              {onlineCount}/{runners.length} online, {busyCount} busy
            </span>
          </div>
        </div>

        <div class="divide-y divide-border/50">
          {#each runners as runner (runner.id)}
            <RunnerRow {runner} />
          {/each}
        </div>
      </div>
    </div>
  </main>
</div>
