export type MemberType = {
  id: string
  name: string
}

export type FormDetails = {
  clubId: string
  memberTypes: MemberType[]
  formId: string
  title: string
  description?: string
  registrationOpens: string
}

export type SubmissionPayload = {
  name: string
  email: string
  phoneNumber: string
  birthDate: string
  memberTypeId: string
  clubId: string
  formId: string
}

export type SubmissionResponse = SubmissionPayload & {
  submissionId: string
  createdAt: string
}

export type ErrorEnvelope = {
  errors: Record<string, string>
}
