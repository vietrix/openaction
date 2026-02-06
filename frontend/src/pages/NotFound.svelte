<script lang="ts">
  import { location, link } from '@/lib/router';
  $: console.error('404 Error: User attempted to access non-existent route:', $location.pathname);
</script>

<div class="dark flex min-h-screen items-center justify-center bg-muted text-foreground">
  <div class="text-center">
    <h1 class="mb-4 text-4xl font-bold">404</h1>
    <p class="mb-4 text-xl text-muted-foreground">Oops! Page not found</p>
    <a href="/" use:link class="text-primary underline hover:text-primary/90">Return to Home</a>
  </div>
</div>
