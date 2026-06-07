import { Wizard } from './wizard/Wizard'
import { WizardProvider } from './wizard/WizardProvider'

function App() {
  return (
    <main className="app-shell">
      <header className="header">
        <h1>Club Membership Registration</h1>
        <p>Complete the 3-step signup wizard.</p>
      </header>

      <WizardProvider>
        <Wizard />
      </WizardProvider>
    </main>
  )
}

export default App
