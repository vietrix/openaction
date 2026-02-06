<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import PipelineNode from './PipelineNode.svelte';
  import { cn } from '@/lib/utils';

  interface Node {
    id: string;
    name: string;
    status: 'success' | 'error' | 'running' | 'pending';
    duration?: string;
    step?: string;
    dependencies: string[];
  }

  export let nodes: Node[] = [];
  export let className = '';
  export let selectedNode: string | null = null;

  const dispatch = createEventDispatcher<{ select: string }>();

  const computePositions = (items: Node[]) => {
    const layers: string[][] = [];

    const getLayer = (nodeId: string, depth = 0): number => {
      const node = items.find((n) => n.id === nodeId);
      if (!node) return depth;
      if (node.dependencies.length === 0) return 0;
      return Math.max(...node.dependencies.map((d) => getLayer(d, depth) + 1));
    };

    items.forEach((node) => {
      const layer = getLayer(node.id);
      if (!layers[layer]) layers[layer] = [];
      layers[layer].push(node.id);
    });

    const newPositions: Record<string, { x: number; y: number }> = {};
    const nodeWidth = 180;
    const nodeHeight = 80;
    const horizontalGap = 80;
    const verticalGap = 24;

    layers.forEach((layer, layerIndex) => {
      const totalHeight = layer.length * nodeHeight + (layer.length - 1) * verticalGap;
      const startY = (300 - totalHeight) / 2;

      layer.forEach((nodeId, nodeIndex) => {
        newPositions[nodeId] = {
          x: layerIndex * (nodeWidth + horizontalGap) + 20,
          y: startY + nodeIndex * (nodeHeight + verticalGap),
        };
      });
    });

    return newPositions;
  };

  $: positions = computePositions(nodes);
  $: graphWidth = Math.max(
    320,
    ...Object.values(positions).map((pos) => pos.x + 200)
  );

  const generatePath = (fromId: string, toId: string) => {
    const from = positions[fromId];
    const to = positions[toId];
    if (!from || !to) return '';

    const nodeWidth = 160;
    const nodeHeight = 60;

    const startX = from.x + nodeWidth;
    const startY = from.y + nodeHeight / 2;
    const endX = to.x;
    const endY = to.y + nodeHeight / 2;

    const controlOffset = (endX - startX) / 2;

    return `M ${startX} ${startY} C ${startX + controlOffset} ${startY}, ${endX - controlOffset} ${endY}, ${endX} ${endY}`;
  };

  $: edges = nodes.flatMap((node) =>
    node.dependencies.map((dep) => {
      const depNode = nodes.find((n) => n.id === dep);
      return {
        from: dep,
        to: node.id,
        active: node.status === 'running' || depNode?.status === 'running',
      };
    })
  );
</script>

<div class={cn('relative h-[320px] w-full overflow-hidden', className)}>
  <div class="h-full w-full overflow-x-auto overflow-y-hidden">
    <div class="relative h-full" style={`width: ${graphWidth}px;`}>
      <svg class="absolute inset-0 h-full w-full pointer-events-none">
        <defs>
          <marker id="arrowhead" markerWidth="10" markerHeight="7" refX="9" refY="3.5" orient="auto">
            <polygon points="0 0, 10 3.5, 0 7" class="fill-border" />
          </marker>
        </defs>
        {#each edges as edge, i (edge.from + edge.to + i)}
          <path
            d={generatePath(edge.from, edge.to)}
            class={cn('pipeline-edge', edge.active && 'pipeline-edge-active')}
            marker-end="url(#arrowhead)"
          />
        {/each}
      </svg>

      {#each nodes as node (node.id)}
        <button
          type="button"
          class="absolute transition-all duration-300 ease-out text-left p-0 m-0 border-0 bg-transparent"
          style={`left: ${positions[node.id]?.x || 0}px; top: ${positions[node.id]?.y || 0}px;`}
          on:click={() => dispatch('select', node.id)}
        >
          <PipelineNode
            id={node.id}
            name={node.name}
            status={node.status}
            duration={node.duration}
            step={node.step}
            className={
              node.id === selectedNode
                ? 'ring-2 ring-primary ring-offset-2 ring-offset-background'
              : ''
            }
          />
        </button>
      {/each}
    </div>
  </div>
</div>
