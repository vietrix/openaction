<script lang="ts">
  import { Bell, Shield, Database, TimerReset } from '@lucide/svelte';
  import PageShell from '@/components/layout/PageShell.svelte';
</script>

<PageShell title="Settings" description="Cấu hình workspace, chính sách bảo mật và thông báo hệ thống">
  <div class="space-y-4">
    <section class="rounded-lg border border-border bg-card p-5">
      <h2 class="text-sm font-medium mb-4 flex items-center gap-2">
        <Bell class="h-4 w-4 text-muted-foreground" />
        Notifications
      </h2>
      <div class="space-y-3 text-sm">
        <label class="flex items-center justify-between gap-3">
          <span class="text-muted-foreground">Notify when pipeline fails</span>
          <input type="checkbox" checked class="h-4 w-4 rounded border-border bg-background" />
        </label>
        <label class="flex items-center justify-between gap-3">
          <span class="text-muted-foreground">Notify when deploy succeeds</span>
          <input type="checkbox" checked class="h-4 w-4 rounded border-border bg-background" />
        </label>
      </div>
    </section>

    <section class="rounded-lg border border-border bg-card p-5">
      <h2 class="text-sm font-medium mb-4 flex items-center gap-2">
        <Shield class="h-4 w-4 text-muted-foreground" />
        Security
      </h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <label class="text-sm text-muted-foreground">Session timeout (minutes)
          <input type="number" value="60" class="mt-1 w-full rounded-md border border-border bg-background px-3 py-2 text-foreground" />
        </label>
        <label class="text-sm text-muted-foreground">Audit log retention (days)
          <input type="number" value="90" class="mt-1 w-full rounded-md border border-border bg-background px-3 py-2 text-foreground" />
        </label>
      </div>
    </section>

    <section class="rounded-lg border border-border bg-card p-5">
      <h2 class="text-sm font-medium mb-4 flex items-center gap-2">
        <Database class="h-4 w-4 text-muted-foreground" />
        Data & Retention
      </h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
        <label class="text-sm text-muted-foreground">Artifact retention (days)
          <input type="number" value="30" class="mt-1 w-full rounded-md border border-border bg-background px-3 py-2 text-foreground" />
        </label>
        <label class="text-sm text-muted-foreground">Pipeline history limit
          <input type="number" value="500" class="mt-1 w-full rounded-md border border-border bg-background px-3 py-2 text-foreground" />
        </label>
      </div>
    </section>

    <div class="flex items-center justify-end gap-2">
      <button class="inline-flex items-center gap-1.5 px-3 py-2 rounded-md border border-border text-sm text-foreground hover:bg-muted transition-colors">
        <TimerReset class="h-4 w-4" />
        Reset
      </button>
      <button class="px-4 py-2 rounded-md text-sm font-medium bg-primary text-primary-foreground hover:bg-primary/90 transition-colors">
        Save Changes
      </button>
    </div>
  </div>
</PageShell>
