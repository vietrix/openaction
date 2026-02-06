<script lang="ts">
  import { BookOpen, FileText, Rocket, Terminal, Shield } from '@lucide/svelte';
  import PageShell from '@/components/layout/PageShell.svelte';

  const docs = [
    { title: 'Getting Started', description: 'Thiết lập project, kết nối repo và chạy pipeline đầu tiên.', icon: Rocket },
    { title: 'Pipeline DSL', description: 'Định nghĩa step, dependency graph và strategy cho runner.', icon: Terminal },
    { title: 'Release Process', description: 'Quy trình promote staging -> production và rollback.', icon: FileText },
    { title: 'Security Baseline', description: 'Chuẩn bảo mật cho secrets, artifacts và access policy.', icon: Shield },
  ];
</script>

<PageShell title="Docs" description="Tài liệu nội bộ cho CI/CD workflow, bảo mật và vận hành">
  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    {#each docs as item (item.title)}
      <article class="rounded-lg border border-border bg-card p-5 hover:border-foreground/20 transition-colors">
        <div class="flex items-start gap-3">
          <div class="h-9 w-9 rounded-md bg-muted flex items-center justify-center">
            <svelte:component this={item.icon} class="h-4.5 w-4.5 text-muted-foreground" />
          </div>
          <div class="min-w-0">
            <h2 class="text-sm font-medium text-foreground">{item.title}</h2>
            <p class="text-sm text-muted-foreground mt-1">{item.description}</p>
            <button class="mt-3 inline-flex items-center gap-1.5 text-xs text-muted-foreground hover:text-foreground transition-colors">
              <BookOpen class="h-3.5 w-3.5" />
              Open document
            </button>
          </div>
        </div>
      </article>
    {/each}
  </div>
</PageShell>
