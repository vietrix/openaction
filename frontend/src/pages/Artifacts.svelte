<script lang="ts">
  import { Package, Download, HardDrive, CalendarClock, Trash2 } from '@lucide/svelte';
  import PageShell from '@/components/layout/PageShell.svelte';

  const artifacts = [
    { name: 'forge-api-linux-amd64.tar.gz', pipeline: 'forge-api #1250', size: '128 MB', retention: '6 days', createdAt: '11 min ago' },
    { name: 'forge-web-static-bundle.zip', pipeline: 'forge-web #1248', size: '42 MB', retention: '13 days', createdAt: '48 min ago' },
    { name: 'forge-worker-container-metadata.json', pipeline: 'forge-worker #1247', size: '1.2 MB', retention: '29 days', createdAt: '1 hr ago' },
  ];
</script>

<PageShell title="Artifacts" description="Theo dõi artifact build, retention và tải xuống bản đóng gói">
  <div class="grid grid-cols-1 md:grid-cols-3 gap-3 mb-6">
    <div class="rounded-lg border border-border bg-card p-4">
      <div class="flex items-center gap-2 text-muted-foreground">
        <Package class="h-4 w-4" />
        <span class="text-xs uppercase tracking-wide">Artifacts</span>
      </div>
      <p class="text-2xl font-semibold mt-1">{artifacts.length}</p>
    </div>
    <div class="rounded-lg border border-border bg-card p-4">
      <div class="flex items-center gap-2 text-muted-foreground">
        <HardDrive class="h-4 w-4" />
        <span class="text-xs uppercase tracking-wide">Total Size</span>
      </div>
      <p class="text-2xl font-semibold mt-1">171.2 MB</p>
    </div>
    <div class="rounded-lg border border-border bg-card p-4">
      <div class="flex items-center gap-2 text-muted-foreground">
        <CalendarClock class="h-4 w-4" />
        <span class="text-xs uppercase tracking-wide">Default Retention</span>
      </div>
      <p class="text-2xl font-semibold mt-1">30 days</p>
    </div>
  </div>

  <div class="rounded-lg border border-border bg-card overflow-hidden">
    <div class="px-4 py-3 border-b border-border bg-muted/20">
      <h2 class="text-sm font-medium">Recent Artifacts</h2>
    </div>
    <div class="divide-y divide-border/50">
      {#each artifacts as artifact (artifact.name)}
        <div class="grid grid-cols-1 lg:grid-cols-[1.5fr_1fr_0.8fr_0.8fr_auto] gap-3 px-4 py-4">
          <div class="min-w-0">
            <p class="text-sm font-medium truncate">{artifact.name}</p>
            <p class="text-xs text-muted-foreground mt-0.5">{artifact.pipeline}</p>
          </div>
          <div class="text-xs text-muted-foreground">{artifact.size}</div>
          <div class="text-xs text-muted-foreground">{artifact.retention}</div>
          <div class="text-xs text-muted-foreground">{artifact.createdAt}</div>
          <div class="flex items-center gap-1">
            <button class="h-8 w-8 inline-flex items-center justify-center rounded-md hover:bg-muted text-muted-foreground hover:text-foreground transition-colors">
              <Download class="h-4 w-4" />
            </button>
            <button class="h-8 w-8 inline-flex items-center justify-center rounded-md hover:bg-muted text-muted-foreground hover:text-status-error transition-colors">
              <Trash2 class="h-4 w-4" />
            </button>
          </div>
        </div>
      {/each}
    </div>
  </div>
</PageShell>
