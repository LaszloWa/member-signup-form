import {
  useEffect,
  useRef,
  useState,
  type FormEvent,
} from 'react'

import { validateStep2 } from './validation'
import { useWizardContext } from './wizard-context-hook'
import type { Step2Errors } from './types'

export function Step2PersonalInfo() {
  const {
    state: { draft },
    updateDraft,
    setStep,
  } = useWizardContext()

  const [errors, setErrors] = useState<Step2Errors>({})
  const headingRef = useRef<HTMLHeadingElement>(null)

  useEffect(() => {
    headingRef.current?.focus()
  }, [])

  function handleContinue(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()

    const nextErrors = validateStep2(draft)
    setErrors(nextErrors)

    if (Object.keys(nextErrors).length > 0) {
      return
    }

    setStep(3)
  }

  const nameErrorId = errors.name ? 'step2-name-error' : undefined
  const emailErrorId = errors.email ? 'step2-email-error' : undefined
  const phoneErrorId = errors.phoneNumber ? 'step2-phone-error' : undefined
  const birthDateErrorId = errors.birthDate ? 'step2-birth-date-error' : undefined
  const maxBirthDate = new Date().toISOString().slice(0, 10)

  return (
    <form className="card" onSubmit={handleContinue} noValidate>
      <h2 ref={headingRef} tabIndex={-1} data-testid="step-2-heading">
        Step 2: Personal Information
      </h2>
      <p className="subtle">Tell us a bit about yourself.</p>

      <label htmlFor="personal-info-name">Name</label>
      <input
        id="personal-info-name"
        type="text"
        name="name"
        value={draft.name}
        onChange={(event) => updateDraft({ name: event.target.value })}
        autoComplete="name"
        required
        maxLength={120}
        minLength={1}
        aria-invalid={Boolean(errors.name)}
        aria-describedby={nameErrorId}
      />
      {errors.name ? (
        <span id="step2-name-error" className="field-error">
          {errors.name}
        </span>
      ) : null}

      <label htmlFor="personal-info-email">Email</label>
      <input
        id="personal-info-email"
        type="email"
        name="email"
        value={draft.email}
        onChange={(event) => updateDraft({ email: event.target.value })}
        autoComplete="email"
        required
        maxLength={254}
        minLength={3}
        aria-invalid={Boolean(errors.email)}
        aria-describedby={emailErrorId}
      />
      {errors.email ? (
        <span id="step2-email-error" className="field-error">
          {errors.email}
        </span>
      ) : null}

      <label htmlFor="personal-info-phone-number">Phone number</label>
      <input
        id="personal-info-phone-number"
        type="tel"
        name="phoneNumber"
        value={draft.phoneNumber}
        onChange={(event) => updateDraft({ phoneNumber: event.target.value })}
        autoComplete="tel"
        required
        minLength={6}
        maxLength={20}
        aria-invalid={Boolean(errors.phoneNumber)}
        aria-describedby={phoneErrorId}
      />
      {errors.phoneNumber ? (
        <span id="step2-phone-error" className="field-error">
          {errors.phoneNumber}
        </span>
      ) : null}

      <label htmlFor="personal-info-birth-date">Birth date</label>
      <input
        id="personal-info-birth-date"
        type="date"
        name="birthDate"
        value={draft.birthDate}
        onChange={(event) => updateDraft({ birthDate: event.target.value })}
        required
        max={maxBirthDate}
        aria-invalid={Boolean(errors.birthDate)}
        aria-describedby={birthDateErrorId}
      />
      {errors.birthDate ? (
        <span id="step2-birth-date-error" className="field-error">
          {errors.birthDate}
        </span>
      ) : null}

      <div className="actions split">
        <button type="button" className="secondary" onClick={() => setStep(1)}>
          Back
        </button>
        <button type="submit">Continue</button>
      </div>
    </form>
  )
}
