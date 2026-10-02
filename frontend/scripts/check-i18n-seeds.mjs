import {execFileSync} from 'node:child_process';

const locales = ['en', 'de'];
const seedDirectory = 'backend/conf/dbseeds/i18n';
const sourceRoots = ['frontend/apps/portal', 'frontend/libs/portal', 'frontend/libs/shared'];
const sourcePathPattern = /\.(?:[cm]?[jt]sx?|vue)$/;
const testPathPattern = /(?:^|\/)(?:tests?|__tests__)\/|\.(?:test|spec)\.[^.]+$/;
const translationCallPattern = /(?:\b(?:t|i18n\.t)|\$t)\(\s*(['"])([^'"\r\n]+)\1(?=\s*[),])/g;

function git(args, options = {}) {
  return execFileSync('git', args, {encoding: 'utf8', maxBuffer: 512 * 1024 * 1024, ...options});
}

function listSourceBlobs(args, oidIndex) {
  const blobs = new Map();
  for (const entry of git(args).split('\0')) {
    const tab = entry.indexOf('\t');
    if (tab < 0) continue;
    const path = entry.slice(tab + 1);
    const meta = entry.slice(0, tab).split(' ');
    if (meta[0] === '160000' || !sourcePathPattern.test(path) || testPathPattern.test(path)) continue;
    blobs.set(path, meta[oidIndex]);
  }
  return blobs;
}

function collectTranslationKeys(oids) {
  const keys = new Set();
  if (!oids.length) return keys;

  const output = git(['cat-file', '--batch'], {encoding: null, input: `${oids.join('\n')}\n`});
  let offset = 0;
  while (offset < output.length) {
    const headerEnd = output.indexOf(10, offset);
    const size = Number(output.toString('utf8', offset, headerEnd).split(' ')[2]);
    const start = headerEnd + 1;
    for (const match of output.toString('utf8', start, start + size).matchAll(translationCallPattern)) {
      keys.add(match[2]);
    }
    offset = start + size + 1;
  }
  return keys;
}

function readCatalog(path, trackedPaths) {
  if (!trackedPaths.has(path)) {
    throw new Error(`Missing translation catalog: ${path} (staged)`);
  }

  const catalog = JSON.parse(git(['show', `:${path}`]));
  if (!catalog || typeof catalog !== 'object' || Array.isArray(catalog)) {
    throw new Error(`Expected a flat translation dictionary: ${path}`);
  }
  for (const [key, value] of Object.entries(catalog)) {
    if (typeof value !== 'string') {
      throw new Error(`Expected a string translation: ${path}: ${key}`);
    }
  }
  return catalog;
}

function checkSeeds() {
  const seedPaths = new Set(git(['ls-files', '-z', '--', seedDirectory]).split('\0'));
  const seedCatalogs = locales.map((locale) => {
    const path = `${seedDirectory}/portal.${locale}.json`;
    return [path, readCatalog(path, seedPaths)];
  });

  const hasHead = git(['rev-parse', '--revs-only', 'HEAD']).trim();
  const stagedBlobs = listSourceBlobs(['ls-files', '-s', '-z', '--', ...sourceRoots], 1);
  const headBlobs = hasHead ? listSourceBlobs(['ls-tree', '-r', '-z', 'HEAD', '--', ...sourceRoots], 2) : new Map();
  const changedOids = [...stagedBlobs].filter(([path, oid]) => headBlobs.get(path) !== oid).map(([, oid]) => oid);
  const candidateKeys = collectTranslationKeys(changedOids);
  const headKeys = candidateKeys.size ? collectTranslationKeys([...new Set(headBlobs.values())]) : new Set();
  const addedKeys = [...candidateKeys].filter((key) => !headKeys.has(key)).sort();

  const missing = [];
  for (const [path, seedCatalog] of seedCatalogs) {
    for (const key of addedKeys) {
      if (!Object.hasOwn(seedCatalog, key)) missing.push(`${path}: ${key}`);
    }
  }

  if (missing.length) {
    console.error(`New translation keys are missing from staged seed files:\n${missing.join('\n')}`);
    console.error('Add the keys to both portal seed catalogs and stage the seed files before committing.');
    process.exitCode = 1;
    return;
  }
  console.log(`Translation seed check passed (${addedKeys.length} new keys).`);
}

try {
  checkSeeds();
} catch (error) {
  console.error(`Translation seed check failed: ${error.message}`);
  process.exitCode = 1;
}
