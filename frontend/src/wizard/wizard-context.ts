import { createContext } from 'react'

import type { WizardDraft, WizardState } from './types'

export type WizardContextValue = {
  state: WizardState
  setStep: (step: 1 | 2 | 3) => void
  updateDraft: (patch: Partial<WizardDraft>) => void
  resetWizard: () => void
}

export const WizardContext = createContext<WizardContextValue | null>(null)