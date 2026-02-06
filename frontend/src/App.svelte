<script lang="ts">
  import Index from '@/pages/Index.svelte';
  import Pipelines from '@/pages/Pipelines.svelte';
  import PipelineDetails from '@/pages/PipelineDetails.svelte';
  import Repositories from '@/pages/Repositories.svelte';
  import Artifacts from '@/pages/Artifacts.svelte';
  import Runners from '@/pages/Runners.svelte';
  import Releases from '@/pages/Releases.svelte';
  import PublicReleases from '@/pages/PublicReleases.svelte';
  import Team from '@/pages/Team.svelte';
  import Secrets from '@/pages/Secrets.svelte';
  import Docs from '@/pages/Docs.svelte';
  import Settings from '@/pages/Settings.svelte';
  import NotFound from '@/pages/NotFound.svelte';
  import Toaster from '@/components/ui/Toaster.svelte';
  import Sonner from '@/components/ui/Sonner.svelte';
  import { location } from '@/lib/router';

  let currentComponent:
    | typeof Index
    | typeof Pipelines
    | typeof PipelineDetails
    | typeof Repositories
    | typeof Artifacts
    | typeof Runners
    | typeof Releases
    | typeof PublicReleases
    | typeof Team
    | typeof Secrets
    | typeof Docs
    | typeof Settings
    | typeof NotFound = Index;

  $: path = $location.pathname;

  $: {
    const pipelineMatch = path.match(/^\/pipelines\/([^/]+)\/([^/]+)$/);

    if (path === '/') {
      currentComponent = Index;
    } else if (path === '/pipelines') {
      currentComponent = Pipelines;
    } else if (path === '/runners') {
      currentComponent = Runners;
    } else if (path === '/releases') {
      currentComponent = Releases;
    } else if (path === '/public/releases') {
      currentComponent = PublicReleases;
    } else if (path === '/repos') {
      currentComponent = Repositories;
    } else if (path === '/artifacts') {
      currentComponent = Artifacts;
    } else if (path === '/team') {
      currentComponent = Team;
    } else if (path === '/secrets') {
      currentComponent = Secrets;
    } else if (path === '/docs') {
      currentComponent = Docs;
    } else if (path === '/settings') {
      currentComponent = Settings;
    } else if (pipelineMatch) {
      currentComponent = PipelineDetails;
    } else {
      currentComponent = NotFound;
    }
  }
</script>

<Toaster />
<Sonner />

<svelte:component this={currentComponent} />
