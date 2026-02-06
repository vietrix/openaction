<script lang="ts">
  import { onMount } from 'svelte';
  import { cn } from '@/lib/utils';

  export let align: 'start' | 'end' = 'end';
  export let className = '';

  let open = false;
  let root: HTMLDivElement | null = null;

  const close = () => {
    open = false;
  };

  const toggle = (event: MouseEvent) => {
    event.stopPropagation();
    open = !open;
  };

  const handleTriggerKeydown = (event: KeyboardEvent) => {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      open = !open;
    }
  };

  const handleMenuKeydown = (event: KeyboardEvent) => {
    if (event.key === 'Escape') {
      event.preventDefault();
      close();
    }
  };

  const handleDocumentClick = (event: MouseEvent) => {
    if (!open || !root) return;
    if (!root.contains(event.target as Node)) {
      close();
    }
  };

  const handleKeydown = (event: KeyboardEvent) => {
    if (event.key === 'Escape') {
      close();
    }
  };

  onMount(() => {
    document.addEventListener('mousedown', handleDocumentClick);
    document.addEventListener('keydown', handleKeydown);
    return () => {
      document.removeEventListener('mousedown', handleDocumentClick);
      document.removeEventListener('keydown', handleKeydown);
    };
  });
</script>

<div class="relative inline-flex" bind:this={root}>
  <div
    role="button"
    tabindex="0"
    on:click={toggle}
    on:keydown={handleTriggerKeydown}
    aria-haspopup="menu"
    aria-expanded={open}
  >
    <slot name="trigger" />
  </div>
  {#if open}
    <div
      class={cn(
        'absolute top-full mt-2 z-50 min-w-[8rem] overflow-hidden rounded-md border bg-popover p-1 text-popover-foreground shadow-md',
        align === 'end' ? 'right-0' : 'left-0',
        className
      )}
      role="menu"
      tabindex="-1"
      on:click={close}
      on:keydown={handleMenuKeydown}
    >
      <slot name="content" />
    </div>
  {/if}
</div>
