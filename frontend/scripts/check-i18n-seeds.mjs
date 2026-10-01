import {execFileSync} from 'node:child_process';

const locales = ['en', 'de'];
const seedDirectory = 'backend/conf/dbseeds/i18n';
const sourcePathPattern = /^frontend\/(?:apps\/portal|libs\/(?:portal|shared))\/.*\.(?:[cm]?[jt]sx?|vue)$/;
const testPathPattern = /(?:^|\/)(?:tests?|__tests__)\/|\.(?:test|spec)\.[^.]+$/;
const translationCallPattern = /(?:\b(?:t|i18n\.t)|\$t)\(\s*(['"])([^'"\r\n]+)\1(?=\s*[),])/g;

function git(args) {
  return execFileSync('git', args, {encoding: 'utf8', maxBuffer: 10 * 1024 * 1024}).trim();
}

function readFile(revision, path) {
  return git(['show', `${revision}:${path}`]);
}

function readCatalog(revision, path, trackedPaths) {
  if (!trackedPaths.has(path)) {
    throw new Error(`Missing translation catalog: ${path} (${revision || 'staged'})`);
  }

  const catalog = JSON.parse(readFile(revision, path));
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

function collectTranslationKeys(revision, paths) {
  const keys = new Set();
  for (const path of paths) {
    if (!sourcePathPattern.test(path) || testPathPattern.test(path)) continue;
    for (const match of readFile(revision, path).matchAll(translationCallPattern)) {
      keys.add(match[2]);
    }
  }
  return keys;
}

function checkSeeds() {
  const stagedPaths = new Set(git(['ls-files', '-z']).split('\0'));
  const hasHead = git(['rev-parse', '--revs-only', 'HEAD']);
  const headPaths = new Set(hasHead ? git(['ls-tree', '-r', '--name-only', '-z', 'HEAD']).split('\0') : []);
  const stagedKeys = collectTranslationKeys('', stagedPaths);
  const headKeys = collectTranslationKeys('HEAD', headPaths);
  const addedKeys = new Set([...stagedKeys].filter((key) => !headKeys.has(key)));

  const missing = [];
  for (const locale of locales) {
    const path = `${seedDirectory}/portal.${locale}.json`;
    const seedCatalog = readCatalog('', path, stagedPaths);
    for (const key of [...addedKeys].sort()) {
      if (!Object.hasOwn(seedCatalog, key)) missing.push(`${path}: ${key}`);
    }
  }

  if (missing.length) {
    console.error(`New translation keys are missing from staged seed files:\n${missing.join('\n')}`);
    console.error('Add the keys to both portal seed catalogs and stage the seed files before committing.');
    process.exitCode = 1;
    return;
  }
  console.log(`Translation seed check passed (${addedKeys.size} new keys).`);
}

try {
  checkSeeds();
} catch (error) {
  console.error(`Translation seed check failed: ${error.message}`);
  process.exitCode = 1;
}
