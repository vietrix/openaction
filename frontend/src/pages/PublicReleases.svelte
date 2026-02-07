<script lang="ts">
  import { onMount } from 'svelte';
  import { Download, Home, Package, Layers } from '@lucide/svelte';
  import { marked } from 'marked';
  import { cn } from '@/lib/utils';
  import { link } from '@/lib/router';
  import { api } from '@/lib/api';

  type Release = {
    id: string;
    project_id: string;
    version: string;
    build: string;
    patch: string;
    created_at: number;
    update_path: string;
  };

  type Artifact = {
    id: string;
    filename: string;
    size_bytes: number;
  };

  type Patch = { id: string; name: string; releaseId: string; assetCount: number };
  type Build = { id: string; name: string; patches: Patch[] };
  type Version = { id: string; name: string; builds: Build[] };

  let buildTree: Version[] = [];
  let releases: Release[] = [];
  let selected = { versionId: '', buildId: '', patchId: '', releaseId: '' };
  let selectedUpdate = '';
  let selectedArtifacts: Artifact[] = [];
  let loading = true;
  let error = '';

  const formatBytes = (value: number) => {
    if (!value) return '0 B';
    if (value < 1024) return `${value} B`;
    const kb = value / 1024;
    if (kb < 1024) return `${kb.toFixed(1)} KB`;
    const mb = kb / 1024;
    if (mb < 1024) return `${mb.toFixed(1)} MB`;
    const gb = mb / 1024;
    return `${gb.toFixed(1)} GB`;
  };

  const buildTreeFromReleases = (items: Release[], artifactCounts: Record<string, number>) => {
    const versionMap = new Map<string, Version>();

    items.forEach((release) => {
      if (!versionMap.has(release.version)) {
        versionMap.set(release.version, { id: release.version, name: release.version, builds: [] });
      }
      const version = versionMap.get(release.version)!;
      let build = version.builds.find((b) => b.id === release.build);
      if (!build) {
        build = { id: release.build, name: release.build, patches: [] };
        version.builds = [...version.builds, build];
      }
      build.patches = [
        ...build.patches,
        {
          id: release.patch,
          name: release.patch,
          releaseId: release.id,
          assetCount: artifactCounts[release.id] ?? 0,
        },
      ];
    });

    return Array.from(versionMap.values()).map((version) => ({
      ...version,
      builds: [...version.builds]
        .map((build) => ({
          ...build,
          patches: [...build.patches].sort((a, b) => a.name.localeCompare(b.name)),
        }))
        .sort((a, b) => a.name.localeCompare(b.name)),
    }));
  };

  const selectPatch = async (versionId: string, buildId: string, patch: Patch) => {
    selected = { versionId, buildId, patchId: patch.id, releaseId: patch.releaseId };
    await loadReleaseDetails(patch.releaseId);
  };

  const loadReleaseDetails = async (releaseId: string) => {
    try {
      const [release, artifacts] = await Promise.all([
        api.getPublicRelease(releaseId),
        api.getPublicArtifacts(releaseId),
      ]);
      selectedUpdate = release.update_md || '';
      selectedArtifacts = artifacts.map((item) => ({
        id: item.id,
        filename: item.filename,
        size_bytes: item.size_bytes,
      }));
    } catch (err) {
      selectedUpdate = '';
      selectedArtifacts = [];
      error = err instanceof Error ? err.message : 'Không thể tải release';
    }
  };

  onMount(async () => {
    loading = true;
    error = '';
    try {
      releases = await api.getPublicReleases();
      const artifactsByRelease = await Promise.all(
        releases.map(async (release) => {
          try {
            const artifacts = await api.getPublicArtifacts(release.id);
            return { id: release.id, count: artifacts.length };
          } catch {
            return { id: release.id, count: 0 };
          }
        })
      );
      const counts = artifactsByRelease.reduce<Record<string, number>>((acc, item) => {
        acc[item.id] = item.count;
        return acc;
      }, {});
      buildTree = buildTreeFromReleases(releases, counts);
      const latest = [...releases].sort((a, b) => a.created_at - b.created_at).at(-1);
      if (latest) {
        selected = {
          versionId: latest.version,
          buildId: latest.build,
          patchId: latest.patch,
          releaseId: latest.id,
        };
        await loadReleaseDetails(latest.id);
      }
    } catch (err) {
      error = err instanceof Error ? err.message : 'Không thể tải danh sách release';
    } finally {
      loading = false;
    }
  });

  $: updateHtml = selectedUpdate
    ? marked.parse(selectedUpdate)
    : '<p>Chọn một bản vá để xem nội dung cập nhật.</p>';
</script>

