import type { FormDetails } from '../api/contracts'

export type WizardDraft = {
  memberTypeId: string
  name: string
  email: string
  phoneNumber: string
  birthDate: string
}

export type Step2Errors = Partial<Record<'name' | 'email' | 'phoneNumber' | 'birthDate', string>>

export type WizardState = {
  step: 1 | 2 | 3
  draft: WizardDraft
  formDetails: FormDetails | null
  loadingForm: boolean
  formLoadError: string | null
  registrationLocked: boolean
}

export const emptyDraft: WizardDraft = {
  memberTypeId: '',
  name: '',
  email: '',
  phoneNumber: '',
  birthDate: '',
}
