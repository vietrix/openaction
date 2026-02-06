<script lang="ts">
  import { cn } from '@/lib/utils';

  export let content = '';
  export type Side = 'top' | 'right' | 'bottom' | 'left';
  export let side: Side = 'top';
  export let className = '';

  let open = false;

  const sideClasses: Record<Side, string> = {
    top: 'bottom-full mb-2 left-1/2 -translate-x-1/2',
    right: 'left-full ml-2 top-1/2 -translate-y-1/2',
    bottom: 'top-full mt-2 left-1/2 -translate-x-1/2',
    left: 'right-full mr-2 top-1/2 -translate-y-1/2',
  };
</script>

<span
  class="relative inline-flex"
  role="presentation"
  on:mouseenter={() => (open = true)}
  on:mouseleave={() => (open = false)}
  on:focusin={() => (open = true)}
  on:focusout={() => (open = false)}
>
  <slot />
  {#if open}
    <span
      role="tooltip"
      class={cn(
        'z-50 absolute whitespace-nowrap rounded-md border bg-popover px-3 py-1.5 text-sm text-popover-foreground shadow-md animate-in fade-in-0 zoom-in-95 pointer-events-none',
        sideClasses[side],
        className
      )}
    >
      {content}
    </span>
  {/if}
</span>
