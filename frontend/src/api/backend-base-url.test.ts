import { describe, expect, it } from 'vitest'

import { joinBackendEndpoint } from './backend-base-url'

describe('joinBackendEndpoint', () => {
  it('strips trailing slashes from base URL', () => {
    expect(joinBackendEndpoint('https://api.example.com///', '/health')).toBe('https://api.example.com/health')
  })

  it('ensures API path starts with slash', () => {
    expect(joinBackendEndpoint('https://api.example.com', 'api/v1/forms')).toBe('https://api.example.com/api/v1/forms')
  })

  it('returns normalized path when base URL is empty', () => {
    expect(joinBackendEndpoint('', 'api/v1/forms')).toBe('/api/v1/forms')
  })
})
