import { ref, provide, inject } from 'vue'

interface Auth {
  isLoggedIn: ref<boolean>;
  login: () => void;
  logout: () => void;
}

const AuthSymbol = Symbol('auth')

export function provideAuth(): void {
  const isLoggedIn = ref<boolean>(false)

  const login = (): void => {
    isLoggedIn.value = true
  }

  const logout = (): void => {
    isLoggedIn.value = false
  }

  const auth: Auth = {
    isLoggedIn,
    login,
    logout
  }

  provide(AuthSymbol, auth)
}

export function useAuth(): Auth {
  const auth = inject<Auth | undefined>(AuthSymbol)
  if (!auth) {
    throw new Error('No auth provided')
  }
  return auth
}