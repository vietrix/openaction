<script lang="ts">
  import { Download, Package, Rocket, GitCommit, CalendarClock, ShieldCheck } from '@lucide/svelte';
  import PageShell from '@/components/layout/PageShell.svelte';

  const releases = [
    {
      version: 'v2.4.1',
      project: 'forge-api',
      commit: 'abc1234',
      channel: 'stable',
      createdAt: '2024-01-16',
      notes: 'Improved rate limiting + faster startup.',
      assets: [
        { name: 'forge-api-linux-amd64.tar.gz', size: '128 MB' },
        { name: 'forge-api-darwin-arm64.tar.gz', size: '122 MB' },
      ],
    },
    {
      version: 'v2.4.0',
      project: 'forge-web',
      commit: 'def5678',
      channel: 'stable',
      createdAt: '2024-01-12',
      notes: 'New runner dashboard and stability fixes.',
      assets: [
        { name: 'forge-web-static-bundle.zip', size: '42 MB' },
      ],
    },
    {
      version: 'v2.5.0-rc.1',
      project: 'forge-worker',
      commit: 'ghi9012',
      channel: 'rc',
      createdAt: '2024-01-10',
      notes: 'Preview: improved caching strategy.',
      assets: [
        { name: 'forge-worker-linux-amd64.tar.gz', size: '98 MB' },
      ],
    },
  ];
</script>

<PageShell title="Releases" description="Khám phá và tải các bản build đã phát hành sau khi CI hoàn tất">
  <div class="grid grid-cols-1 md:grid-cols-3 gap-3 mb-6">
    <div class="rounded-lg border border-border bg-card p-4">
      <div class="flex items-center gap-2 text-muted-foreground">
        <Rocket class="h-4 w-4" />
        <span class="text-xs uppercase tracking-wide">Releases</span>
      </div>
      <p class="text-2xl font-semibold mt-1">{releases.length}</p>
    </div>
    <div class="rounded-lg border border-border bg-card p-4">
      <div class="flex items-center gap-2 text-muted-foreground">
        <Package class="h-4 w-4" />
        <span class="text-xs uppercase tracking-wide">Artifacts</span>
      </div>
      <p class="text-2xl font-semibold mt-1">4</p>
    </div>
    <div class="rounded-lg border border-border bg-card p-4">
      <div class="flex items-center gap-2 text-muted-foreground">
        <ShieldCheck class="h-4 w-4 text-status-success" />
        <span class="text-xs uppercase tracking-wide">Verified</span>
      </div>
      <p class="text-2xl font-semibold mt-1">100%</p>
    </div>
  </div>

  <div class="rounded-lg border border-border bg-card overflow-hidden">
    <div class="px-4 py-3 border-b border-border bg-muted/20 flex items-center justify-between">
      <h2 class="text-sm font-medium">Release History</h2>
      <button class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs border border-border text-muted-foreground hover:text-foreground hover:bg-muted transition-colors">
        <CalendarClock class="h-3.5 w-3.5" />
        Filter
      </button>
    </div>

    <div class="divide-y divide-border/50">
      {#each releases as release (release.version)}
        <div class="px-4 py-4">
          <div class="flex flex-wrap items-center gap-3">
            <span class="text-sm font-semibold text-foreground">{release.project}</span>
            <span class="text-xs text-muted-foreground">{release.version}</span>
            <span class="text-[10px] uppercase tracking-wide px-2 py-0.5 rounded-full bg-muted text-muted-foreground">
              {release.channel}
            </span>
            <span class="text-xs text-muted-foreground flex items-center gap-1.5">
              <GitCommit class="h-3.5 w-3.5" />
              {release.commit}
            </span>
            <span class="text-xs text-muted-foreground">{release.createdAt}</span>
          </div>
          <p class="text-sm text-muted-foreground mt-2">{release.notes}</p>

          <div class="mt-4 grid grid-cols-1 md:grid-cols-2 gap-3">
            {#each release.assets as asset (asset.name)}
              <div class="flex items-center justify-between rounded-md border border-border bg-muted/20 px-3 py-2">
                <div class="min-w-0">
                  <p class="text-sm font-medium truncate">{asset.name}</p>
                  <p class="text-xs text-muted-foreground mt-0.5">{asset.size}</p>
                </div>
                <button class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-md text-xs bg-primary text-primary-foreground hover:bg-primary/90 transition-colors">
                  <Download class="h-3.5 w-3.5" />
                  Download
                </button>
              </div>
            {/each}
          </div>
        </div>
      {/each}
    </div>
  </div>
</PageShell>
