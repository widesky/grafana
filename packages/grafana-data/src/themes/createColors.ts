import { merge } from 'lodash';

import { alpha, darken, emphasize, getContrastRatio, lighten } from './colorManipulator';
import { CustomColors } from './createTheme';
import { palette } from './palette';
import { BrandGradient, DeepPartial, ThemeRichColor } from './types';

/** @internal */
export type ThemeColorsMode = 'light' | 'dark' | 'WideSky';

/** @internal */
export interface ThemeColorsBase<TColor> {
  mode: ThemeColorsMode;

  primary: TColor;
  secondary: TColor;
  info: TColor;
  error: TColor;
  success: TColor;
  warning: TColor;

  text: {
    primary: string;
    secondary: string;
    disabled: string;
    link: string;
    /** Used for auto white or dark text on colored backgrounds */
    maxContrast: string;
  };

  background: {
    /** Dashboard and body background */
    canvas: string;
    /** Primary content pane background (panels etc) */
    primary: string;
    /** Cards and elements that need to stand out on the primary background */
    secondary: string;
  };

  border: {
    weak: string;
    medium: string;
    strong: string;
  };

  gradients: {
    brandVertical: string;
    brandHorizontal: string;
  };

  action: {
    /** Used for selected menu item / select option */
    selected: string;
    /**
     * @alpha (Do not use from plugins)
     * Used for selected items when background only change is not enough (Currently only used for FilterPill)
     **/
    selectedBorder: string;
    /** Used for hovered menu item / select option */
    hover: string;
    /** Used for button/colored background hover opacity */
    hoverOpacity: number;
    /** Used focused menu item / select option */
    focus: string;
    /** Used for disabled buttons and inputs */
    disabledBackground: string;
    /** Disabled text */
    disabledText: string;
    /** Disablerd opacity */
    disabledOpacity: number;
  };

  hoverFactor: number;
  contrastThreshold: number;
  tonalOffset: number;
}

export interface ThemeHoverStrengh {}

/** @beta */
export interface ThemeColors extends ThemeColorsBase<ThemeRichColor> {
  /** Returns a text color for the background */
  getContrastText(background: string, threshold?: number): string;
  /* Brighten or darken a color by specified factor (0-1) */
  emphasize(color: string, amount?: number): string;
}

/** @internal */
export type ThemeColorsInput = DeepPartial<ThemeColorsBase<ThemeRichColor>>;

class DarkColors implements ThemeColorsBase<Partial<ThemeRichColor>> {
  mode: ThemeColorsMode = 'dark';

  // Used to get more white opacity colors
  whiteBase = '204, 204, 220';

  border = {
    weak: `rgba(${this.whiteBase}, 0.12)`,
    medium: `rgba(${this.whiteBase}, 0.20)`,
    strong: `rgba(${this.whiteBase}, 0.30)`,
  };

  text = {
    primary: `rgb(${this.whiteBase})`,
    secondary: `rgba(${this.whiteBase}, 0.65)`,
    disabled: `rgba(${this.whiteBase}, 0.6)`,
    link: palette.blueDarkText,
    maxContrast: palette.white,
  };

  primary = {
    main: palette.blueDarkMain,
    text: palette.blueDarkText,
    border: palette.blueDarkText,
  };

  secondary = {
    main: `rgba(${this.whiteBase}, 0.10)`,
    shade: `rgba(${this.whiteBase}, 0.14)`,
    transparent: `rgba(${this.whiteBase}, 0.08)`,
    text: this.text.primary,
    contrastText: `rgb(${this.whiteBase})`,
    border: `rgba(${this.whiteBase}, 0.08)`,
  };

  info = this.primary;

  error = {
    main: palette.redDarkMain,
    text: palette.redDarkText,
  };

  success = {
    main: palette.greenDarkMain,
    text: palette.greenDarkText,
  };

  warning = {
    main: palette.orangeDarkMain,
    text: palette.orangeDarkText,
  };

  background = {
    canvas: palette.gray05,
    primary: palette.gray10,
    secondary: palette.gray15,
  };

  action = {
    hover: `rgba(${this.whiteBase}, 0.16)`,
    selected: `rgba(${this.whiteBase}, 0.12)`,
    selectedBorder: palette.orangeDarkMain,
    focus: `rgba(${this.whiteBase}, 0.16)`,
    hoverOpacity: 0.08,
    disabledText: this.text.disabled,
    disabledBackground: `rgba(${this.whiteBase}, 0.04)`,
    disabledOpacity: 0.38,
  };

  gradients = {
    brandHorizontal: 'linear-gradient(270deg, #4e9ed7 0%, #0073bc 100%)',
    brandVertical: 'linear-gradient(0.01deg, #4e9ed7 0.01%, #0073bc 99.99%)',
  };

  contrastThreshold = 3;
  hoverFactor = 0.03;
  tonalOffset = 0.15;
}

class LightColors implements ThemeColorsBase<Partial<ThemeRichColor>> {
  mode: ThemeColorsMode = 'light';

  blackBase = '36, 41, 46';

