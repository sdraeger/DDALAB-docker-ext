// API Types
export interface Service {
  name: string;
  status: string;
}

export interface Status {
  running: boolean;
  services: Service[];
  version: string;
  path: string;
}

export interface PathValidationResult {
  valid: boolean;
  path: string;
  message: string;
  has_compose: boolean;
  has_ddalab_script: boolean;
}

export interface PathSelectionRequest {
  path: string;
}

export interface ExtensionConfig {
  selected_path?: string;
  known_paths?: string[];
  discovered_paths?: string[];
}

export interface Alert {
  message: string;
  severity: 'error' | 'warning' | 'info' | 'success';
}

export interface EnvConfig {
  url: string;
  host: string;
  port: string;
  scheme: string;
  domain: string;
}

export interface ApiResponse<T = any> {
  status: string;
  data?: T;
  error?: string;
}

export interface HealthCheck {
  service: string;
  healthy: boolean;
  status: string;
  message?: string;
  details?: Record<string, string>;
}

export interface SystemHealth {
  overall: boolean;
  services: HealthCheck[];
  timestamp: string;
  config_path: string;
}

// Component Props Types
export interface PathSelectorProps {
  currentPath: string;
  onPathChange: (path: string) => void;
}

export interface ServiceListProps {
  services: Service[];
  onServiceAction: (serviceName: string, action: string) => Promise<void>;
  disabled?: boolean;
}

export interface StatusCardProps {
  title: string;
  value: string;
  subtitle?: string;
  type?: 'status' | 'health' | 'default';
  icon?: React.ReactNode;
}

// Environment Configuration Types
export interface EnvVar {
  key: string;
  value: string;
  comment: string;
  section: string;
  required: boolean;
  secret: boolean;
  line_num: number;
}

export interface EnvFile {
  variables: EnvVar[];
  path: string;
  modified: boolean;
}

export interface UpdateEnvRequest {
  variables: EnvVar[];
}

export interface ValidationError {
  key: string;
  message: string;
}

export interface ValidationResult {
  valid: boolean;
  errors: ValidationError[];
}