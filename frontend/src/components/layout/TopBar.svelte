<script lang="ts">
  import { Search, ChevronRight, Bell, Command } from '@lucide/svelte';
  import Avatar from '@/components/ui/Avatar.svelte';
  import Button from '@/components/ui/Button.svelte';

  export let sidebarCollapsed: boolean;
</script>

<header
  class={`fixed top-0 z-30 flex h-14 items-center justify-between border-b border-border bg-background/80 backdrop-blur-sm transition-all duration-200 ${
    sidebarCollapsed ? 'left-14' : 'left-56'
  } right-0 px-4`}
>
  <nav class="flex items-center gap-1 text-sm">
    <span class="text-muted-foreground">Dashboard</span>
    <ChevronRight class="h-3.5 w-3.5 text-muted-foreground/50" />
    <span class="text-foreground font-medium">Pipelines</span>
    <ChevronRight class="h-3.5 w-3.5 text-muted-foreground/50" />
    <span class="text-foreground font-medium">forge-api</span>
  </nav>

  <div class="flex items-center gap-3">
    <Button
      variant="outline"
      className="btn-glow h-8 gap-2 border-border bg-secondary/50 px-3 text-sm text-muted-foreground hover:bg-secondary hover:text-foreground"
    >
      <Search class="h-3.5 w-3.5" />
      <span class="hidden sm:inline">Search...</span>
      <div class="hidden items-center gap-0.5 sm:flex">
        <kbd class="kbd">
          <Command class="h-2.5 w-2.5" />
        </kbd>
        <kbd class="kbd">K</kbd>
      </div>
    </Button>

    <Button
      variant="ghost"
      size="icon"
      className="relative h-8 w-8 text-muted-foreground hover:text-foreground"
    >
      <Bell class="h-4 w-4" />
      <span class="absolute right-1.5 top-1.5 h-1.5 w-1.5 rounded-full bg-status-error"></span>
    </Button>

    <Avatar
      className="h-7 w-7 border border-border"
      src="https://github.com/shadcn.png"
      alt="User"
      fallback="JD"
      fallbackClass="bg-muted text-xs font-medium"
    />
  </div>
</header>
