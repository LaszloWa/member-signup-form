import type { WizardDraft } from './types'
import type { Step2Errors } from './types'

const phonePattern = /^[0-9+()\-\s]{6,20}$/

function validateEmail(email: string): boolean {
  return /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
}

function validateBirthDate(value: string): string | null {
  if (!value) {
    return 'Birth date is required.'
  }

  const selected = new Date(`${value}T00:00:00Z`)
  if (Number.isNaN(selected.getTime())) {
    return 'Birth date must be a valid date.'
  }

  const now = new Date()
  if (selected > now) {
    return 'Birth date cannot be in the future.'
  }

  const oldestValid = new Date()
  oldestValid.setUTCFullYear(oldestValid.getUTCFullYear() - 120)
  if (selected < oldestValid) {
    return 'Birth date looks unrealistically old.'
  }

  return null
}

export function validateStep2(draft: WizardDraft): Step2Errors {
  const nextErrors: Step2Errors = {}

  const name = draft.name.trim()
  const email = draft.email.trim()
  const phone = draft.phoneNumber.trim()

  if (!name) {
    nextErrors.name = 'Name is required.'
  } else if (name.length > 120) {
    nextErrors.name = 'Name must be 120 characters or fewer.'
  }

  if (!email) {
    nextErrors.email = 'Email is required.'
  } else if (!validateEmail(email)) {
    nextErrors.email = 'Email is invalid.'
  }

  if (!phone) {
    nextErrors.phoneNumber = 'Phone number is required.'
  } else if (!phonePattern.test(phone)) {
    nextErrors.phoneNumber = 'Phone number format is invalid.'
  }

  const birthDateError = validateBirthDate(draft.birthDate)
  if (birthDateError) {
    nextErrors.birthDate = birthDateError
  }

  return nextErrors
}
