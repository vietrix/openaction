<script lang="ts">
  import { GitBranch, ShieldCheck, ExternalLink, RefreshCw } from '@lucide/svelte';
  import PageShell from '@/components/layout/PageShell.svelte';

  const repos = [
    { name: 'forge-api', provider: 'GitHub', branch: 'main', protection: 'required checks', lastSync: '1 min ago' },
    { name: 'forge-web', provider: 'GitHub', branch: 'main', protection: 'required reviews', lastSync: '4 min ago' },
    { name: 'forge-worker', provider: 'GitLab', branch: 'main', protection: 'required checks', lastSync: '9 min ago' },
    { name: 'forge-infra', provider: 'GitHub', branch: 'main', protection: 'admin only', lastSync: '14 min ago' },
  ];
</script>

<PageShell title="Repositories" description="Kết nối và quản trị repository dùng cho build/deploy">
  <div class="rounded-lg border border-border bg-card overflow-hidden">
    <div class="px-4 py-3 border-b border-border bg-muted/20 flex items-center justify-between">
      <h2 class="text-sm font-medium">Connected Repositories</h2>
      <button class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs border border-border hover:bg-muted transition-colors">
        <RefreshCw class="h-3.5 w-3.5" />
        Sync now
      </button>
    </div>

    <div class="divide-y divide-border/50">
      {#each repos as repo (repo.name)}
        <div class="grid grid-cols-1 md:grid-cols-[1.2fr_0.8fr_1fr_1fr_auto] gap-3 px-4 py-4">
          <div class="min-w-0">
            <p class="text-sm font-medium truncate">{repo.name}</p>
            <p class="text-xs text-muted-foreground mt-0.5">{repo.provider}</p>
          </div>
          <div class="text-xs text-muted-foreground flex items-center gap-1.5">
            <GitBranch class="h-3.5 w-3.5" />
            {repo.branch}
          </div>
          <div class="text-xs text-muted-foreground flex items-center gap-1.5">
            <ShieldCheck class="h-3.5 w-3.5 text-status-success" />
            {repo.protection}
          </div>
          <div class="text-xs text-muted-foreground">{repo.lastSync}</div>
          <button class="text-xs inline-flex items-center gap-1.5 text-muted-foreground hover:text-foreground transition-colors">
            Details
            <ExternalLink class="h-3.5 w-3.5" />
          </button>
        </div>
      {/each}
    </div>
  </div>
</PageShell>
