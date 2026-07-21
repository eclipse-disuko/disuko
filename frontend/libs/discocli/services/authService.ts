import {getApi} from '@cli/api';

type LoginRequestPayload = {
  projectUUID: string;
  token: string;
};
const {api} = getApi();

class AuthService {
  private static instance: AuthService;
  private authProjectUuid: string | null = null;
  private authIsGroup: boolean | null = null;

  private constructor() {
    api.defaults.withCredentials = true;
  }

  public static getInstance(): AuthService {
    if (!AuthService.instance) {
      AuthService.instance = new AuthService();
    }
    return AuthService.instance;
  }

  /**
   * Perform login request. Backend sets HttpOnly cookies; we persist projectUuid
   * so the router guard can treat the user as authenticated.
   */
  async login(projectUuid: string, token: string): Promise<string> {
    const payload: LoginRequestPayload = {projectUUID: projectUuid, token};
    const response = await api.post('/auth/login', payload, {withCredentials: true});

    if (response.status !== 200) {
      throw new Error('Login failed');
    }

    const authFromInfo = await this.fetchAuthInfo();
    if (authFromInfo) {
      return authFromInfo;
    }
    this.authProjectUuid = projectUuid;
    return projectUuid;
  }

  async fetchAuthInfo(): Promise<string | null> {
    try {
      const response = await api.get('/auth/info', {withCredentials: true});
      const projectUuid = response?.data?.projectUUID || response?.data?.projectUuid;
      const isGroup = response?.data?.isGroup;
      if (projectUuid) {
        this.authProjectUuid = projectUuid;
        this.authIsGroup = typeof isGroup === 'boolean' ? isGroup : null;
        return projectUuid;
      }
    } catch {
      // ignore
    }
    this.authProjectUuid = null;
    this.authIsGroup = null;
    return null;
  }

  async logout(): Promise<void> {
    try {
      await api.get('/auth/logout', {withCredentials: true});
    } finally {
      this.authProjectUuid = null;
      this.authIsGroup = null;
    }
  }

  isAuthenticated(): boolean {
    return Boolean(this.authProjectUuid);
  }

  getAuth(): string | null {
    return this.authProjectUuid;
  }

  getIsGroup(): boolean | null {
    return this.authIsGroup;
  }
}

export const authService = AuthService.getInstance();
export default AuthService;
