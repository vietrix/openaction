<script lang="ts">
  import { cn } from '@/lib/utils';
  import Sidebar from '@/components/layout/Sidebar.svelte';
  import TopBar from '@/components/layout/TopBar.svelte';

  export let title: string;
  export let description: string;

  let sidebarCollapsed = false;
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
        <h1 class="text-xl font-semibold text-foreground">{title}</h1>
        <p class="text-sm text-muted-foreground mt-1">{description}</p>
      </div>

      <slot />
    </div>
  </main>
</div>
