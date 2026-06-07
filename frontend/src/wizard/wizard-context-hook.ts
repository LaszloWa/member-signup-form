import { useContext } from 'react'

import { WizardContext, type WizardContextValue } from './wizard-context'

export function useWizardContext(): WizardContextValue {
  const context = useContext(WizardContext)
  if (!context) {
    throw new Error('useWizardContext must be used inside WizardProvider')
  }

  return context
}
