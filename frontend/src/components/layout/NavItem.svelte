<script lang="ts">
  import { link, location } from '@/lib/router';
  import Tooltip from '@/components/ui/Tooltip.svelte';
  import { cn } from '@/lib/utils';

  export let icon: any;
  export let label: string;
  export let href: string;
  export let collapsed: boolean;

  $: pathname = $location.pathname;
  $: isActive = href === '/' ? pathname === '/' : pathname.startsWith(href);
</script>

{#if collapsed}
  <Tooltip content={label} side="right" className="text-xs">
    <a
      href={href}
      use:link
      class={cn(
        'sidebar-icon flex items-center gap-3 rounded-md px-2.5 py-2 text-sm',
        isActive
          ? 'bg-sidebar-accent text-sidebar-accent-foreground font-medium'
          : 'text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground'
      )}
    >
      <svelte:component this={icon} class="h-4 w-4 shrink-0" />
    </a>
  </Tooltip>
{:else}
  <a
    href={href}
    use:link
    class={cn(
      'sidebar-icon flex items-center gap-3 rounded-md px-2.5 py-2 text-sm',
      isActive
        ? 'bg-sidebar-accent text-sidebar-accent-foreground font-medium'
        : 'text-sidebar-foreground hover:bg-sidebar-accent hover:text-sidebar-accent-foreground'
    )}
  >
    <svelte:component this={icon} class="h-4 w-4 shrink-0" />
    <span>{label}</span>
  </a>
{/if}
