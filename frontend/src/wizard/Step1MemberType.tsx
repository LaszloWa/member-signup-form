import {
  useEffect,
  useRef,
  useState,
  type SubmitEventHandler,
} from 'react'

import { useWizardContext } from './wizard-context-hook'

export function Step1MemberType() {
  const {
    state: { draft, formDetails },
    updateDraft,
    setStep,
  } = useWizardContext()

  const headingRef = useRef<HTMLHeadingElement>(null)
  const [hasAttemptedContinue, setHasAttemptedContinue] = useState(false)

  useEffect(() => {
    headingRef.current?.focus()
  }, [])

  if (!formDetails) {
    return null
  }

  const handleContinue: SubmitEventHandler<HTMLFormElement> = (event) => {
    event.preventDefault()
    if (!draft.memberTypeId) {
      setHasAttemptedContinue(true)
      return
    }

    setStep(2)
  }

  const showMemberTypeError = hasAttemptedContinue && !draft.memberTypeId
  const memberTypeErrorId = showMemberTypeError ? 'step1-member-type-error' : undefined

  return (
    <form className="card" onSubmit={handleContinue} data-testid="step-1">
      <h2 ref={headingRef} tabIndex={-1} data-testid="step-1-heading">
        Step 1: {formDetails.title}
      </h2>
      <p className="subtle">{formDetails.description ?? 'Choose which member type you want to sign up as.'}</p>

      <fieldset
        className="member-type-grid"
        aria-invalid={showMemberTypeError}
        aria-describedby={memberTypeErrorId}
        data-testid="member-type-options"
      >
        <legend>Member types</legend>
        {formDetails.memberTypes.map((memberType) => {
          const checked = draft.memberTypeId === memberType.id

          return (
            <label key={memberType.id} className={`member-option ${checked ? 'selected' : ''}`}>
              <input
                type="radio"
                name="memberTypeId"
                value={memberType.id}
                checked={checked}
                onChange={() => updateDraft({ memberTypeId: memberType.id })}
              />
              <span>{memberType.name}</span>
            </label>
          )
        })}
      </fieldset>

      {showMemberTypeError ? (
        <p id="step1-member-type-error" className="field-error" role="alert">
          Select a member type to continue.
        </p>
      ) : null}

      <div className="actions">
        <button type="submit">
          Continue
        </button>
      </div>
    </form>
  )
}
