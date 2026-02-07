const API_BASE = '/actions';
const PUBLIC_BASE = '/public';

type RequestOptions = {
  method?: string;
  body?: unknown;
  headers?: Record<string, string>;
};

const getCookie = (name: string) => {
  if (typeof document === 'undefined') return '';
  const match = document.cookie.match(new RegExp(`(^| )${name}=([^;]+)`));
  return match ? decodeURIComponent(match[2]) : '';
};

const buildHeaders = (method: string, headers?: Record<string, string>) => {
  const result: Record<string, string> = { ...(headers ?? {}) };
  if (method !== 'GET' && method !== 'HEAD') {
    const csrf = getCookie('oa_csrf');
    if (csrf) {
      result['X-CSRF-Token'] = csrf;
    }
  }
  return result;
};

const request = async <T>(path: string, options: RequestOptions = {}) => {
  const method = options.method ?? 'GET';
  const headers = buildHeaders(method, options.headers);
  const init: RequestInit = {
    method,
    headers,
    credentials: 'include',
  };

  if (options.body !== undefined) {
    headers['Content-Type'] = 'application/json';
    init.body = JSON.stringify(options.body);
  }

  const response = await fetch(`${API_BASE}${path}`, init);
  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || response.statusText);
  }
  if (response.status === 204) return null as T;
  return (await response.json()) as T;
};

const requestText = async (path: string) => {
  const response = await fetch(`${API_BASE}${path}`, { credentials: 'include' });
  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || response.statusText);
  }
  return response.text();
};

const publicRequest = async <T>(path: string) => {
  const response = await fetch(`${PUBLIC_BASE}${path}`);
  if (!response.ok) {
    const message = await response.text();
    throw new Error(message || response.statusText);
  }
  if (response.status === 204) return null as T;
  return (await response.json()) as T;
};

export const api = {
  getMe: () => request<{ id: string; email: string; name: string }>('/auth/me'),
  login: (email: string, password: string) => request('/auth/login', { method: 'POST', body: { email, password } }),
  logout: () => request('/auth/logout', { method: 'POST' }),
  getProjects: () => request<Array<{ id: string; name: string; repo_url: string; default_branch: string }>>('/projects'),
  getProjectPipelines: (projectId: string) =>
    request<Array<{
      id: string;
      status: string;
      commit_hash: string;
      branch: string;
      triggered_by: string;
      started_at: number;
      finished_at: number;
    }>>(`/projects/${projectId}/pipelines`),
  getPipeline: (pipelineId: string) =>
    request<{
      id: string;
      project_id: string;
      status: string;
      commit_hash: string;
      branch: string;
      triggered_by: string;
      started_at: number;
      finished_at: number;
    }>(`/pipelines/${pipelineId}`),
  getPipelineSteps: (pipelineId: string) =>
    request<Array<{
      id: string;
      name: string;
      status: string;
      started_at: number;
      finished_at: number;
      log_path: string;
    }>>(`/pipelines/${pipelineId}/steps`),
  getLogSnapshot: (pipelineId: string, stepId: string) =>
    requestText(`/pipelines/${pipelineId}/logs?step_id=${stepId}`),
  getPublicReleases: () =>
    publicRequest<
      Array<{
        id: string;
        project_id: string;
        version: string;
        build: string;
        patch: string;
        created_at: number;
        update_path: string;
      }>
    >('/releases'),
  getPublicRelease: (id: string) =>
    publicRequest<{
      id: string;
      project_id: string;
      version: string;
      build: string;
      patch: string;
      created_at: number;
      update_path: string;
      update_md: string;
    }>(`/releases/${id}`),
  getPublicArtifacts: (releaseId: string) =>
    publicRequest<
      Array<{
        id: string;
        filename: string;
        size_bytes: number;
        blob_path: string;
        created_at: number;
      }>
    >(`/releases/${releaseId}/artifacts`),
};

export const buildLogStreamUrl = (pipelineId: string, logPath: string) => {
  const url = new URL(`/actions/pipelines/${pipelineId}/logs/stream`, window.location.origin);
  url.searchParams.set('path', logPath);
  url.protocol = url.protocol === 'https:' ? 'wss:' : 'ws:';
  return url.toString();
};
