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

Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: jest.fn().mockImplementation(query => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: jest.fn(), // Deprecated
    removeListener: jest.fn(), // Deprecated
    addEventListener: jest.fn(),
    removeEventListener: jest.fn(),
    dispatchEvent: jest.fn(),
  })),
});

// Clean up nock between tests
beforeEach(() => {
  nock.cleanAll()
})

afterEach(() => {
  nock.cleanAll()
});