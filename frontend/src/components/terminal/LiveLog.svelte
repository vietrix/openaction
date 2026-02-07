<script lang="ts">
  import { afterUpdate, onDestroy } from 'svelte';
  import { Copy, Check, ArrowDown } from '@lucide/svelte';
  import Button from '@/components/ui/Button.svelte';
  import { cn } from '@/lib/utils';

  interface LogLine {
    id: number;
    timestamp: string;
    level: 'info' | 'warn' | 'error' | 'success' | 'debug';
    message: string;
  }

  export let logs: LogLine[] = [];
  export let title = 'Build Logs';
  export let className = '';
  export let streamUrl = '';
  export let loading = false;

  let containerRef: HTMLDivElement | null = null;
  let copied = false;
  let autoScroll = true;
  let isHovered = false;
  let copyTimeout: ReturnType<typeof setTimeout> | null = null;
  let socket: WebSocket | null = null;
  let currentStreamUrl = '';
  let entries: LogLine[] = [];
  let lastLogsRef: LogLine[] | null = null;

  afterUpdate(() => {
    if (autoScroll && containerRef) {
      containerRef.scrollTop = containerRef.scrollHeight;
    }
  });

  onDestroy(() => {
    if (copyTimeout) clearTimeout(copyTimeout);
    if (socket) socket.close();
  });

  $: if (logs !== lastLogsRef) {
    lastLogsRef = logs;
    entries = [...logs];
  }

  const appendLine = (message: string) => {
    const timestamp = new Date().toLocaleTimeString('en-US', { hour12: false }).slice(0, 8);
    entries = [
      ...entries,
      {
        id: entries.length + 1,
        timestamp,
        level: parseLevel(message),
        message,
      },
    ];
  };

  const parseLevel = (line: string): LogLine['level'] => {
    if (/ERROR|FAILED|FAIL/i.test(line)) return 'error';
    if (/WARN|WARNING/i.test(line)) return 'warn';
    if (/SUCCESS|PASSED|OK/i.test(line)) return 'success';
    if (/DEBUG/i.test(line)) return 'debug';
    return 'info';
  };

  const connectStream = () => {
    if (!streamUrl) return;
    if (streamUrl === currentStreamUrl && socket) return;
    if (socket) {
      socket.close();
      socket = null;
    }
    currentStreamUrl = streamUrl;
    socket = new WebSocket(streamUrl);
    socket.onmessage = (event) => {
      const payload = typeof event.data === 'string' ? event.data : '';
      payload.split('\n').forEach((line) => {
        const trimmed = line.trim();
        if (trimmed) appendLine(trimmed);
      });
    };
    socket.onclose = () => {
      socket = null;
    };
  };

  $: if (streamUrl && streamUrl !== currentStreamUrl) {
    connectStream();
  } else if (!streamUrl && socket) {
    socket.close();
    socket = null;
    currentStreamUrl = '';
  }

  const handleScroll = () => {
    if (!containerRef) return;
    const { scrollTop, scrollHeight, clientHeight } = containerRef;
    const isAtBottom = scrollHeight - scrollTop - clientHeight < 50;
    autoScroll = isAtBottom;
  };

  const copyLogs = async () => {
    const text = entries.map((log) => `[${log.timestamp}] ${log.message}`).join('\n');
    try {
      await navigator.clipboard.writeText(text);
      copied = true;
      if (copyTimeout) clearTimeout(copyTimeout);
      copyTimeout = setTimeout(() => {
        copied = false;
      }, 2000);
    } catch {
      copied = false;
    }
  };

  const scrollToBottom = () => {
    if (containerRef) {
      containerRef.scrollTop = containerRef.scrollHeight;
      autoScroll = true;
    }
  };

  const getLevelClass = (level: LogLine['level']) => {
    switch (level) {
      case 'error':
        return 'terminal-error';
      case 'warn':
        return 'terminal-warning';
      case 'success':
        return 'terminal-success';
      case 'info':
        return 'terminal-info';
      default:
        return 'text-terminal-foreground';
    }
  };

  const highlightKeywords = (message: string) =>
    message
      .replace(/(ERROR|FATAL|FAILED)/gi, '<span class="terminal-error font-semibold">$1</span>')
      .replace(/(WARN|WARNING)/gi, '<span class="terminal-warning font-semibold">$1</span>')
      .replace(/(SUCCESS|PASSED|OK)/gi, '<span class="terminal-success font-semibold">$1</span>')
      .replace(/(INFO|DEBUG)/gi, '<span class="terminal-info">$1</span>');
</script>

<div
  class={cn(
    'flex flex-col rounded-lg border border-terminal-border bg-terminal overflow-hidden',
    className
  )}
  role="region"
  on:mouseenter={() => (isHovered = true)}
  on:mouseleave={() => (isHovered = false)}
>
  <div class="flex items-center justify-between border-b border-terminal-border px-4 py-2.5">
    <div class="flex items-center gap-2">
      <div class="flex gap-1.5">
        <span class="h-3 w-3 rounded-full bg-status-error/80"></span>
        <span class="h-3 w-3 rounded-full bg-status-warning/80"></span>
        <span class="h-3 w-3 rounded-full bg-status-success/80"></span>
      </div>
      <span class="ml-2 text-xs font-medium text-terminal-foreground">{title}</span>
    </div>

    <div
      class={cn(
        'flex items-center gap-1 transition-opacity duration-150',
        isHovered ? 'opacity-100' : 'opacity-0'
      )}
    >
      <Button
        variant="ghost"
        size="icon"
        className="h-7 w-7 text-terminal-foreground/60 hover:bg-terminal-border hover:text-terminal-foreground"
        on:click={copyLogs}
      >
        {#if copied}
          <Check class="h-3.5 w-3.5 text-status-success" />
        {:else}
          <Copy class="h-3.5 w-3.5" />
        {/if}
      </Button>
    </div>
  </div>

  <div
    bind:this={containerRef}
    on:scroll={handleScroll}
    class="flex-1 overflow-y-auto p-4 min-h-[200px] max-h-[400px]"
  >
    <div class="terminal-text space-y-0.5">
      {#if loading && entries.length === 0}
        <div class="text-xs text-muted-foreground">Đang tải log...</div>
      {:else}
        {#each entries as log (log.id)}
          <div class="flex gap-3 leading-relaxed animate-fade-in">
            <span class="shrink-0 text-muted-foreground/50 select-none">{log.timestamp}</span>
            <span class={getLevelClass(log.level)}>{@html highlightKeywords(log.message)}</span>
          </div>
        {/each}
      {/if}
    </div>
  </div>

  {#if !autoScroll}
    <button
      on:click={scrollToBottom}
      class="absolute bottom-4 right-4 flex items-center gap-1.5 rounded-full bg-terminal-border/90 px-3 py-1.5 text-xs text-terminal-foreground backdrop-blur-sm transition-all hover:bg-terminal-border"
    >
      <ArrowDown class="h-3 w-3" />
      New logs
    </button>
  {/if}
</div>