<div class="dark min-h-screen bg-background text-foreground">
  <div class="flex min-h-screen">
    <aside class="w-72 border-r border-border bg-card p-5">
      <div class="flex items-center gap-3 text-primary">
        <Package class="h-5 w-5" />
        <span class="text-sm font-semibold">OpenAction</span>
      </div>
      <nav class="mt-6 space-y-1 text-sm">
        <a
          href="/"
          use:link
          class="flex items-center gap-2 rounded-md px-2.5 py-2 text-muted-foreground hover:text-foreground hover:bg-muted/30"
        >
          <Home class="h-4 w-4" />
          Trang chính
        </a>
        <div class="mt-4 text-xs uppercase tracking-wide text-muted-foreground">Build Tree</div>
        <div class="mt-4 tree">
          {#if loading}
            <div class="text-xs text-muted-foreground">Đang tải dữ liệu...</div>
          {:else if error}
            <div class="text-xs text-status-error">{error}</div>
          {:else}
            {#each buildTree as version (version.id)}
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
                                selected.releaseId === patch.releaseId ? 'tree-selected' : 'tree-unselected'
                              )}
                              on:click={() => selectPatch(version.id, build.id, patch)}
                            >
                              <span class="flex items-center gap-2">
                                <span class="font-medium text-foreground">{patch.name}</span>
                              </span>
                              <span class="text-[11px] text-muted-foreground">{patch.assetCount} assets</span>
                            </button>
                          </div>
                        {/each}
                      </div>
                    </div>
                  {/each}
                </div>
              </div>
            {/each}
          {/if}
        </div>
      </nav>
    </aside>

    <main class="flex-1 p-8 bg-background">
      <div class="max-w-5xl">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-xl font-semibold text-foreground">Public Releases</h1>
            <p class="text-sm text-muted-foreground mt-1">
              Khám phá và tải bản build đã được phát hành.
            </p>
          </div>
        </div>

        <div class="mt-6 grid gap-6">
          <section class="rounded-xl border border-border bg-card overflow-hidden">
            <div class="px-5 py-4 border-b border-border bg-muted/20">
              <h2 class="text-sm font-medium text-foreground">Update Notes</h2>
            </div>
            <div class="px-5 py-4">
              <div class="update-content">{@html updateHtml}</div>
            </div>
          </section>

          <section class="rounded-xl border border-border bg-card overflow-hidden">
            <div class="px-5 py-4 border-b border-border bg-muted/20 flex items-center justify-between">
              <h2 class="text-sm font-medium text-foreground">Artifacts</h2>
              <span class="text-xs text-muted-foreground">
                {selectedArtifacts.length} files
              </span>
            </div>
            <div class="divide-y divide-border/50">
              {#if selectedArtifacts.length > 0}
                {#each selectedArtifacts as asset (asset.id)}
                  <div class="flex items-center justify-between gap-4 px-5 py-4">
                    <div class="min-w-0">
                      <p class="text-sm font-medium text-foreground truncate">{asset.filename}</p>
                      <p class="text-xs text-muted-foreground mt-0.5">{formatBytes(asset.size_bytes)}</p>
                    </div>
                    <a
                      href={`/public/artifacts/${asset.id}/download`}
                      class="inline-flex items-center gap-2 px-3 py-1.5 rounded-md text-xs bg-secondary text-secondary-foreground hover:bg-secondary/80 transition-colors"
                    >
                      <Download class="h-3.5 w-3.5" />
                      Download
                    </a>
                  </div>
                {/each}
              {:else}
                <div class="px-5 py-6 text-sm text-muted-foreground">
                  Chọn một bản vá để xem artifacts.
                </div>
              {/if}
            </div>
          </section>
        </div>
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
    padding-left: 1.35rem;
    position: relative;
  }
  .tree-children::before {
    content: '';
    position: absolute;
    left: 0.3rem;
    top: 0.12rem;
    bottom: -0.1rem;
    width: 1px;
    background: hsl(var(--border));
  }
  .tree-row {
    position: relative;
  }
  .tree-row-last::after {
    content: '';
    position: absolute;
    left: -0.85rem;
    top: calc(50% + 0.35rem);
    bottom: -0.55rem;
    width: 1px;
    background: hsl(var(--background));
  }
  .tree-branch,
  .tree-leaf {
    padding-left: 0.35rem;
  }
  .tree-branch::before,
  .tree-leaf::before {
    content: '';
    position: absolute;
    left: -0.85rem;
    top: 50%;
    width: 0.85rem;
    height: 1px;
    background: hsl(var(--border));
  }
  .tree-leaf {
    justify-content: space-between;
    border-radius: 0.45rem;
    padding: 0.45rem 0.6rem 0.45rem 0.4rem;
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
