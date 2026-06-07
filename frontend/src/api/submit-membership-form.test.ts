import { afterEach, describe, expect, it, vi } from 'vitest'

import type { SubmissionPayload } from './contracts'
import { submitMembershipForm } from './submit-membership-form'

vi.mock('./backend-base-url', () => ({
  buildBackendEndpoint: (path: string) => path,
}))

function makeResponse(options: {
  ok: boolean
  status: number
  payload?: unknown
  throwsOnJson?: boolean
}): Response {
  return {
    ok: options.ok,
    status: options.status,
    json: async () => {
      if (options.throwsOnJson) {
        throw new Error('invalid json')
      }
      return options.payload
    },
  } as Response
}

const payload: SubmissionPayload = {
  name: 'Test User',
  email: 'test@example.com',
  phoneNumber: '+47 12345678',
  birthDate: '1990-04-21',
  memberTypeId: 'member-type-1',
  clubId: 'club-1',
  formId: 'form-1',
}

afterEach(() => {
  vi.restoreAllMocks()
})

describe('submitMembershipForm', () => {
  it('returns fallback error for malformed success payload', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(makeResponse({ ok: true, status: 200, payload: {} }))

    const result = await submitMembershipForm(payload)

    expect(result.success).toBe(false)
    if (result.success) {
      throw new Error('expected failure result')
    }
    expect(result.errors.server).toBe('Unexpected submission response.')
  })

  it('returns explicit network error when fetch fails', async () => {
    vi.spyOn(globalThis, 'fetch').mockRejectedValue(new Error('network down'))

    const result = await submitMembershipForm(payload)

    expect(result.success).toBe(false)
    if (result.success) {
      throw new Error('expected failure result')
    }
    expect(result.errors.server).toBe('Network error. Please try again.')
  })

  it('returns backend validation errors when envelope is valid', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      makeResponse({
        ok: false,
        status: 400,
        payload: {
          errors: {
            email: 'email is invalid',
          },
        },
      }),
    )

    const result = await submitMembershipForm(payload)

    expect(result.success).toBe(false)
    if (result.success) {
      throw new Error('expected failure result')
    }
    expect(result.errors).toEqual({ email: 'email is invalid' })
  })

  it('falls back to status message when error payload shape is invalid', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      makeResponse({
        ok: false,
        status: 500,
        payload: {
          message: 'oops',
        },
      }),
    )

    const result = await submitMembershipForm(payload)

    expect(result.success).toBe(false)
    if (result.success) {
      throw new Error('expected failure result')
    }
    expect(result.errors.server).toBe('Submission failed with status 500')
  })
})
