import type { Page } from '@playwright/test'

function randomSuffix(): string {
  return `${Date.now()}-${Math.random().toString(36).slice(2, 8)}`
}

export async function goToReviewStep(page: Page, suffix = randomSuffix()) {
  await page.getByLabel('Active Member').click()
  await page.getByRole('button', { name: 'Continue' }).click()

  await page.getByLabel('Name').fill(`Test User ${suffix}`)
  await page.getByLabel('Email').fill(`test.${suffix}@example.com`)
  await page.getByLabel('Phone number').fill('+47 12345678')
  await page.getByLabel('Birth date').fill('1990-04-21')
  await page.getByRole('button', { name: 'Continue' }).click()
}
