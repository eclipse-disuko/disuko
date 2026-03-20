import fs from 'node:fs';
import path from 'node:path';
import {fileURLToPath} from 'node:url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

/**
 * Vite plugin to auto-generate Tailwind CSS color definitions from Colors.ts
 * This ensures a single source of truth for all color definitions
 */
export function generateTailwindColors() {
  const colorsFilePath = path.resolve(__dirname, './src/plugins/Colors.ts');
  const tailwindFilePath = path.resolve(__dirname, './src/styles/tailwind.css');

  const generateCSSColors = (colorsContent) => {
    const cssColors = [];

    const simpleColorRegex = /export const (\w+) = '(#[0-9A-Fa-f]{6})';/g;
    let simpleMatch;

    while ((simpleMatch = simpleColorRegex.exec(colorsContent)) !== null) {
      const colorName = simpleMatch[1];
      const hexValue = simpleMatch[2];
      cssColors.push(`    --color-${colorName}: ${hexValue};`);
    }

    const colorScaleRegex = /export const (\w+) = \{([\s\S]*?)\};/g;
    let match;

    while ((match = colorScaleRegex.exec(colorsContent)) !== null) {
      const colorName = match[1];
      const scaleContent = match[2];

      if (colorName === 'dataVisColors') continue;

      const shadeRegex = /\[(\d+)\]:\s*['"]([#\w]+)['"]/g;
      let shadeMatch;
      const shades = [];

      while ((shadeMatch = shadeRegex.exec(scaleContent)) !== null) {
        shades.push({shade: shadeMatch[1], hex: shadeMatch[2]});
      }

      if (shades.length > 0) {
        cssColors.push(`\n    /* ${colorName.charAt(0).toUpperCase() + colorName.slice(1)} scale */`);
        shades.forEach(({shade, hex}) => {
          cssColors.push(`    --color-${colorName}-${shade}: ${hex};`);
        });
      }
    }

    return cssColors.join('\n');
  };

  const updateTailwindCSS = () => {
    try {
      const colorsContent = fs.readFileSync(colorsFilePath, 'utf-8');
      const cssColors = generateCSSColors(colorsContent);
      let tailwindContent = fs.readFileSync(tailwindFilePath, 'utf-8');
      const themeRegex = /@theme\s*\{([\s\S]*?)\}/;
      const themeMatch = tailwindContent.match(themeRegex);

      if (themeMatch) {
        const currentThemeContent = themeMatch[1];
        const breakpointMatch = currentThemeContent.match(
          /([\s\S]*?)(?=\n\s*\/\*.*scale|--color-black|--color-white|$)/,
        );
        const breakpoints = breakpointMatch ? breakpointMatch[1].trim() : '';
        const newThemeContent = `@theme {
    ${breakpoints}
${cssColors}
}`;
        tailwindContent = tailwindContent.replace(themeRegex, newThemeContent);
        fs.writeFileSync(tailwindFilePath, tailwindContent, 'utf-8');
        console.log('✓ Tailwind colors generated from Colors.ts');
      } else {
        console.warn('⚠ Could not find @theme block in tailwind.css');
      }
    } catch (error) {
      console.error('Error generating Tailwind colors:', error);
    }
  };

  return {
    name: 'generate-tailwind-colors',
    buildStart() {
      updateTailwindCSS();
    },
    configureServer(server) {
      const watcher = server.watcher;
      watcher.add(colorsFilePath);

      watcher.on('change', (file) => {
        if (file === colorsFilePath) {
          console.log('Colors.ts changed, regenerating Tailwind colors...');
          updateTailwindCSS();
          const tailwindModule = server.moduleGraph.getModuleById(tailwindFilePath);
          if (tailwindModule) {
            server.moduleGraph.invalidateModule(tailwindModule);
            server.ws.send({
              type: 'full-reload',
              path: '*',
            });
          }
        }
      });
    },
  };
}
