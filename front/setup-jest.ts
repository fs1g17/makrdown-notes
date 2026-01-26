import '@testing-library/jest-dom'
import nock from 'nock'
import axios from 'axios'

// Force axios to use http adapter so nock can intercept
axios.defaults.adapter = 'http'

// Disable real network connections
nock.disableNetConnect()

// Mock next/navigation
jest.mock('next/navigation', () => ({
  useRouter: jest.fn(() => ({
    push: jest.fn(),
    replace: jest.fn(),
    back: jest.fn(),
  })),
  useParams: jest.fn(() => ({})),
  usePathname: jest.fn(() => '/'),
  useSearchParams: jest.fn(() => new URLSearchParams()),
}))

// Clean up nock between tests
beforeEach(() => {
  nock.cleanAll()
})

afterEach(() => {
  nock.cleanAll()
  nock.restore()
})