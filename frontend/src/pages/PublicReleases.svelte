<script lang="ts">
  import { Download, Home, Package, Layers } from '@lucide/svelte';
  import { marked } from 'marked';
  import { cn } from '@/lib/utils';
  import { link } from '@/lib/router';

  type Asset = { name: string; size: string };
  type Patch = { id: string; name: string; assets: Asset[] };
  type Build = { id: string; name: string; patches: Patch[] };
  type Version = { id: string; name: string; builds: Build[] };

  const buildTree: Version[] = [
    {
      id: 'v0.0.3',
      name: 'v0.0.3',
      builds: [
        {
          id: 'dev',
          name: 'dev',
          patches: [
            {
              id: 'a',
              name: 'a',
              assets: [
                { name: 'forge-ci-dev-a-linux-amd64.tar.gz', size: '128 MB' },
                { name: 'forge-ci-dev-a-darwin-arm64.tar.gz', size: '122 MB' },
                { name: 'forge-ci-dev-a-windows-x64.zip', size: '130 MB' },
              ],
            },
            {
              id: 'b',
              name: 'b',
              assets: [
                { name: 'forge-ci-dev-b-linux-amd64.tar.gz', size: '129 MB' },
                { name: 'forge-ci-dev-b-darwin-arm64.tar.gz', size: '123 MB' },
                { name: 'forge-ci-dev-b-windows-x64.zip', size: '131 MB' },
              ],
            },
          ],
        },
        {
          id: 'prod',
          name: 'prod',
          patches: [
            {
              id: 'v0.0.3b',
              name: 'v0.0.3b',
              assets: [
                { name: 'forge-ci-prod-v0.0.3b-linux-amd64.tar.gz', size: '130 MB' },
                { name: 'forge-ci-prod-v0.0.3b-darwin-arm64.tar.gz', size: '124 MB' },
                { name: 'forge-ci-prod-v0.0.3b-windows-x64.zip', size: '132 MB' },
              ],
            },
          ],
        },
      ],
    },
    {
      id: 'bixie',
      name: 'bixie',
      builds: [
        {
          id: 'dev',
          name: 'dev',
          patches: [
            {
              id: 'a',
              name: 'a',
              assets: [
                { name: 'forge-ci-bixie-dev-a-linux-amd64.tar.gz', size: '134 MB' },
                { name: 'forge-ci-bixie-dev-a-darwin-arm64.tar.gz', size: '129 MB' },
                { name: 'forge-ci-bixie-dev-a-windows-x64.zip', size: '136 MB' },
              ],
            },
          ],
        },
        {
          id: 'prod',
          name: 'prod',
          patches: [
            {
              id: 'a',
              name: 'a',
              assets: [
                { name: 'forge-ci-bixie-prod-a-linux-amd64.tar.gz', size: '140 MB' },
                { name: 'forge-ci-bixie-prod-a-darwin-arm64.tar.gz', size: '135 MB' },
                { name: 'forge-ci-bixie-prod-a-windows-x64.zip', size: '142 MB' },
              ],
            },
          ],
        },
      ],
    },
  ];

  const updateModules = import.meta.glob('/docs/updates/**/*.md', {
    query: '?raw',
    import: 'default',
    eager: true,
  });
  const updateMap = updateModules as Record<string, string>;

  const getUpdateKey = (versionId: string, buildId: string, patchId: string) =>
    `/docs/updates/${versionId}/${buildId}/${patchId}.md`;

  const getLatestSelection = () => {
    const version = buildTree[buildTree.length - 1];
    const build = version?.builds[version.builds.length - 1];
    const patch = build?.patches[build.patches.length - 1];
    return {
      versionId: version?.id ?? '',
      buildId: build?.id ?? '',
      patchId: patch?.id ?? '',
    };
  };

  let selected = getLatestSelection();

  const selectPatch = (versionId: string, buildId: string, patchId: string) => {
    selected = { versionId, buildId, patchId };
  };

  $: selectedVersion = buildTree.find((version) => version.id === selected.versionId);
  $: selectedBuild = selectedVersion?.builds.find((build) => build.id === selected.buildId);
  $: selectedPatch = selectedBuild?.patches.find((patch) => patch.id === selected.patchId);
  $: updateKey = selectedPatch ? getUpdateKey(selected.versionId, selected.buildId, selected.patchId) : '';
  $: updateMarkdown = updateKey && updateMap[updateKey] ? updateMap[updateKey] : '';
  $: updateHtml = updateMarkdown ? marked.parse(updateMarkdown) : '';
