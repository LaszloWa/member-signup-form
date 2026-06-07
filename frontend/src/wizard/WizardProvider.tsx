import {
  useEffect,
  useState,
  type PropsWithChildren,
} from 'react'

import { fetchFormDetails } from '../api/fetch-form-details'
import { clearDraft, loadDraft, saveDraft } from './draftStorage'
import { WizardContext, type WizardContextValue } from './wizard-context'
import { emptyDraft, type WizardDraft, type WizardState } from './types'

function getInitialDraft(): WizardDraft {
  return loadDraft() ?? emptyDraft
}

export function WizardProvider({ children }: PropsWithChildren) {
  const [state, setState] = useState<WizardState>({
    step: 1,
    draft: getInitialDraft(),
    formDetails: null,
    loadingForm: true,
    formLoadError: null,
    registrationLocked: false,
  })

  function patchState(patch: Partial<WizardState>) {
    setState((prev) => ({ ...prev, ...patch }))
  }

  useEffect(() => {
    const controller = new AbortController()

    fetchFormDetails(controller.signal)
      .then((formDetails) => {
        const opensAt = Date.parse(formDetails.registrationOpens)
        const registrationLocked = !Number.isNaN(opensAt) && opensAt > Date.now()

        patchState({
          formDetails,
          loadingForm: false,
          formLoadError: null,
          registrationLocked,
        })
      })
      .catch(() => {
        if (controller.signal.aborted) {
          return
        }

        patchState({
          loadingForm: false,
          formLoadError: 'Could not load form details. Please try again.',
        })
      })

    return () => {
      controller.abort()
    }
  }, [])

  function setStep(step: 1 | 2 | 3) {
    patchState({ step })
  }

  function updateDraft(patch: Partial<WizardDraft>) {
    setState((prev) => {
      const nextDraft = { ...prev.draft, ...patch }
      saveDraft(nextDraft)
      return {
        ...prev,
        draft: nextDraft,
      }
    })
  }

  function resetWizard() {
    clearDraft()
    setState((prev) => ({
      ...prev,
      step: 1,
      draft: emptyDraft,
    }))
  }

  const contextValue: WizardContextValue = {
    state,
    setStep,
    updateDraft,
    resetWizard,
  }

  return <WizardContext.Provider value={contextValue}>{children}</WizardContext.Provider>
}
