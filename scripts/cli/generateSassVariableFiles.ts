import fs from 'fs';
import { writeFile } from 'node:fs/promises';
import path, { resolve } from 'path';
import * as sass from 'sass';

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
    console.log('Created custom SASS WideSky file');
  } catch (error) {
    console.log(`Unable to create custom SASS WideSky theme due to: ${error.message}`);
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

const resolveAndInlineImports = (scssContent: string, baseDir: string): string => {
  return scssContent.replace(/@import ['"]([^'"]+\.(css|scss))['"];/g, (_, importPath, ext) => {
    const resolvedPath = require.resolve(importPath, {
      paths: [baseDir, path.resolve('node_modules')],
    });

    if (!fs.existsSync(resolvedPath)) {
      throw new Error(`Cannot find imported file: ${importPath}`);
    }

    const importedContent = fs.readFileSync(resolvedPath, 'utf8');

    if (ext === 'css') {
      // For .css files, directly inline the content without further processing
      return `/* Inlined from ${importPath} */\n${importedContent}\n\n`;
    } else if (ext === 'scss') {
      // For .scss files, resolve imports recursively
      return `/* Inlined from ${importPath} */\n${resolveAndInlineImports(
        importedContent,
        path.dirname(resolvedPath)
      )}`;
    }

    return `@import '${importPath}';`; // Keep as-is for unexpected cases
  });
};

const compileStyles = async () => {
  const input = resolve(__dirname, 'public', 'sass', 'grafana.WideSky.scss');

  try {
    const buildDir = './public/build'; // Directory containing the hashed output file
    const outputFileNamePattern = /^grafana\.WideSky\.[a-f0-9]+\.css$/;

    console.log(`Processing ${input}...`);

    // Compile SCSS to CSS, resolving imports from node_modules
    const result = sass.compile(input, {
      sourceMap: false,
      charset: false,
      loadPaths: [path.resolve('node_modules')],
    });

    const inlinedContent = resolveAndInlineImports(result.css, path.dirname(input));
    const result2 = sass.compileString(inlinedContent, {
      sourceMap: false,
      charset: false,
      style: 'compressed',
      loadPaths: [path.resolve('node_modules')],
    });

    // Locate the existing file in the build directory
    const files = fs.readdirSync(buildDir);
    const targetFile = files.find((file) => outputFileNamePattern.test(file));

    if (!targetFile) {
      throw new Error(`No file matching the pattern found in ${buildDir}`);
    }

    const targetFilePath = path.join(buildDir, targetFile);

    // Overwrite the existing file
    fs.writeFileSync(targetFilePath, result2.css);
    console.log(`Output written to ${targetFilePath}`);
  } catch (error) {
    console.error(`Error processing ${input}:`, error.message);
    console.error(error.stack);
  }
};

generateSassVariableFiles()
  .then(compileStyles)
  .catch((error) => {
    console.error('Build process failed:', error.message);
    process.exit(1);
  });
