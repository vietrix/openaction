<script lang="ts">
  import { onMount } from 'svelte';
  import { cn } from '@/lib/utils';
  import { GripVertical } from '@lucide/svelte';

  export let initial = 45;
  export let minLeft = 30;
  export let minRight = 35;
  export let withHandle = true;
  export let className = '';

  let container: HTMLDivElement | null = null;
  let leftPercent = initial;
  let dragging = false;

  const clamp = (value: number) => Math.min(100 - minRight, Math.max(minLeft, value));

  const onPointerMove = (event: PointerEvent) => {
    if (!dragging || !container) return;
    const rect = container.getBoundingClientRect();
    const next = ((event.clientX - rect.left) / rect.width) * 100;
    leftPercent = clamp(next);
  };

  const onPointerUp = () => {
    dragging = false;
    document.body.style.cursor = '';
    document.body.style.userSelect = '';
    window.removeEventListener('pointermove', onPointerMove);
    window.removeEventListener('pointerup', onPointerUp);
  };

  const onPointerDown = (event: PointerEvent) => {
    event.preventDefault();
    dragging = true;
    document.body.style.cursor = 'col-resize';
    document.body.style.userSelect = 'none';
    window.addEventListener('pointermove', onPointerMove);
    window.addEventListener('pointerup', onPointerUp);
  };

  onMount(() => {
    leftPercent = clamp(initial);
  });
</script>

<div
  class={cn('flex h-full w-full', className)}
  bind:this={container}
>
  <div class="min-w-0 flex-1" style={`flex-basis: ${leftPercent}%`}>
    <slot name="left" />
  </div>

  <div
    class="relative flex w-px items-center justify-center bg-border after:absolute after:inset-y-0 after:left-1/2 after:w-1 after:-translate-x-1/2 focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring focus-visible:ring-offset-1"
    on:pointerdown={onPointerDown}
    role="separator"
    aria-orientation="vertical"
  >
    {#if withHandle}
      <div class="z-10 flex h-4 w-3 items-center justify-center rounded-sm border bg-border">
        <GripVertical class="h-2.5 w-2.5" />
      </div>
    {/if}
  </div>

  <div class="min-w-0 flex-1" style={`flex-basis: ${100 - leftPercent}%`}>
    <slot name="right" />
  </div>
</div>
