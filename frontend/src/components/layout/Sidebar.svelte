<script lang="ts">
  import {
    LayoutDashboard,
    GitBranch,
    Box,
    Settings,
    Users,
    Key,
    BookOpen,
    ChevronLeft,
    ChevronRight,
    Workflow,
    HardDrive,
    Rocket,
  } from '@lucide/svelte';
  import { link } from '@/lib/router';
  import { cn } from '@/lib/utils';
  import NavItem from './NavItem.svelte';

  export let collapsed: boolean;
  export let onToggle: () => void;

  const navItems = [
    { icon: LayoutDashboard, label: 'Dashboard', href: '/' },
    { icon: Workflow, label: 'Pipelines', href: '/pipelines' },
    { icon: GitBranch, label: 'Repositories', href: '/repos' },
    { icon: Box, label: 'Artifacts', href: '/artifacts' },
    { icon: Rocket, label: 'Releases', href: '/releases' },
    { icon: HardDrive, label: 'Runners', href: '/runners' },
  ];

  const bottomItems = [
    { icon: Users, label: 'Team', href: '/team' },
    { icon: Key, label: 'Secrets', href: '/secrets' },
    { icon: BookOpen, label: 'Docs', href: '/docs' },
    { icon: Settings, label: 'Settings', href: '/settings' },
  ];
</script>

<aside
  class={cn(
    'fixed left-0 top-0 z-40 h-screen border-r border-sidebar-border bg-sidebar transition-all duration-200 ease-out',
    collapsed ? 'w-14' : 'w-56'
  )}
>
  <div class="flex h-14 items-center justify-between border-b border-sidebar-border px-3">
    <a href="/" use:link class="flex items-center gap-2 overflow-hidden">
      <div class="flex h-8 w-8 shrink-0 items-center justify-center rounded-md bg-foreground">
        <span class="text-sm font-bold text-background">F</span>
      </div>
      {#if !collapsed}
        <span class="text-sm font-semibold text-foreground">Forge CI</span>
      {/if}
    </a>
  </div>

  <nav class="flex h-[calc(100vh-3.5rem)] flex-col justify-between p-2">
    <div class="space-y-1">
      {#each navItems as item}
        <NavItem
          icon={item.icon}
          label={item.label}
          href={item.href}
          {collapsed}
        />
      {/each}
    </div>

    <div class="space-y-1">
      {#each bottomItems as item}
        <NavItem
          icon={item.icon}
          label={item.label}
          href={item.href}
          {collapsed}
        />
      {/each}

      <button
        on:click={onToggle}
        class="flex w-full items-center gap-3 rounded-md px-2.5 py-2 text-sm text-sidebar-foreground transition-colors duration-150 hover:bg-sidebar-accent hover:text-sidebar-accent-foreground"
      >
        {#if collapsed}
          <ChevronRight class="h-4 w-4 shrink-0" />
        {:else}
          <ChevronLeft class="h-4 w-4 shrink-0" />
          <span>Collapse</span>
        {/if}
      </button>
    </div>
  </nav>
</aside>