  primary = {
    main: palette.blueLightMain,
    border: palette.blueLightText,
    text: palette.blueLightText,
  };

  text = {
    primary: `rgba(${this.blackBase}, 1)`,
    secondary: `rgba(${this.blackBase}, 0.75)`,
    disabled: `rgba(${this.blackBase}, 0.50)`,
    link: this.primary.text,
    maxContrast: palette.black,
  };

  border = {
    weak: `rgba(${this.blackBase}, 0.12)`,
    medium: `rgba(${this.blackBase}, 0.30)`,
    strong: `rgba(${this.blackBase}, 0.40)`,
  };

  secondary = {
    main: `rgba(${this.blackBase}, 0.08)`,
    shade: `rgba(${this.blackBase}, 0.15)`,
    transparent: `rgba(${this.blackBase}, 0.08)`,
    contrastText: `rgba(${this.blackBase},  1)`,
    text: this.text.primary,
    border: this.border.weak,
  };

  info = {
    main: palette.blueLightMain,
    text: palette.blueLightText,
  };

  error = {
    main: palette.redLightMain,
    text: palette.redLightText,
    border: palette.redLightText,
  };

  success = {
    main: palette.greenLightMain,
    text: palette.greenLightText,
  };

  warning = {
    main: palette.orangeLightMain,
    text: palette.orangeLightText,
  };

  background = {
    canvas: palette.gray90,
    primary: palette.white,
    secondary: palette.gray100,
  };

  action = {
    hover: `rgba(${this.blackBase}, 0.12)`,
    selected: `rgba(${this.blackBase}, 0.08)`,
    selectedBorder: palette.orangeLightMain,
    hoverOpacity: 0.08,
    focus: `rgba(${this.blackBase}, 0.12)`,
    disabledBackground: `rgba(${this.blackBase}, 0.04)`,
    disabledText: this.text.disabled,
    disabledOpacity: 0.38,
  };

  gradients = {
    brandHorizontal: 'linear-gradient(90deg, #0073bc 0%, #4e9ed7 100%);',
    brandVertical: 'linear-gradient(0.01deg, #4e9ed7 -31.2%, #0073bc 113.07%);',
  };

  contrastThreshold = 3;
  hoverFactor = 0.03;
  tonalOffset = 0.2;
}

export class WideSkyColors implements ThemeColorsBase<Partial<ThemeRichColor>> {
  mode: ThemeColorsMode = 'WideSky';

  background = {
    canvas: CustomColors.wideSkyTheme!.backgroundCanvas,
    primary: CustomColors.wideSkyTheme!.backgroundPrimary,
    secondary: CustomColors.wideSkyTheme!.backgroundSecondary,
  };

  gradients = {
    brandHorizontal: this.getGradient('linear-gradient(90deg, $0 0%, $1 100%);', CustomColors.wideSkyTheme!.horizontal),
    brandVertical: this.getGradient(
      'linear-gradient(0.01deg, $0 -31.2%, $1 113.07%);',
      CustomColors.wideSkyTheme!.vertical
    ),
  };

  primary = {
    main: CustomColors.wideSkyTheme!.primary.main,
    shade: CustomColors.wideSkyTheme!.primary.shade,
    text: CustomColors.wideSkyTheme!.primary.text,
    contrastText: CustomColors.wideSkyTheme!.primary.text,
  };
  secondary = {
    main: this.fromBaseOrValue(CustomColors.wideSkyTheme!.secondary.main, 0.16),
    shade: this.fromBaseOrValue(CustomColors.wideSkyTheme!.secondary.shade, 0.2),
    text: this.fromBaseOrValue(CustomColors.wideSkyTheme!.secondary.text, 1),
    contrastText: this.fromBaseOrValue(CustomColors.wideSkyTheme!.secondary.text, 1),
  };

  info = {
    main: CustomColors.wideSkyTheme!.info.main,
    shade: CustomColors.wideSkyTheme!.info.shade,
    contrastText: CustomColors.wideSkyTheme!.info.text,
  };

  error = {
    main: CustomColors.wideSkyTheme!.error.main,
    shade: CustomColors.wideSkyTheme!.error.shade,
    contrastText: CustomColors.wideSkyTheme!.error.text,
  };

  success = {
    main: CustomColors.wideSkyTheme!.success.main,
    shade: CustomColors.wideSkyTheme!.success.shade,
    contrastText: CustomColors.wideSkyTheme!.success.text,
  };

  warning = {
    main: CustomColors.wideSkyTheme!.warning.main,
    shade: CustomColors.wideSkyTheme!.warning.shade,
    contrastText: CustomColors.wideSkyTheme!.warning.text,
  };

  text = {
    primary: this.fromBaseOrValue(CustomColors.wideSkyTheme!.textPrimary, 1),
    secondary: this.fromBaseOrValue(CustomColors.wideSkyTheme!.textSecondary, 0.75),
    disabled: this.fromBaseOrValue(CustomColors.wideSkyTheme!.textDisabled, 0.5),
    link: this.asValueOrDefault(
      CustomColors.wideSkyTheme!.textLink,
      this.fromBaseOrValue(CustomColors.wideSkyTheme!.textPrimary, 1)
    ),
    maxContrast: palette.white,
  };

