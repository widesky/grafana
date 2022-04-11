import { merge } from 'lodash';
import {
  BootData,
  BuildInfo,
  createTheme,
  DataSourceInstanceSettings,
  FeatureToggles,
  GrafanaConfig,
  GrafanaTheme,
  GrafanaTheme2,
  LicenseInfo,
  MapLayerOptions,
  OAuthSettings,
  PanelPluginMeta,
  PreloadPlugin,
  systemDateFormats,
  SystemDateFormatSettings,
} from '@grafana/data';

export interface AzureSettings {
  cloud?: string;
  managedIdentityEnabled: boolean;
}

export class GrafanaBootConfig implements GrafanaConfig {
  datasources: { [str: string]: DataSourceInstanceSettings } = {};
  panels: { [key: string]: PanelPluginMeta } = {};
  minRefreshInterval = '';
  appUrl = '';
  appSubUrl = '';
  windowTitlePrefix = '';
  buildInfo: BuildInfo;
  newPanelTitle = '';
  bootData: BootData;
  externalUserMngLinkUrl = '';
  externalUserMngLinkName = '';
  externalUserMngInfo = '';
  allowOrgCreate = false;
  disableLoginForm = false;
  defaultDatasource = ''; // UID
  alertingEnabled = false;
  alertingErrorOrTimeout = '';
  alertingNoDataOrNullValues = '';
  alertingMinInterval = 1;
  angularSupportEnabled = false;
  authProxyEnabled = false;
  exploreEnabled = false;
  helpEnabled = false;
  profileEnabled = false;
  ldapEnabled = false;
  sigV4AuthEnabled = false;
  samlEnabled = false;
  samlName = '';
  autoAssignOrg = true;
  verifyEmailEnabled = false;
  oauth: OAuthSettings = {};
  disableUserSignUp = false;
  loginHint = '';
  passwordHint = '';
  loginError = undefined;
  navTree: any;
  viewersCanEdit = false;
  editorsCanAdmin = false;
  disableSanitizeHtml = false;
  liveEnabled = true;
  theme: GrafanaTheme;
  theme2: GrafanaTheme2;
  pluginsToPreload: PreloadPlugin[] = [];
  featureToggles: FeatureToggles = {};
  licenseInfo: LicenseInfo = {} as LicenseInfo;
  rendererAvailable = false;
  dashboardPreviews: {
    systemRequirements: {
      met: boolean;
      requiredImageRendererPluginVersion: string;
    };
    thumbnailsExist: boolean;
  } = { systemRequirements: { met: false, requiredImageRendererPluginVersion: '' }, thumbnailsExist: false };
  rendererVersion = '';
  http2Enabled = false;
  dateFormats?: SystemDateFormatSettings;
  sentry = {
    enabled: false,
    dsn: '',
    customEndpoint: '',
    sampleRate: 1,
  };
  pluginCatalogURL = 'https://grafana.com/grafana/plugins/';
  pluginAdminEnabled = true;
  pluginAdminExternalManageEnabled = false;
  pluginCatalogHiddenPlugins: string[] = [];
  expressionsEnabled = false;
  customTheme?: any;
  awsAllowedAuthProviders: string[] = [];
  awsAssumeRoleEnabled = false;
  azure: AzureSettings = {
    managedIdentityEnabled: false,
  };
  caching = {
    enabled: false,
  };
  geomapDefaultBaseLayerConfig?: MapLayerOptions;
  geomapDisableCustomBaseLayer?: boolean;
  unifiedAlertingEnabled = false;
  applicationInsightsConnectionString?: string;
  applicationInsightsEndpointUrl?: string;
  recordedQueries = {
    enabled: true,
  };
  featureHighlights = {
    enabled: false,
  };
  reporting = {
    enabled: true,
  };

  // WideSky Grafana white labeling
  footerElement1Icon = '';
  footerElement2Icon = '';
  footerElement3Icon = '';
  footerElement1Text = '';
  footerElement2Text = '';
  footerElement3Text = '';
  footerElement1Link = '';
  footerElement2Link = '';
  footerElement3Link = '';
  footerElement1TextColor = '';
  footerElement2TextColor = '';
  footerElement3TextColor = '';
  footerPipeColor = '';
  browserTabTitle = '';
  appSidebarLogo = '';
  applicationName = '';
  loginBackground = '';
  loginBoxLogo = '';
  loginBoxLogoMaxWidth = '';
  loginBoxLogoMaxWidthMediaBreakpoint = '';
  loginBoxColour = '';
  loginBoxTextColour = '';
  loginBoxTextboxPlaceholderColour = '';
  loginBoxTextboxBorderColour = '';
  loginBoxButtonBgColour = '';
  loginBoxButtonHoverBgColour = '';
  loginBoxButtonTextColour = '';
  loginBoxButtonHoverTextColour = '';
  loginBoxLinkButtonBgColour = '';
  loginBoxLinkButtonHoverBgColour = '';
  loginBoxLinkButtonTextColour = '';
  loginBoxLinkButtonHoverTextColour = '';

  constructor(options: GrafanaBootConfig) {
    const mode = options.bootData.user.lightTheme ? 'light' : 'dark';
    this.theme2 = createTheme({ colors: { mode } });
    this.theme = this.theme2.v1;
    this.bootData = options.bootData;
    this.buildInfo = options.buildInfo;

    const defaults = {
      datasources: {},
      windowTitlePrefix: 'Widesky - ',
      panels: {},
      newPanelTitle: 'Panel Title',
      playlist_timespan: '1m',
      unsaved_changes_warning: true,
      appUrl: '',
      appSubUrl: '',
      buildInfo: {
        version: 'v1.0',
        commit: '1',
        env: 'production',
      },
      viewersCanEdit: false,
      editorsCanAdmin: false,
      disableSanitizeHtml: false,
      hideVersion: true,
    };

    merge(this, defaults, options);

    if (this.dateFormats) {
      systemDateFormats.update(this.dateFormats);
    }
  }
}

const bootData = (window as any).grafanaBootData || {
  settings: {},
  user: {},
  navTree: [],
};

const options = bootData.settings;
options.bootData = bootData;

/**
 * Use this to access the {@link GrafanaBootConfig} for the current running Grafana instance.
 *
 * @public
 */
export const config = new GrafanaBootConfig(options);
