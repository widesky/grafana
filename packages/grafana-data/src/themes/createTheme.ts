import { createBreakpoints } from './breakpoints';
import { createColors, ThemeColorsInput } from './createColors';
import { createComponents } from './createComponents';
import { createShadows } from './createShadows';
import { createShape, ThemeShapeInput } from './createShape';
import { createSpacing, ThemeSpacingOptions } from './createSpacing';
import { createTransitions } from './createTransitions';
import { createTypography, ThemeTypographyInput } from './createTypography';
import { createV1Theme } from './createV1Theme';
import { createVisualizationColors } from './createVisualizationColors';
import { GrafanaTheme2, WideSkyCustomTheme } from './types';
import { zIndex } from './zIndex';

/** @internal */
export interface NewThemeOptions {
  name?: string;
  colors?: ThemeColorsInput;
  spacing?: ThemeSpacingOptions;
  shape?: ThemeShapeInput;
  typography?: ThemeTypographyInput;
}

/** @internal */
export function createTheme(options: NewThemeOptions = {}): GrafanaTheme2 {
  const {
    colors: colorsInput = {},
    spacing: spacingInput = {},
    shape: shapeInput = {},
    typography: typographyInput = {},
  } = options;

  const colors = createColors(colorsInput);
  const breakpoints = createBreakpoints();
  const spacing = createSpacing(spacingInput);
  const shape = createShape(shapeInput);
  const typography = createTypography(colors, typographyInput);
  const shadows = createShadows(colors);
  const transitions = createTransitions();
  const components = createComponents(colors, shadows);
  const visualization = createVisualizationColors(colors);

  const theme = {
    name: colors.mode[0].toUpperCase() + colors.mode.substring(1),
    isDark: colors.mode === 'dark',
    isLight: colors.mode === 'light',
    isWideSky: colors.mode === 'WideSky',
    colors,
    breakpoints,
    spacing,
    shape,
    components,
    typography,
    shadows,
    transitions,
    visualization,
    zIndex: {
      ...zIndex,
    },
    flags: {},
  };

  return {
    ...theme,
    v1: createV1Theme(theme),
  };
}

export class CustomColors {
  public static wideSkyTheme?: WideSkyCustomTheme;
  static wideSkyThemePromise?: Promise<void>;

  public static initCustomTheme = async () => {
    if (CustomColors.wideSkyTheme !== undefined) {
      return;
    }

    if (CustomColors.wideSkyThemePromise === undefined) {
      CustomColors.wideSkyThemePromise = CustomColors.requestWideSkyTheme();
    }

    await CustomColors.wideSkyThemePromise;
  };

  private static requestWideSkyTheme = async () => {
    let processEnv: NodeJS.ProcessEnv | undefined = undefined;

    if (typeof window !== 'undefined') {
      processEnv = window.process?.env;
    } else if (typeof process !== 'undefined' && process.versions != null && process.versions.node != null) {
      processEnv = process.env;
    } else {
      console.error('Unknown environment, using default auth and backend url:port');
    }

    const basicAuth = processEnv?.BASIC_AUTH ?? 'admin:admin';
    const url = `${processEnv?.GF_BACKEND ?? 'http://localhost:3000'}/api/frontend/settings`;

    const resp = await fetch(url, {
      method: 'GET',
      headers: { Authorization: 'Basic ' + btoa(basicAuth) },
    }).catch(() => {
      throw Error(`Failed to fetch data as host ${url} is incorrect or unavailable`);
    });

    const json = await resp.json();

    if (json.wideSkyTheme === undefined) {
      console.error('WideSky custom theme is not present in frontend settings');
      return;
    }

    CustomColors.wideSkyTheme = json.wideSkyTheme as WideSkyCustomTheme;
  };
}
