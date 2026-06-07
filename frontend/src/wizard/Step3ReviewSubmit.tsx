import {
	useActionState,
	useEffect,
	useRef,
} from 'react'
import { useFormStatus } from 'react-dom'

import { submitMembershipForm } from '../api/submit-membership-form'
import type {
	SubmissionPayload,
	SubmissionResponse,
} from '../api/contracts'
import { clearDraft } from './draftStorage'
import { useWizardContext } from './wizard-context-hook'

function SubmitButton() {
	const { pending } = useFormStatus()

	return (
		<button type="submit" disabled={pending} data-testid="submit-registration">
			{pending ? 'Submitting...' : 'Submit registration'}
		</button>
	)
}

function formDataValue(formData: FormData, key: string): string {
	const value = formData.get(key)
	return typeof value === 'string' ? value : ''
}

function buildSubmissionPayload(formData: FormData): SubmissionPayload {
	return {
		name: formDataValue(formData, 'name').trim(),
		email: formDataValue(formData, 'email').trim().toLowerCase(),
		phoneNumber: formDataValue(formData, 'phoneNumber').trim(),
		birthDate: formDataValue(formData, 'birthDate').trim(),
		memberTypeId: formDataValue(formData, 'memberTypeId').trim(),
		clubId: formDataValue(formData, 'clubId').trim(),
		formId: formDataValue(formData, 'formId').trim(),
	}
}

function formatBirthDateForDisplay(value: string): string {
	return value.split('-').reverse().join('-')
}

type SubmitState =
	| { status: 'idle' }
	| { status: 'error'; errors: Record<string, string> }
	| { status: 'success'; submission: SubmissionResponse }

export function Step3ReviewSubmit() {
	const {
		state: { draft, formDetails },
		setStep,
		resetWizard,
	} = useWizardContext()

	const headingRef = useRef<HTMLHeadingElement>(null)

	useEffect(() => {
		headingRef.current?.focus()
	}, [])

	const [submitState, formAction] = useActionState(
		async (_previousSubmitState: SubmitState, formData: FormData): Promise<SubmitState> => {
			const payload = buildSubmissionPayload(formData)

			const result = await submitMembershipForm(payload)
			if (!result.success) {
				return {
					status: 'error',
					errors: result.errors,
				}
			}

			clearDraft()

			return {
				status: 'success',
				submission: result.data,
			}
		},
		{ status: 'idle' },
	)

	if (!formDetails) {
		return null
	}

	if (submitState.status === 'success') {
		return (
			<section className="card success-banner" aria-live="polite" data-testid="wizard-success">
				<h2>Success</h2>
				<p>Your membership registration has been submitted.</p>
				<p data-testid="submission-id">Submission ID: {submitState.submission.submissionId}</p>
				<button type="button" onClick={resetWizard}>
					Start new registration
				</button>
			</section>
		)
	}

	const selectedMemberType = formDetails.memberTypes.find((memberType) => memberType.id === draft.memberTypeId)
	const selectedMemberTypeName = selectedMemberType?.name ?? 'Unknown member type'
	const submitErrors = submitState.status === 'error' ? submitState.errors : {}
	const hasSubmitErrors = submitState.status === 'error' && Object.keys(submitErrors).length > 0

	return (
		<form className="card" action={formAction} data-testid="step-3">
			<h2 ref={headingRef} tabIndex={-1} data-testid="step-3-heading">
				Step 3: Review and Submit
			</h2>
			<p className="subtle">Review your details before submitting.</p>

			<input type="hidden" name="memberTypeId" value={draft.memberTypeId} />
			<input type="hidden" name="name" value={draft.name} />
			<input type="hidden" name="email" value={draft.email} />
			<input type="hidden" name="phoneNumber" value={draft.phoneNumber} />
			<input type="hidden" name="birthDate" value={draft.birthDate} />
			<input type="hidden" name="clubId" value={formDetails.clubId} />
			<input type="hidden" name="formId" value={formDetails.formId} />

			<dl className="review-list">
				<div>
					<dt>Member type</dt>
					<dd>{selectedMemberTypeName}</dd>
				</div>
				<div>
					<dt>Name</dt>
					<dd>{draft.name}</dd>
				</div>
				<div>
					<dt>Email</dt>
					<dd>{draft.email}</dd>
				</div>
				<div>
					<dt>Phone number</dt>
					<dd>{draft.phoneNumber}</dd>
				</div>
				<div>
					<dt>Birth date</dt>
					<dd data-testid="review-birth-date">{formatBirthDateForDisplay(draft.birthDate)}</dd>
				</div>
			</dl>

			{hasSubmitErrors ? (
				<div className="error-banner" role="alert" aria-live="assertive">
					<p>Submission failed. Please fix the following:</p>
					<ul>
						{Object.entries(submitErrors).map(([field, message]) => (
							<li key={field}>
								{field}: {message}
							</li>
						))}
					</ul>
				</div>
			) : null}

			<div className="actions split">
				<button type="button" className="secondary" onClick={() => setStep(2)}>
					Back
				</button>
				<SubmitButton />
			</div>
		</form>
	)
}
