import { defineConfig, devices } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',
  timeout: 30_000,
  expect: {
    timeout: 5_000,
  },
  use: {
    baseURL: 'http://localhost:5173',
    trace: 'on-first-retry',
  },
  webServer: [
    {
      command:
        "BACKEND_ALLOWED_ORIGINS='http://localhost:5173' BACKEND_REPOSITORY=sqlite BACKEND_SQLITE_PATH='./data/e2e/backend-e2e.db' go run ./cmd/api",
      url: 'http://localhost:8080/health',
      cwd: '../backend',
      reuseExistingServer: !process.env.CI,
      timeout: 120_000,
    },
    {
      command: 'npm run dev',
      url: 'http://localhost:5173',
      reuseExistingServer: !process.env.CI,
      timeout: 120_000,
    },
  ],
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
})
