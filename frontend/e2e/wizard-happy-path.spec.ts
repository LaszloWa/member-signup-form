import { expect, test } from '@playwright/test'

import { goToReviewStep } from './helpers'

test('completes happy path and shows success', async ({ page }) => {
  await page.goto('/')
  await expect(page.getByTestId('step-1')).toBeVisible()
  await expect(page.getByText('Select a member type to continue.')).toHaveCount(0)

  await page.getByRole('button', { name: 'Continue' }).click()
  await expect(page.getByText('Select a member type to continue.')).toBeVisible()

  await goToReviewStep(page)
  await expect(page.getByTestId('step-3')).toBeVisible()
  await expect(page.getByTestId('review-birth-date')).toHaveText('21-04-1990')

  await page.getByTestId('submit-registration').click()

  await expect(page.getByTestId('wizard-success')).toBeVisible()
  await expect(page.getByTestId('submission-id')).toBeVisible()

  await expect.poll(async () => page.evaluate(() => sessionStorage.getItem('spond-membership-draft-v1'))).toBeNull()

  await page.getByRole('button', { name: 'Start new registration' }).click()
  await expect(page.getByTestId('step-1')).toBeVisible()
  await expect(page.locator('input[name="memberTypeId"]:checked')).toHaveCount(0)
  await expect(page.getByText('Select a member type to continue.')).toHaveCount(0)

  await page.getByLabel('Active Member').click()
  await page.getByRole('button', { name: 'Continue' }).click()
  await expect(page.getByTestId('step-2-heading')).toBeVisible()
  await expect(page.getByLabel('Name')).toHaveValue('')
  await expect(page.getByLabel('Email')).toHaveValue('')
  await expect(page.getByLabel('Phone number')).toHaveValue('')
  await expect(page.getByLabel('Birth date')).toHaveValue('')
})
