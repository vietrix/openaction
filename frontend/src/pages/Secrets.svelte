<script lang="ts">
  import { KeyRound, ShieldCheck, EyeOff, Clock3 } from '@lucide/svelte';
  import PageShell from '@/components/layout/PageShell.svelte';

  const secrets = [
    { key: 'AWS_ACCESS_KEY_ID', scope: 'production', updatedBy: 'minh', updatedAt: '2 days ago' },
    { key: 'AWS_SECRET_ACCESS_KEY', scope: 'production', updatedBy: 'minh', updatedAt: '2 days ago' },
    { key: 'SENTRY_DSN', scope: 'workspace', updatedBy: 'linh', updatedAt: '7 days ago' },
    { key: 'DOCKERHUB_TOKEN', scope: 'staging', updatedBy: 'bot', updatedAt: '12 days ago' },
  ];
</script>

<PageShell title="Secrets" description="Lưu trữ thông tin nhạy cảm theo scope với audit trail rõ ràng">
  <div class="grid grid-cols-1 md:grid-cols-3 gap-3 mb-6">
    <div class="rounded-lg border border-border bg-card p-4">
      <div class="flex items-center gap-2 text-muted-foreground">
        <KeyRound class="h-4 w-4" />
        <span class="text-xs uppercase tracking-wide">Stored Secrets</span>
      </div>
      <p class="text-2xl font-semibold mt-1">{secrets.length}</p>
    </div>
    <div class="rounded-lg border border-border bg-card p-4">
      <div class="flex items-center gap-2 text-muted-foreground">
        <ShieldCheck class="h-4 w-4 text-status-success" />
        <span class="text-xs uppercase tracking-wide">Encrypted</span>
      </div>
      <p class="text-2xl font-semibold mt-1">AES-256</p>
    </div>
    <div class="rounded-lg border border-border bg-card p-4">
      <div class="flex items-center gap-2 text-muted-foreground">
        <Clock3 class="h-4 w-4" />
        <span class="text-xs uppercase tracking-wide">Last Rotation</span>
      </div>
      <p class="text-2xl font-semibold mt-1">2 days</p>
    </div>
  </div>

  <div class="rounded-lg border border-border bg-card overflow-hidden">
    <div class="px-4 py-3 border-b border-border bg-muted/20 flex items-center justify-between">
      <h2 class="text-sm font-medium">Secret Inventory</h2>
      <button class="px-3 py-1.5 rounded-md text-xs bg-primary text-primary-foreground hover:bg-primary/90 transition-colors">
        Add Secret
      </button>
    </div>

    <div class="divide-y divide-border/50">
      {#each secrets as secret (secret.key)}
        <div class="grid grid-cols-1 md:grid-cols-[1.2fr_0.8fr_0.8fr_0.8fr_auto] gap-3 px-4 py-4">
          <div class="min-w-0">
            <p class="text-sm font-medium truncate">{secret.key}</p>
            <p class="text-xs text-muted-foreground mt-0.5">••••••••••••••••</p>
          </div>
          <div class="text-xs text-muted-foreground">{secret.scope}</div>
          <div class="text-xs text-muted-foreground">{secret.updatedBy}</div>
          <div class="text-xs text-muted-foreground">{secret.updatedAt}</div>
          <button class="inline-flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors">
            <EyeOff class="h-3.5 w-3.5" />
            Rotate
          </button>
        </div>
      {/each}
    </div>
  </div>
</PageShell>
