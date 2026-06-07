import { Step1MemberType } from './Step1MemberType'
import { Step2PersonalInfo } from './Step2PersonalInfo'
import { Step3ReviewSubmit } from './Step3ReviewSubmit'
import { useWizardContext } from './wizard-context-hook'

export function Wizard() {
  const {
    state: {
      step,
      formDetails,
      loadingForm,
      formLoadError,
      registrationLocked,
    },
  } = useWizardContext()

  if (loadingForm) {
    return (
      <p className="status" data-testid="wizard-loading">
        Loading membership form...
      </p>
    )
  }

  if (formLoadError) {
    return (
      <div className="error-banner" role="alert" data-testid="wizard-load-error">
        {formLoadError}
      </div>
    )
  }

  if (!formDetails) {
    return null
  }

  if (registrationLocked) {
    return (
      <section className="card info-banner" role="status" aria-live="polite" data-testid="registration-locked">
        <h2>Registration Is Not Open Yet</h2>
        <p>
          This form opens on{' '}
          {new Date(formDetails.registrationOpens).toLocaleString(undefined, {
            dateStyle: 'medium',
            timeStyle: 'short',
          })}
          .
        </p>
      </section>
    )
  }

  return (
    <>
      <nav className="stepper" aria-label="Form steps" data-testid="stepper">
        <span className={step >= 1 ? 'active' : ''} aria-current={step === 1 ? 'step' : undefined}>
          1. Member type
        </span>
        <span className={step >= 2 ? 'active' : ''} aria-current={step === 2 ? 'step' : undefined}>
          2. Personal info
        </span>
        <span className={step >= 3 ? 'active' : ''} aria-current={step === 3 ? 'step' : undefined}>
          3. Review
        </span>
      </nav>

      {step === 1 ? <Step1MemberType /> : null}
      {step === 2 ? <Step2PersonalInfo /> : null}
      {step === 3 ? <Step3ReviewSubmit /> : null}
    </>
  )
}
