import { GrafanaTheme } from '../types/theme';

import { ThemeBreakpoints } from './breakpoints';
import { ThemeColors } from './createColors';
import { ThemeComponents } from './createComponents';
import { ThemeShadows } from './createShadows';
import { ThemeShape } from './createShape';
import { ThemeSpacing } from './createSpacing';
import { ThemeTransitions } from './createTransitions';
import { ThemeTypography } from './createTypography';
import { ThemeVisualizationColors } from './createVisualizationColors';
import { ThemeZIndices } from './zIndex';

/**
 * @beta
 * Next gen theme model introduced in Grafana v8.
 */
export interface GrafanaTheme2 {
  name: string;
  isDark: boolean;
  isLight: boolean;
  isWideSky: boolean;
  colors: ThemeColors;
  breakpoints: ThemeBreakpoints;
  spacing: ThemeSpacing;
  shape: ThemeShape;
  components: ThemeComponents;
  typography: ThemeTypography;
  zIndex: ThemeZIndices;
  shadows: ThemeShadows;
  visualization: ThemeVisualizationColors;
  transitions: ThemeTransitions;
  /** @deprecated Will be removed in a future version */
  v1: GrafanaTheme;
  /** feature flags that might impact component looks */
  flags: {};
}

/** @alpha */
export interface ThemeRichColor {
  /** color intent (primary, secondary, info, error, etc) */
  name: string;
  /** Main color */
  main: string;
  /** Used for hover */
  shade: string;
  /** Used for text */
  text: string;
  /** Used for borders */
  border: string;
  /** Used subtly colored backgrounds */
  transparent: string;
  /** Used for weak colored borders like larger alert/banner boxes and smaller badges and tags */
  borderTransparent: string;
  /** Text color for text ontop of main */
  contrastText: string;
}

/** @internal */
export type DeepPartial<T> = {
  [P in keyof T]?: DeepPartial<T[P]>;
};

export type BrandGradient = {
  colorA: string;
  colorB: string;
  all: string;
};

export type WideSkyCustomTheme = {
  name: string;
  baseColor: string;
  borderColor: string;

  fontHeader: string;
  fontBody: string;

  horizontal: BrandGradient;
  vertical: BrandGradient;

  primary: ThemeRichColor;
  secondary: ThemeRichColor;
  info: ThemeRichColor;
  success: ThemeRichColor;
  warning: ThemeRichColor;
  error: ThemeRichColor;

  textPrimary: string;
  textSecondary: string;
  textDisabled: string;
  textLink: string;

  backgroundCanvas: string;
  backgroundPrimary: string;
  backgroundSecondary: string;
  backgroundNavBar: string;

  borderWeak: string;
  borderMedium: string;
  borderStrong: string;

  actionHover: string;
  actionFocus: string;
  actionSelected: string;
  actionDisabledBackground: string;
};