  border = {
    weak: this.fromBorderOrValue(CustomColors.wideSkyTheme!.borderWeak, 0.12),
    medium: this.fromBorderOrValue(CustomColors.wideSkyTheme!.borderMedium, 0.3),
    strong: this.fromBorderOrValue(CustomColors.wideSkyTheme!.borderStrong, 0.4),
  };

  action = {
    hover: this.fromBaseOrValue(CustomColors.wideSkyTheme!.actionHover, 0.12),
    selected: this.fromBaseOrValue(CustomColors.wideSkyTheme!.actionSelected, 0.08),
    focus: this.fromBaseOrValue(CustomColors.wideSkyTheme!.actionFocus, 0.12),
    disabledBackground: this.fromBaseOrValue(CustomColors.wideSkyTheme!.actionDisabledBackground, 0.04),
    disabledText: this.text.disabled,
    hoverOpacity: 0.08,
    disabledOpacity: 0.38,
    selectedBorder: palette.orangeDarkMain,
  };

  contrastThreshold = 3;
  hoverFactor = 0.03;
  tonalOffset = 0.15;

  static WS_NOT_SET = 'NOT_SET';

  /**
   * Apply a gradient CSS. If all is not set, comprise gradient from the template. Otherwise use the all string.
   * @param template Template to add colorA and colorB to.
   * @param colorA Color for $0 in template.
   * @param colorB Color for $1 in template.
   * @param all A complete gradient CSS.
   */
  getGradient(template: string, gradient: BrandGradient) {
    if (gradient.all === WideSkyColors.WS_NOT_SET) {
      // just customise the colors
      template = template.replace('$0', gradient.colorA);
      template = template.replace('$1', gradient.colorB);

      return template;
    } else {
      // allow full user customisation
      return gradient.all;
    }
  }

  fromBaseOrValue(value: string, opacity: number) {
    return CustomColors.wideSkyTheme!.baseColor === WideSkyColors.WS_NOT_SET
      ? value
      : `rgba(${CustomColors.wideSkyTheme!.baseColor}, ${opacity})`;
  }

  fromBorderOrValue(value: string, opacity: number) {
    return CustomColors.wideSkyTheme!.borderColor === WideSkyColors.WS_NOT_SET
      ? value
      : `rgba(${CustomColors.wideSkyTheme!.borderColor}, ${opacity})`;
  }

  asValueOrDefault(value: string, defaultVal: string) {
    return value === WideSkyColors.WS_NOT_SET ? defaultVal : value;
  }
}

export function createColors(colors: ThemeColorsInput): ThemeColors {
  const dark = new DarkColors();
  const light = new LightColors();
  let base: ThemeColorsBase<Partial<ThemeRichColor>> = (colors.mode ?? 'dark') === 'dark' ? dark : light;
  if (CustomColors.wideSkyTheme !== undefined && colors.mode === 'WideSky') {
    base = new WideSkyColors();
  }
  const {
    primary = base.primary,
    secondary = base.secondary,
    info = base.info,
    warning = base.warning,
    success = base.success,
    error = base.error,
    tonalOffset = base.tonalOffset,
    hoverFactor = base.hoverFactor,
    contrastThreshold = base.contrastThreshold,
    ...other
  } = colors;

  function getContrastText(background: string, threshold: number = contrastThreshold) {
    const contrastText =
      getContrastRatio(dark.text.maxContrast, background, base.background.primary) >= threshold
        ? dark.text.maxContrast
        : light.text.maxContrast;
    // todo, need color framework
    return contrastText;
  }

  const getRichColor = ({ color, name }: GetRichColorProps): ThemeRichColor => {
    color = { ...color, name };
    if (!color.main) {
      throw new Error(`Missing main color for ${name}`);
    }
    if (!color.text) {
      color.text = color.main;
    }
    if (!color.border) {
      color.border = color.text;
    }
    if (!color.shade) {
      color.shade = base.mode === 'light' ? darken(color.main, tonalOffset) : lighten(color.main, tonalOffset);
    }
    if (!color.transparent) {
      color.transparent = alpha(color.main, 0.15);
    }
    if (!color.contrastText) {
      color.contrastText = getContrastText(color.main);
    }
    if (!color.borderTransparent) {
      color.borderTransparent = alpha(color.border, 0.25);
    }
    return color as ThemeRichColor;
  };

  return merge(
    {
      ...base,
      primary: getRichColor({ color: primary, name: 'primary' }),
      secondary: getRichColor({ color: secondary, name: 'secondary' }),
      info: getRichColor({ color: info, name: 'info' }),
      error: getRichColor({ color: error, name: 'error' }),
      success: getRichColor({ color: success, name: 'success' }),
      warning: getRichColor({ color: warning, name: 'warning' }),
      getContrastText,
      emphasize: (color: string, factor?: number) => {
        return emphasize(color, factor ?? hoverFactor);
      },
    },
    other
  );
}

interface GetRichColorProps {
  color: Partial<ThemeRichColor>;
  name: string;
}
