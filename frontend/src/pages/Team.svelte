<script lang="ts">
  import { Users, Shield, Clock3, Mail } from '@lucide/svelte';
  import PageShell from '@/components/layout/PageShell.svelte';
  import StatusBadge from '@/components/ui/StatusBadge.svelte';

  const members = [
    { name: 'Le Hung Quang Minh', role: 'Owner', email: 'minh@forge.dev', lastActive: 'now', status: 'success' as const },
    { name: 'Linh Tran', role: 'Backend', email: 'linh@forge.dev', lastActive: '5 min ago', status: 'success' as const },
    { name: 'Duc Nguyen', role: 'Infra', email: 'duc@forge.dev', lastActive: '17 min ago', status: 'warning' as const },
    { name: 'CI Bot', role: 'Automation', email: 'bot@forge.dev', lastActive: '1 min ago', status: 'running' as const },
  ];
</script>

<PageShell title="Team" description="Quản lý quyền truy cập và hoạt động của thành viên trong workspace">
  <div class="grid grid-cols-1 md:grid-cols-3 gap-3 mb-6">
    <div class="rounded-lg border border-border bg-card p-4">
      <div class="flex items-center gap-2 text-muted-foreground">
        <Users class="h-4 w-4" />
        <span class="text-xs uppercase tracking-wide">Members</span>
      </div>
      <p class="text-2xl font-semibold mt-1">{members.length}</p>
    </div>
    <div class="rounded-lg border border-border bg-card p-4">
      <div class="flex items-center gap-2 text-muted-foreground">
        <Shield class="h-4 w-4" />
        <span class="text-xs uppercase tracking-wide">Admins</span>
      </div>
      <p class="text-2xl font-semibold mt-1">1</p>
    </div>
    <div class="rounded-lg border border-border bg-card p-4">
      <div class="flex items-center gap-2 text-muted-foreground">
        <Clock3 class="h-4 w-4" />
        <span class="text-xs uppercase tracking-wide">Active Today</span>
      </div>
      <p class="text-2xl font-semibold mt-1">4</p>
    </div>
  </div>

  <div class="rounded-lg border border-border bg-card overflow-hidden">
    <div class="px-4 py-3 border-b border-border bg-muted/20 flex items-center justify-between">
      <h2 class="text-sm font-medium">Member Directory</h2>
      <button class="px-3 py-1.5 rounded-md text-xs bg-primary text-primary-foreground hover:bg-primary/90 transition-colors">
        Invite Member
      </button>
    </div>

    <div class="divide-y divide-border/50">
      {#each members as member (member.email)}
        <div class="grid grid-cols-1 md:grid-cols-[1.2fr_0.8fr_1fr_0.8fr_auto] gap-3 px-4 py-4">
          <div class="min-w-0">
            <p class="text-sm font-medium truncate">{member.name}</p>
            <p class="text-xs text-muted-foreground mt-0.5">{member.role}</p>
          </div>
          <div class="text-xs text-muted-foreground flex items-center gap-1.5">
            <Mail class="h-3.5 w-3.5" />
            <span class="truncate">{member.email}</span>
          </div>
          <div class="text-xs text-muted-foreground">{member.lastActive}</div>
          <StatusBadge status={member.status} />
          <button class="text-xs text-muted-foreground hover:text-foreground transition-colors">Manage</button>
        </div>
      {/each}
    </div>
  </div>
</PageShell>