</script>

<div class="dark min-h-screen bg-background text-foreground">
  <div class="min-h-screen lg:flex">
    <aside class="w-full lg:w-80 border-b border-border bg-card/40 lg:border-b-0 lg:border-r">
      <div class="px-5 py-6 border-b border-border">
        <div class="flex items-center gap-2 text-sm font-semibold text-foreground">
          <div class="h-9 w-9 rounded-lg bg-primary/15 text-primary flex items-center justify-center font-semibold">
            CI
          </div>
          Public Releases
        </div>
        <a
          href="/"
          use:link
          class="mt-4 inline-flex items-center gap-2 text-xs text-muted-foreground hover:text-foreground transition-colors"
        >
          <Home class="h-3.5 w-3.5" />
          Home
        </a>
      </div>

      <div class="px-4 py-4 overflow-y-auto max-h-[calc(100vh-120px)]">
        <div class="text-[11px] uppercase tracking-widest text-muted-foreground">Build Tree</div>
        <div class="mt-4 tree">
          {#each buildTree as version, versionIndex (version.id)}
            <div class="tree-version">
              <div class="tree-node tree-root">
                <Package class="h-4 w-4 text-primary" />
                <span class="text-sm font-semibold text-foreground">{version.name}</span>
              </div>
              <div class="tree-children">
                {#each version.builds as build, buildIndex (build.id)}
                  <div class={cn('tree-row', buildIndex === version.builds.length - 1 && 'tree-row-last')}>
                    <div class="tree-node tree-branch">
                      <Layers class="h-3.5 w-3.5 text-muted-foreground" />
                      <span class="text-xs uppercase tracking-wide text-muted-foreground">{build.name}</span>
                    </div>
                    <div class="tree-children">
                      {#each build.patches as patch, patchIndex (patch.id)}
                        <div class={cn('tree-row', patchIndex === build.patches.length - 1 && 'tree-row-last')}>
                          <button
                            type="button"
                            class={cn(
                              'tree-node tree-leaf w-full text-left text-xs transition-colors',
                              selected.versionId === version.id &&
                                selected.buildId === build.id &&
                                selected.patchId === patch.id
                                ? 'tree-selected'
                                : 'tree-unselected'
                            )}
                            on:click={() => selectPatch(version.id, build.id, patch.id)}
                          >
                            <span class="font-medium">{patch.name}</span>
                            <span class="text-[10px] text-muted-foreground">
                              {patch.assets.length} assets
                            </span>
                          </button>
                        </div>
                      {/each}
                    </div>
                  </div>
                {/each}
              </div>
            </div>
          {/each}
        </div>
      </div>
    </aside>

    <main class="flex-1 px-6 py-8">
      <div class="mx-auto max-w-5xl">
        <div class="flex flex-wrap items-center gap-3 mb-6">
          <span class="inline-flex items-center gap-2 rounded-full bg-muted px-3 py-1 text-xs text-muted-foreground">
            Public Release Notes
          </span>
          {#if selectedPatch}
            <span class="inline-flex items-center gap-2 rounded-full border border-border px-3 py-1 text-xs text-foreground">
              {selectedVersion?.name} / {selectedBuild?.name} / {selectedPatch.name}
            </span>
          {/if}
        </div>

        <section class="rounded-xl border border-border bg-card p-6">
          <div class="flex items-center justify-between mb-4">
            <h1 class="text-lg font-semibold text-foreground">Update Notes</h1>
            {#if selectedPatch}
              <span class="text-xs text-muted-foreground">Patch {selectedPatch.name}</span>
            {/if}
          </div>
          {#if updateHtml}
            <div class="update-content">{@html updateHtml}</div>
          {:else}
            <div class="text-sm text-muted-foreground">
              No update content found for this patch.
            </div>
          {/if}
        </section>

        <section class="mt-6 rounded-xl border border-border bg-card overflow-hidden">
          <div class="px-5 py-4 border-b border-border bg-muted/20 flex items-center justify-between">
            <h2 class="text-sm font-medium text-foreground">Artifacts</h2>
            <span class="text-xs text-muted-foreground">
              {selectedPatch ? selectedPatch.assets.length : 0} files
            </span>
          </div>
          <div class="divide-y divide-border/50">
            {#if selectedPatch}
              {#each selectedPatch.assets as asset (asset.name)}
                <div class="flex items-center justify-between gap-4 px-5 py-4">
                  <div class="min-w-0">
                    <p class="text-sm font-medium text-foreground truncate">{asset.name}</p>
                    <p class="text-xs text-muted-foreground mt-0.5">{asset.size}</p>
                  </div>
                  <button class="inline-flex items-center gap-2 px-3 py-1.5 rounded-md text-xs bg-secondary text-secondary-foreground hover:bg-secondary/80 transition-colors">
                    <Download class="h-3.5 w-3.5" />
                    Download
                  </button>
                </div>
              {/each}
            {:else}
              <div class="px-5 py-6 text-sm text-muted-foreground">
                Select a patch from the left to view artifacts.
              </div>
            {/if}
          </div>
        </section>
      </div>
    </main>
  </div>
</div>

<style>
  .update-content {
    color: hsl(var(--foreground));
    font-size: 0.95rem;
    line-height: 1.65;
  }
  .update-content :global(h1),
  .update-content :global(h2),
  .update-content :global(h3) {
    font-weight: 600;
    margin: 1.2rem 0 0.6rem;
  }
  .update-content :global(h1) {
    font-size: 1.25rem;
  }
  .update-content :global(h2) {
    font-size: 1.1rem;
  }
  .update-content :global(p) {
    margin: 0.6rem 0;
    color: hsl(var(--foreground));
  }
  .update-content :global(ul),
  .update-content :global(ol) {
    margin: 0.5rem 0 0.8rem 1.2rem;
  }
  .update-content :global(li) {
    margin: 0.2rem 0;
  }
  .update-content :global(code) {
    background: hsl(var(--muted));
    padding: 0.15rem 0.35rem;
    border-radius: 0.3rem;
    font-size: 0.85rem;
  }
  .update-content :global(pre) {
    background: hsl(var(--muted));
    padding: 0.75rem;
    border-radius: 0.5rem;
    overflow-x: auto;
  }
  .update-content :global(a) {
    color: hsl(var(--primary));
    text-decoration: underline;
  }

  .tree {
    display: grid;
    gap: 1.1rem;
  }
  .tree-version {
    display: grid;
    gap: 0.5rem;
  }
  .tree-node {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    position: relative;
  }
  .tree-root {
    padding-left: 0.25rem;
  }
  .tree-children {
    display: grid;
    gap: 0.35rem;
    margin-left: 0.45rem;
    padding-left: 1.25rem;
    position: relative;
  }
  .tree-children::before {
    content: '';
    position: absolute;
    left: 0.22rem;
    top: 0.092rem;
    bottom: -0.02rem;
    width: 1px;
    background: hsl(var(--border));
  }
  .tree-row {
    position: relative;
  }
  .tree-row-last::after {
    content: '';
    position: absolute;
    left: -0.75rem;
    top: 50%;
    bottom: -0.55rem;
    width: 1px;
    background: hsl(var(--background));
  }
  .tree-branch,
  .tree-leaf {
    padding-left: 0.25rem;
  }
  .tree-branch::before,
  .tree-leaf::before {
    content: '';
    position: absolute;
    left: -0.75rem;
    top: 50%;
    width: 0.75rem;
    height: 1px;
    background: hsl(var(--border));
  }
  .tree-leaf {
    justify-content: space-between;
    border-radius: 0.45rem;
    padding: 0.45rem 0.6rem 0.45rem 0.35rem;
  }
  .tree-unselected {
    color: hsl(var(--foreground));
    background: transparent;
  }
  .tree-unselected:hover {
    background: hsl(var(--muted));
  }
  .tree-selected {
    color: hsl(var(--primary));
    background: hsl(var(--primary) / 0.12);
    border: 1px solid hsl(var(--primary) / 0.3);
  }
</style>
