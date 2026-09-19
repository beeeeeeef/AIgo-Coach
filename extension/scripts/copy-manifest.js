import { copyFileSync } from 'fs';
import { resolve, dirname } from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

const src = resolve(__dirname, '../manifest.json');
const dest = resolve(__dirname, '../dist/manifest.json');

copyFileSync(src, dest);
console.log('✓ manifest.json copied to dist/');
