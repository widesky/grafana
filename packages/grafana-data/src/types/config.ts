import { SystemDateFormatSettings } from '../datetime';
import { MapLayerOptions } from '../geo/layer';
import { GrafanaTheme2, WideSkyCustomTheme } from '../themes';

import { DataSourceInstanceSettings } from './datasource';
import { FeatureToggles } from './featureToggles.gen';
import { PanelPluginMeta } from './panel';

import { GrafanaTheme, IconName, NavLinkDTO, OrgRole } from '.';

/**
 * Describes the build information that will be available via the Grafana configuration.
 *
 * @public
 */
export interface BuildInfo {
  // This MUST be a semver-ish version string, such as "11.0.0-54321"
  version: string;
  // Version to show in the UI instead of version
  versionString: string;
  commit: string;
  env: string;
  edition: GrafanaEdition;
  latestVersion: string;
  hasUpdate: boolean;
  hideVersion: boolean;
}

/**
 * @internal
 */
export enum GrafanaEdition {
  OpenSource = 'Open Source',
  Pro = 'Pro',
  Enterprise = 'Enterprise',
}

/**
 * Describes the license information about the current running instance of Grafana.
 *
 * @public
 */
export interface LicenseInfo {
  expiry: number;
  licenseUrl: string;
  stateInfo: string;
  edition: GrafanaEdition;
  enabledFeatures: { [key: string]: boolean };
  trialExpiry?: number;
}

/**
 * Describes GrafanaJavascriptAgentConfig integration config
 *
 * @public
 */
export interface GrafanaJavascriptAgentConfig {
  enabled: boolean;
  customEndpoint: string;
  errorInstrumentalizationEnabled: boolean;
  consoleInstrumentalizationEnabled: boolean;
  webVitalsInstrumentalizationEnabled: boolean;
  apiKey: string;
}

export interface UnifiedAlertingConfig {
  minInterval: string;
  // will be undefined if alerStateHistory is not enabled
  alertStateHistoryBackend?: string;
  // will be undefined if implementation is not "multiple"
  alertStateHistoryPrimary?: string;
}

/** Supported OAuth services
 *
 * @public
 */
export type OAuth =
  | 'github'
  | 'gitlab'
  | 'google'
  | 'generic_oauth'
  // | 'grafananet' Deprecated. Key always changed to "grafana_com"
  | 'grafana_com'
  | 'azuread'
  | 'okta';

/** Map of enabled OAuth services and their respective names
 *
 * @public
 */
export type OAuthSettings = Partial<Record<OAuth, { name: string; icon?: IconName }>>;

/**
 * Information needed for analytics providers
 *
 * @internal
 */
export interface AnalyticsSettings {
  identifier: string;
  intercomIdentifier?: string;
}

/** Current user info included in bootData
 *
 * @internal
 */
export interface CurrentUserDTO {
  isSignedIn: boolean;
  id: number;
  uid: string;
  externalUserId: string;
  login: string;
  email: string;
  name: string;
  theme: string; // dark | light | system
  orgCount: number;
  orgId: number;
  orgName: string;
  orgRole: OrgRole | '';
  isGrafanaAdmin: boolean;
  gravatarUrl: string;
  timezone: string;
  weekStart: string;
  locale: string;
  language: string;
  permissions?: Record<string, boolean>;
  analytics: AnalyticsSettings;
  authenticatedBy: string;

  /** @deprecated Use theme instead */
  lightTheme: boolean;
}

/** Contains essential user and config info
 *
 * @internal
 */
export interface BootData {
  user: CurrentUserDTO;
  settings: GrafanaConfig;
  navTree: NavLinkDTO[];
  assets: {
    light: string;
    dark: string;
    WideSky: string;
  };
}

/**
 * Describes all the different Grafana configuration values available for an instance.
 *
 * @internal
 */
export interface GrafanaConfig {
  publicDashboardAccessToken?: string;
  publicDashboardsEnabled: boolean;
  snapshotEnabled: boolean;
  datasources: { [str: string]: DataSourceInstanceSettings };
  panels: { [key: string]: PanelPluginMeta };
  auth: AuthSettings;
  minRefreshInterval: string;
  appSubUrl: string;
  windowTitlePrefix: string;
  buildInfo: BuildInfo;
  newPanelTitle: string;
  bootData: BootData;
  externalUserMngLinkUrl: string;
  externalUserMngLinkName: string;
  externalUserMngInfo: string;
  allowOrgCreate: boolean;
  disableLoginForm: boolean;
  defaultDatasource: string;
  authProxyEnabled: boolean;
  exploreEnabled: boolean;
  queryHistoryEnabled: boolean;
  helpEnabled: boolean;
  profileEnabled: boolean;
  newsFeedEnabled: boolean;
  ldapEnabled: boolean;
  sigV4AuthEnabled: boolean;
  azureAuthEnabled: boolean;
  samlEnabled: boolean;
  autoAssignOrg: boolean;
  verifyEmailEnabled: boolean;
  oauth: OAuthSettings;
  /** @deprecated always set to true. */
  rbacEnabled: boolean;
  disableUserSignUp: boolean;
  loginHint: string;
  passwordHint: string;
  loginError?: string;
  viewersCanEdit: boolean;
  editorsCanAdmin: boolean;
  disableSanitizeHtml: boolean;
  trustedTypesDefaultPolicyEnabled: boolean;
  cspReportOnlyEnabled: boolean;
  liveEnabled: boolean;
  /** @deprecated Use `theme2` instead. */
  theme: GrafanaTheme;
  theme2: GrafanaTheme2;
  anonymousEnabled: boolean;
  anonymousDeviceLimit: number | undefined;
  featureToggles: FeatureToggles;
  licenseInfo: LicenseInfo;
  http2Enabled: boolean;
  dateFormats?: SystemDateFormatSettings;
  grafanaJavascriptAgent: GrafanaJavascriptAgentConfig;
  customTheme?: any;
  geomapDefaultBaseLayer?: MapLayerOptions;
  geomapDisableCustomBaseLayer?: boolean;
  unifiedAlertingEnabled: boolean;
  unifiedAlerting: UnifiedAlertingConfig;
  angularSupportEnabled: boolean;
  feedbackLinksEnabled: boolean;
  secretsManagerPluginEnabled: boolean;
  supportBundlesEnabled: boolean;
  secureSocksDSProxyEnabled: boolean;
  googleAnalyticsId: string | undefined;
  googleAnalytics4Id: string | undefined;
  googleAnalytics4SendManualPageViews: boolean;
  rudderstackWriteKey: string | undefined;
  rudderstackDataPlaneUrl: string | undefined;
  rudderstackSdkUrl: string | undefined;
  rudderstackConfigUrl: string | undefined;
  rudderstackIntegrationsUrl: string | undefined;
  sqlConnectionLimits: SqlConnectionLimits;
  sharedWithMeFolderUID?: string;
  rootFolderUID?: string;
  localFileSystemAvailable?: boolean;
  cloudMigrationIsTarget?: boolean;
  listDashboardScopesEndpoint?: string;
  listScopesEndpoint?: string;

