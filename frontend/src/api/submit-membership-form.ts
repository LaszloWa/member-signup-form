import type { ErrorEnvelope, SubmissionPayload, SubmissionResponse } from './contracts'
import { buildBackendEndpoint } from './backend-base-url'

const UNEXPECTED_SUBMISSION_RESPONSE_MESSAGE = 'Unexpected submission response.'
const NETWORK_SUBMISSION_ERROR_MESSAGE = 'Network error. Please try again.'

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null
}

function isStringRecord(value: unknown): value is Record<string, string> {
  if (!isRecord(value)) {
    return false
  }

  return Object.values(value).every((item) => typeof item === 'string')
}

function isSubmissionResponse(value: unknown): value is SubmissionResponse {
  if (!isRecord(value)) {
    return false
  }

  return (
    typeof value.submissionId === 'string'
    && typeof value.createdAt === 'string'
    && typeof value.name === 'string'
    && typeof value.email === 'string'
    && typeof value.phoneNumber === 'string'
    && typeof value.birthDate === 'string'
    && typeof value.memberTypeId === 'string'
    && typeof value.clubId === 'string'
    && typeof value.formId === 'string'
  )
}

function isErrorEnvelope(value: unknown): value is ErrorEnvelope {
  if (!isRecord(value)) {
    return false
  }

  return isStringRecord(value.errors)
}

async function parseErrorEnvelope(response: Response): Promise<ErrorEnvelope | null> {
  try {
    const payload: unknown = await response.json()
    return isErrorEnvelope(payload) ? payload : null
  } catch {
    return null
  }
}

export async function submitMembershipForm(
  payload: SubmissionPayload,
): Promise<{ success: true; data: SubmissionResponse } | { success: false; errors: Record<string, string> }> {
  let response: Response

  try {
    response = await fetch(buildBackendEndpoint('/api/v1/forms/public/submissions'), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(payload),
    })
  } catch {
    return {
      success: false,
      errors: {
        server: NETWORK_SUBMISSION_ERROR_MESSAGE,
      },
    }
  }

  if (response.ok) {
    try {
      const responsePayload: unknown = await response.json()
      if (!isSubmissionResponse(responsePayload)) {
        return {
          success: false,
          errors: {
            server: UNEXPECTED_SUBMISSION_RESPONSE_MESSAGE,
          },
        }
      }

      return {
        success: true,
        data: responsePayload,
      }
    } catch {
      return {
        success: false,
        errors: {
          server: UNEXPECTED_SUBMISSION_RESPONSE_MESSAGE,
        },
      }
    }
  }

  const errorEnvelope = await parseErrorEnvelope(response)
  return {
    success: false,
    errors: errorEnvelope?.errors ?? {
      server: `Submission failed with status ${response.status}`,
    },
  }
}