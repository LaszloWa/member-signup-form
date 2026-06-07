import type { WizardDraft } from './types'

const DRAFT_STORAGE_KEY = 'spond-membership-draft-v1'
const FIFTEEN_MINUTES_MS = 15 * 60 * 1000

type StoredDraft = {
  expiresAt: number
  draft: WizardDraft
}

export function saveDraft(draft: WizardDraft): void {
  const payload: StoredDraft = {
    draft,
    expiresAt: Date.now() + FIFTEEN_MINUTES_MS,
  }

  sessionStorage.setItem(DRAFT_STORAGE_KEY, JSON.stringify(payload))
}

export function loadDraft(): WizardDraft | null {
  const rawValue = sessionStorage.getItem(DRAFT_STORAGE_KEY)
  if (!rawValue) {
    return null
  }

  try {
    const parsed = JSON.parse(rawValue) as Partial<StoredDraft>

    if (!parsed.expiresAt || !parsed.draft) {
      clearDraft()
      return null
    }

    if (Date.now() > parsed.expiresAt) {
      clearDraft()
      return null
    }

    return parsed.draft
  } catch {
    clearDraft()
    return null
  }
}

export function clearDraft(): void {
  sessionStorage.removeItem(DRAFT_STORAGE_KEY)
}