  // The namespace to use for kubernetes apiserver requests
  namespace: string;

  /**
   * Language used in Grafana's UI. This is after the user's preference (or deteceted locale) is resolved to one of
   * Grafana's supported language.
   */
  language: string | undefined;

  // WideSky
  wideSkyWhitelabeling?: WideSkyWhitelabeling;
  wideSkyTheme?: WideSkyCustomTheme;
  wideSkyProvisioner?: WideSkyProvisioner;
}

export interface SqlConnectionLimits {
  maxOpenConns: number;
  maxIdleConns: number;
  connMaxLifetime: number;
}

export interface AuthSettings {
  AuthProxyEnableLoginToken?: boolean;
  // @deprecated -- this is no longer used and will be removed in Grafana 11
  OAuthSkipOrgRoleUpdateSync?: boolean;
  // @deprecated -- this is no longer used and will be removed in Grafana 11
  SAMLSkipOrgRoleSync?: boolean;
  // @deprecated -- this is no longer used and will be removed in Grafana 11
  LDAPSkipOrgRoleSync?: boolean;
  // @deprecated -- this is no longer used and will be removed in Grafana 11
  JWTAuthSkipOrgRoleSync?: boolean;
  // @deprecated -- this is no longer used and will be removed in Grafana 11
  GrafanaComSkipOrgRoleSync?: boolean;
  // @deprecated -- this is no longer used and will be removed in Grafana 11
  GithubSkipOrgRoleSync?: boolean;
  // @deprecated -- this is no longer used and will be removed in Grafana 11
  GitLabSkipOrgRoleSync?: boolean;
  // @deprecated -- this is no longer used and will be removed in Grafana 11
  OktaSkipOrgRoleSync?: boolean;
  // @deprecated -- this is no longer used and will be removed in Grafana 11
  AzureADSkipOrgRoleSync?: boolean;
  // @deprecated -- this is no longer used and will be removed in Grafana 11
  GoogleSkipOrgRoleSync?: boolean;
  // @deprecated -- this is no longer used and will be removed in Grafana 11
  GenericOAuthSkipOrgRoleSync?: boolean;

  disableLogin?: boolean;
  basicAuthStrongPasswordPolicy?: boolean;
}

// WideSky
interface HelpOption {
  label: string;
  href: string;
}

interface FooterItem {
  text: string;
  url: string;
  icon: IconName;
  hide: boolean;
  color?: string;
}

interface LoginBox {
  logo: string;
  logoMaxWidth: number;
  logoMaxWidthMediaBreakpoint: number;

  color?: string;
  textColor?: string;
  textboxBorderColor?: string;
  textboxPlaceholderColor?: string;
  textboxBackgroundColor?: string;

  buttonBgColor?: string;
  buttonHoverBgColor?: string;
  buttonTextColor?: string;
  buttonHoverTextColor?: string;

  linkButtonBgColor?: string;
  linkButtonHoverBgColor?: string;
  linkButtonTextColor?: string;
  linkButtonHoverTextColor?: string;
}

interface LoginBackground {
  image: string;
  background?: string;
  position: string;
  color?: string;
  minHeight: string;
}

export interface WideSkyWhitelabeling {
  appSidebarLogo: string;
  applicationName: string;
  applicationLogo: string;
  browserTabTitle: string;
  browserTabSubtitle: string;
  showEnterpriseUpgradeInfo: boolean;
  adminSettingsMessage: string;

  footerItems: FooterItem[];
  footerPipeColor?: string;

  helpOptions: HelpOption[];

  loginBox: LoginBox;
  loginBackground: LoginBackground;

  entityNotFoundText: string;
  entityNotFoundLink: string;
  entityNotFoundLinkText: string;
  templateVariableHelpLink: string;
}

interface LanguagePermissions {
    enabled: boolean;
    languages: { [key: string]: string };
}

export interface WideSkyProvisioner {
  languagePermissions: LanguagePermissions;
}
