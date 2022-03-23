import { DataSourceInstanceSettings } from './datasource';
import { PanelPluginMeta } from './panel';
import { GrafanaTheme } from './theme';
import { SystemDateFormatSettings } from '../datetime';
import { GrafanaTheme2 } from '../themes';
import { MapLayerOptions } from '../geo/layer';
import { FeatureToggles } from './featureToggles.gen';
import { NavLinkDTO, OrgRole } from '.';

/**
 * Describes the build information that will be available via the Grafana configuration.
 *
 * @public
 */
export interface BuildInfo {
  version: string;
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
}

/**
 * Describes Sentry integration config
 *
 * @public
 */
export interface SentryConfig {
  enabled: boolean;
  dsn: string;
  customEndpoint: string;
  sampleRate: number;
}

/**
 * Describes the plugins that should be preloaded prior to start Grafana.
 *
 * @public
 */
export type PreloadPlugin = {
  path: string;
  version: string;
};

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
export type OAuthSettings = Partial<Record<OAuth, { name: string; icon?: string }>>;

/** Current user info included in bootData
 *
 * @internal
 */
export interface CurrentUserDTO {
  isSignedIn: boolean;
  id: number;
  externalUserId: string;
  login: string;
  email: string;
  name: string;
  lightTheme: boolean;
  orgCount: number;
  orgId: number;
  orgName: string;
  orgRole: OrgRole | '';
  isGrafanaAdmin: boolean;
  gravatarUrl: string;
  timezone: string;
  weekStart: string;
  locale: string;
  permissions?: Record<string, boolean>;
}

/** Contains essential user and config info
 *
 * @internal
 */
export interface BootData {
  user: CurrentUserDTO;
  settings: GrafanaConfig;
  navTree: NavLinkDTO[];
  themePaths: {
    light: string;
    dark: string;
  };
}

/**
 * Describes all the different Grafana configuration values available for an instance.
 *
 * @internal
 */
export interface GrafanaConfig {
  datasources: { [str: string]: DataSourceInstanceSettings };
  panels: { [key: string]: PanelPluginMeta };
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
  alertingEnabled: boolean;
  alertingErrorOrTimeout: string;
  alertingNoDataOrNullValues: string;
  alertingMinInterval: number;
  authProxyEnabled: boolean;
  exploreEnabled: boolean;
  helpEnabled: boolean;
  profileEnabled: boolean;
  ldapEnabled: boolean;
  sigV4AuthEnabled: boolean;
  samlEnabled: boolean;
  autoAssignOrg: boolean;
  verifyEmailEnabled: boolean;
  oauth: OAuthSettings;
  disableUserSignUp: boolean;
  loginHint: string;
  passwordHint: string;
  loginError?: string;
  navTree: any;
  viewersCanEdit: boolean;
  editorsCanAdmin: boolean;
  disableSanitizeHtml: boolean;
  liveEnabled: boolean;
  theme: GrafanaTheme;
  theme2: GrafanaTheme2;
  pluginsToPreload: PreloadPlugin[];
  featureToggles: FeatureToggles;
  licenseInfo: LicenseInfo;
  http2Enabled: boolean;
  dateFormats?: SystemDateFormatSettings;
  sentry: SentryConfig;
  customTheme?: any;
  geomapDefaultBaseLayer?: MapLayerOptions;
  geomapDisableCustomBaseLayer?: boolean;
  unifiedAlertingEnabled: boolean;
  angularSupportEnabled: boolean;
  footerElement1Icon?: string;
  footerElement2Icon?: string;
  footerElement3Icon?: string;
  footerElement1Text?: string;
  footerElement2Text?: string;
  footerElement3Text?: string;
  footerElement1Link?: string;
  footerElement2Link?: string;
  footerElement3Link?: string;
  footerElement1TextColor?: string;
  footerElement2TextColor?: string;
  footerElement3TextColor?: string;
  footerPipeColor?: string;
  browserTabTitle?: string;
  appSidebarLogo?: string;
  loginBackground?: string;
  loginBoxLogo?: string;
  loginBoxLogoMaxWidth?: string;
  loginBoxLogoMaxWidthMediaBreakpoint?: string;
  loginBoxColour?: string;
  loginBoxTextColour?: string;
  loginBoxTextboxPlaceholderColour?: string;
  loginBoxButtonBgColour?: string;
  loginBoxButtonHoverBgColour?: string;
  loginBoxButtonTextColour?: string;
  loginBoxButtonHoverTextColour?: string;
  loginBoxLinkButtonBgColour?: string;
  loginBoxLinkButtonHoverBgColour?: string;
  loginBoxLinkButtonTextColour?: string;
  loginBoxLinkButtonHoverTextColour?: string;
}
