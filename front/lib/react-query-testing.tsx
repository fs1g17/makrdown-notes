import { QueryClient, QueryClientProvider as QCProvider } from "@tanstack/react-query";

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: false, // Don't retry failed requests in tests
    },
  },
});

export const QueryClientProvider = ({ children }: { children: React.ReactNode }) => (
  <QCProvider client={queryClient}>{children}</QCProvider>
);
