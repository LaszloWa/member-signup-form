import { expect, test } from '@playwright/test'

test('shows registration-locked state when registrationOpens is in future', async ({ page }) => {
  await page.route('**/api/v1/forms/public', (route) => {
    route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        clubId: 'britsport',
        memberTypes: [
          { id: '8FE4113D4E4020E0DCF887803A886981', name: 'Active Member' },
          { id: '4237C55C5CC3B4B082CBF2540612778E', name: 'Social Member' },
        ],
        formId: 'B171388180BC457D9887AD92B6CCFC86',
        title: 'Coding camp summer 2025',
        registrationOpens: '2099-01-01T00:00:00Z',
      }),
    })
  })

  await page.goto('/')

  await expect(page.getByTestId('registration-locked')).toBeVisible()
  await expect(page.getByTestId('step-1')).toHaveCount(0)
})
