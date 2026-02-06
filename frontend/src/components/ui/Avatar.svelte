<script lang="ts">
  import { cn } from '@/lib/utils';

  export let src = '';
  export let alt = '';
  export let fallback = '';
  export let className = '';
  export let fallbackClass = '';

  let imageLoaded = false;

  $: if (!src) {
    imageLoaded = false;
  }
</script>

<div class={cn('relative flex h-10 w-10 shrink-0 overflow-hidden rounded-full', className)}>
  {#if src}
    <img
      src={src}
      alt={alt}
      class={cn('aspect-square h-full w-full', !imageLoaded && 'hidden')}
      on:load={() => (imageLoaded = true)}
      on:error={() => (imageLoaded = false)}
    />
  {/if}
  <div
    class={cn(
      'flex h-full w-full items-center justify-center rounded-full bg-muted',
      fallbackClass,
      imageLoaded && 'hidden'
    )}
  >
    {fallback}
  </div>
</div>
