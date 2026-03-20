// load and transform all *.md files from the given directory
import { readdirSync, statSync } from "node:fs";

export const getChildren = (
  baseDir,
  directory,
  normalizedName = true,
  collapsed = true
) => {
  return readdirSync(`${baseDir}/${directory}`).map((entry) =>
    statSync(`${baseDir}/${directory}/${entry}`).isDirectory()
      ? {
        text: entry ? normalizeFolderNames(entry) : entry,
        collapsed,
        items: getChildren(baseDir, `${directory}/${entry}`),
      }
      : {
        text: normalizedName ? normalizeNames(entry) : entry,
        link: `/${directory}/${entry}`,
      }
  );
};
//normalize folder name folder-name to Folder Name
export const normalizeFolderNames = (fileName: string) => {
  const regex = /[a-z][a-z0-9]*(-[a-z0-9]+)*/g;
  if (!fileName.match(regex)) {
    return undefined;
  }
  const numberRegex = /[0-9]/g;

  // folder-name => Folder Name
  const normalizeFileName = fileName
    .split("-")
    //.filter((text) => !text.match(numberRegex))
    .map((text) => text.charAt(0).toUpperCase() + text.slice(1))
    .join(" ");
  return normalizeFileName;
};

// normalize file name file_name.md to File Name
export const normalizeNames = (fileName: string) => {
  const regex = /[A-Za-z][a-z0-9]*(_[a-z0-9]+)*(.)md/g;
  if (!fileName.match(regex)) {
    return undefined;
  }
  const numberRegex = /[0-9]/g;

  // a_b_c.md => [a_b_c]
  const fileNameWithoutSuffix = fileName.slice(0, -3);

  // [a_b_c] => A B C
  const normalizeFileName = fileNameWithoutSuffix
    .split("_")
    //.filter((text) => !text.match(numberRegex))
    .map((text) => text.charAt(0).toUpperCase() + text.slice(1))
    .join(" ");
  return normalizeFileName;

  //return fileNameWithoutSuffix;
};