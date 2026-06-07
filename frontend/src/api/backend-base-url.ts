const backendBaseUrl = import.meta.env.VITE_BACKEND_BASE_URL?.trim() ?? ''

export function joinBackendEndpoint(baseUrl: string, path: string): string {
  const normalizedPath = path.startsWith('/') ? path : `/${path}`
  const normalizedBaseUrl = baseUrl.trim().replace(/\/+$/, '')

  if (!normalizedBaseUrl) {
    return normalizedPath
  }

  return `${normalizedBaseUrl}${normalizedPath}`
}

export function buildBackendEndpoint(path: string): string {
  const normalizedPath = path.startsWith('/') ? path : `/${path}`

  if (backendBaseUrl) {
    return joinBackendEndpoint(backendBaseUrl, normalizedPath)
  }

  if (import.meta.env.DEV) {
    return normalizedPath
  }

  throw new Error('Missing VITE_BACKEND_BASE_URL in production environment.')
}