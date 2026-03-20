/// <reference types="vite/client" />

//  turn strictImportMetaEnv on for strict env key
// interface ViteTypeOptions {
// strictImportMetaEnv: unknown;
//}

type EnvironmentTypes = 'development' | 'int' | 'test' | 'prod' | 'aws';
interface ImportMetaEnv {
  readonly VITE_BUILD_DATE: string;
  readonly VITE_VERSION: string;
  readonly VITE_COMMIT: string;
  readonly VITE_BRANCH: string;
  readonly VITE_SERVER_URL: string;
  readonly VITE_OAUTH_LOGIN: string;
  readonly VITE_OAUTH_LOGOUT: string;
  readonly PUBLIC_API_ENDPOINT: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
