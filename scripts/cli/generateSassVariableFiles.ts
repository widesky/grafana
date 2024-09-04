import { writeFile } from 'node:fs/promises';
import { resolve } from 'path';

import { GrafanaTheme2, createTheme } from '@grafana/data';
import { CustomColors } from '@grafana/data/src/themes/createTheme';
import { wideskyThemeVarsTemplate } from '@grafana/ui/src/themes/_variables.WideSky.scss.tmpl';
import { darkThemeVarsTemplate } from '@grafana/ui/src/themes/_variables.dark.scss.tmpl';
import { lightThemeVarsTemplate } from '@grafana/ui/src/themes/_variables.light.scss.tmpl';
import { commonThemeVarsTemplate } from '@grafana/ui/src/themes/_variables.scss.tmpl';

const darkThemeVariablesPath = resolve(__dirname, 'public', 'sass', '_variables.dark.generated.scss');
const lightThemeVariablesPath = resolve(__dirname, 'public', 'sass', '_variables.light.generated.scss');
const defaultThemeVariablesPath = resolve(__dirname, 'public', 'sass', '_variables.generated.scss');
const wideskyThemeVariablesPath = resolve(__dirname, 'public', 'sass', '_variables.WideSky.generated.scss');

async function writeVariablesFile(path: string, data: string) {
  try {
    await writeFile(path, data);
  } catch (error) {
    console.error('\nWriting SASS variable files failed', error);
    process.exit(1);
  }
}

async function generateSassVariableFiles() {
  const darkTheme = createTheme();
  const lightTheme = createTheme({ colors: { mode: 'light' } });

  let wideSkyTheme: GrafanaTheme2 | undefined;
  try {
    await CustomColors.initCustomTheme();
    wideSkyTheme = createTheme({ colors: { mode: 'WideSky' } });
    console.log('Created a custom WideSky theme');
  } catch (error) {
    console.log(`Unable to create custom WideSky theme due to: ${error.message}`);
  }

  try {
    await writeVariablesFile(darkThemeVariablesPath, darkThemeVarsTemplate(darkTheme));
    await writeVariablesFile(lightThemeVariablesPath, lightThemeVarsTemplate(lightTheme));
    await writeVariablesFile(defaultThemeVariablesPath, commonThemeVarsTemplate(darkTheme));
    if (wideSkyTheme) {
      await writeVariablesFile(wideskyThemeVariablesPath, wideskyThemeVarsTemplate(wideSkyTheme));
    }
  } catch (error) {
    console.error('\nWriting SASS variable files failed', error);
    process.exit(1);
  }
}

generateSassVariableFiles();
