import type { FormDetails } from './contracts'
import { buildBackendEndpoint } from './backend-base-url'

function isMemberType(value: unknown): value is { id: string; name: string } {
  if (!value || typeof value !== 'object') {
    return false
  }

  const maybeMemberType = value as { id?: unknown; name?: unknown }
  return typeof maybeMemberType.id === 'string' && typeof maybeMemberType.name === 'string'
}

function isFormDetails(value: unknown): value is FormDetails {
  if (!value || typeof value !== 'object') {
    return false
  }

  const maybeFormDetails = value as {
    clubId?: unknown
    memberTypes?: unknown
    formId?: unknown
    title?: unknown
    description?: unknown
    registrationOpens?: unknown
  }

  if (
    typeof maybeFormDetails.clubId !== 'string' ||
    !Array.isArray(maybeFormDetails.memberTypes) ||
    typeof maybeFormDetails.formId !== 'string' ||
    typeof maybeFormDetails.title !== 'string' ||
    typeof maybeFormDetails.registrationOpens !== 'string'
  ) {
    return false
  }

  if (
    maybeFormDetails.description !== undefined &&
    typeof maybeFormDetails.description !== 'string'
  ) {
    return false
  }

  return maybeFormDetails.memberTypes.every(isMemberType)
}

export async function fetchFormDetails(signal?: AbortSignal): Promise<FormDetails> {
  const response = await fetch(buildBackendEndpoint('/api/v1/forms/public'), {
    method: 'GET',
    signal,
  })

  if (!response.ok) {
    throw new Error(`Failed to load form details: ${response.status}`)
  }

  const data: unknown = await response.json()
  if (!isFormDetails(data)) {
    throw new Error('Failed to load form details: invalid response shape')
  }

  return data
}