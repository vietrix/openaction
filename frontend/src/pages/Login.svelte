<script lang="ts">
  import { navigate } from '@/lib/router';
  import { api } from '@/lib/api';
  import { Shield, LockKeyhole } from '@lucide/svelte';

  let email = 'admin@openaction.local';
  let password = 'admin123';
  let error = '';
  let isSubmitting = false;

  const handleSubmit = async () => {
    if (isSubmitting) return;
    error = '';
    isSubmitting = true;
    try {
      await api.login(email.trim(), password);
      navigate('/pipelines');
    } catch (err) {
      error = err instanceof Error ? err.message : 'Đăng nhập thất bại';
    } finally {
      isSubmitting = false;
    }
  };
</script>

<div class="min-h-screen bg-background text-foreground flex">
  <div class="hidden lg:flex w-[48%] p-10 border-r border-border bg-muted/20">
    <div class="max-w-md">
      <div class="flex items-center gap-3 text-primary">
        <Shield class="h-6 w-6" />
        <span class="text-lg font-semibold">OpenAction</span>
      </div>
      <h1 class="text-3xl font-semibold mt-8 leading-tight">
        Điều khiển CI/CD nhanh, tối giản, và tự lưu trữ.
      </h1>
      <p class="text-sm text-muted-foreground mt-4 leading-relaxed">
        Đăng nhập để theo dõi pipeline, streaming log, và quản lý release trong một giao diện thống nhất.
      </p>
      <div class="mt-8 rounded-xl border border-border bg-card p-5">
        <div class="text-xs text-muted-foreground uppercase tracking-wide">Demo Access</div>
        <div class="mt-3 space-y-2 text-sm">
          <div class="flex items-center justify-between">
            <span>Email</span>
            <span class="font-mono text-primary">admin@openaction.local</span>
          </div>
          <div class="flex items-center justify-between">
            <span>Password</span>
            <span class="font-mono text-primary">admin123</span>
          </div>
        </div>
      </div>
    </div>
  </div>

  <div class="flex-1 flex items-center justify-center p-6">
    <form
      class="w-full max-w-sm rounded-xl border border-border bg-card p-6 shadow-lg"
      on:submit|preventDefault={handleSubmit}
    >
      <div class="flex items-center gap-2 text-sm font-medium text-foreground">
        <LockKeyhole class="h-4 w-4 text-muted-foreground" />
        Đăng nhập
      </div>
      <p class="text-xs text-muted-foreground mt-1">
        Dùng session cookie cho browser, token cho CLI.
      </p>

      <div class="mt-5 space-y-4">
        <label class="block text-xs text-muted-foreground">
          Email
          <input
            type="email"
            class="mt-2 w-full rounded-md border border-border bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/40"
            bind:value={email}
            required
          />
        </label>
        <label class="block text-xs text-muted-foreground">
          Mật khẩu
          <input
            type="password"
            class="mt-2 w-full rounded-md border border-border bg-background px-3 py-2 text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/40"
            bind:value={password}
            required
          />
        </label>
        {#if error}
          <div class="rounded-md border border-status-error/40 bg-status-error-bg px-3 py-2 text-xs text-status-error">
            {error}
          </div>
        {/if}
        <button
          type="submit"
          class="w-full inline-flex items-center justify-center gap-2 rounded-md bg-primary px-3 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90 disabled:opacity-60"
          disabled={isSubmitting}
        >
          {isSubmitting ? 'Đang đăng nhập...' : 'Đăng nhập'}
        </button>
      </div>
    </form>
  </div>
</div>
