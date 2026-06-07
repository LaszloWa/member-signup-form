import { expect, test } from '@playwright/test'

import { goToReviewStep } from './helpers'

test('shows duplicate submission error from backend', async ({ page }) => {
  const uniqueSuffix = Date.now().toString()
  await page.goto('/')
  await goToReviewStep(page, uniqueSuffix)

  await page.getByTestId('submit-registration').click()
  await expect(page.getByTestId('wizard-success')).toBeVisible()

  await page.getByRole('button', { name: 'Start new registration' }).click()
  await expect(page.getByTestId('step-1')).toBeVisible()

  await goToReviewStep(page, uniqueSuffix)

  await page.getByTestId('submit-registration').click()

  await expect(page.getByRole('alert')).toBeVisible()
  await expect(
    page.getByText(
      'submission: a submission with the same form, club, email, phone number, and birth date already exists',
    ),
  ).toBeVisible()
})
