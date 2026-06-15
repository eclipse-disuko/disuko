#!/usr/bin/env python3
# SPDX-FileCopyrightText: 2026 Mercedes-Benz Group AG and Mercedes-Benz AG
#
# SPDX-License-Identifier: Apache-2.0

import argparse
import json
import sys
import uuid
from datetime import datetime, timezone
from pathlib import Path
from urllib.error import HTTPError, URLError
from urllib.request import Request, urlopen


DEFAULT_SOURCE_URL = "https://api.github.com/repos/github/choosealicense.com/contents/_licenses?ref=gh-pages"
DEFAULT_LICENSES_FILE = "backend/conf/dbseeds/disuko/licenses.jsonl"
DEFAULT_OBLIGATIONS_FILE = "backend/conf/dbseeds/disuko/obligations.jsonl"
CHOOSEALICENSE_URL = "https://choosealicense.com/licenses/"

TAG_MAPPING = {
    "permissions.commercial-use": ("Commercial use", "right"),
    "permissions.distribution": ("Distribution", "right"),
    "permissions.modifications": ("Modification", "right"),
    "permissions.patent-use": ("Patent use", "right"),
    "permissions.private-use": ("Private use", "right"),
    "conditions.disclose-source": ("Provide source code location", "obligation"),
    "conditions.document-changes": ("State changes", "obligation"),
    "conditions.include-copyright": ("Display copyright notice", "obligation"),
    "conditions.include-copyright--source": ("License and copyright notice for source", "obligation"),
    "conditions.network-use-disclose": ("Network use is distribution", "obligation"),
    "conditions.same-license": ("Same license", "obligation"),
    "conditions.same-license--file": ("Same license (file)", "obligation"),
    "conditions.same-license--library": ("Same license (library)", "obligation"),
    "limitations.liability": ("Liability", "limitation"),
    "limitations.patent-use": ("Patent use", "limitation"),
    "limitations.trademark-use": ("Trademark use", "limitation"),
    "limitations.warranty": ("Warranty", "limitation"),
}


def fetch_text(url):
    request = Request(url, headers={"User-Agent": "disuko-choosealicense-seed-script"})
    try:
        with urlopen(request, timeout=60) as response:
            return response.read().decode("utf-8")
    except HTTPError as exc:
        raise RuntimeError(f"HTTP {exc.code} while fetching {url}") from exc
    except URLError as exc:
        raise RuntimeError(f"Cannot fetch {url}: {exc.reason}") from exc


def load_json_objects(path):
    text = path.read_text(encoding="utf-8") if path.exists() else ""
    decoder = json.JSONDecoder()
    objects = []
    pos = 0
    while pos < len(text):
        while pos < len(text) and text[pos].isspace():
            pos += 1
        if pos >= len(text):
            break
        obj, end = decoder.raw_decode(text, pos)
        objects.append(obj)
        pos = end
    return objects


def append_json_objects(path, objects):
    if not objects:
        return
    path.parent.mkdir(parents=True, exist_ok=True)
    needs_separator = path.exists() and path.read_text(encoding="utf-8").strip()
    with path.open("a", encoding="utf-8", newline="\n") as file:
        if needs_separator:
            file.write("\n\n")
        file.write("\n\n".join(json.dumps(obj, ensure_ascii=False, indent=2) for obj in objects))
        file.write("\n")


def load_obligation_keys(path):
    keys = {}
    for obj in load_json_objects(path):
        name = obj.get("name")
        obligation_type = obj.get("type")
        key = obj.get("_key")
        if name and obligation_type and key:
            keys[(name, obligation_type)] = key
    return keys


def parse_frontmatter(text):
    lines = text.splitlines()
    if not lines or lines[0].strip() != "---":
        raise RuntimeError("missing frontmatter")
    end = None
    for index in range(1, len(lines)):
        if lines[index].strip() == "---":
            end = index
            break
    if end is None:
        raise RuntimeError("missing frontmatter end")

    data = {}
    current_list = None
    for line in lines[1:end]:
        stripped = line.strip()
        if not stripped:
            continue
        if stripped.startswith("- ") and current_list:
            data.setdefault(current_list, []).append(stripped[2:].strip().strip("\"'"))
            continue
        current_list = None
        if ":" not in stripped:
            continue
        key, value = stripped.split(":", 1)
        key = key.strip()
        value = value.strip().strip("\"'")
        if key in {"permissions", "conditions", "limitations"}:
            current_list = key
            data.setdefault(key, [])
        else:
            data[key] = value

    body = "\n".join(lines[end + 1 :]).strip() + "\n"
    return data, body


def family_from_conditions(conditions):
    if "network-use-disclose" in conditions:
        return "network copyleft"
    if "same-license" in conditions:
        return "strong copyleft"
    if any(tag in conditions for tag in ["same-license--file", "same-license--library"]):
        return "weak copyleft"
    return "permissive"


def map_obligations(metadata, obligation_keys):
    keys = []
    unknown = []
    for group in ["permissions", "conditions", "limitations"]:
        for tag in metadata.get(group, []):
            mapped = TAG_MAPPING.get(f"{group}.{tag}")
            if not mapped:
                unknown.append(f"{group}.{tag}")
                continue
            obligation_key = obligation_keys.get(mapped)
            if not obligation_key:
                unknown.append(f"{group}.{tag} -> {mapped[0]} / {mapped[1]}")
                continue
            if obligation_key not in keys:
                keys.append(obligation_key)
    return keys, unknown


def load_choosealicense_files(source_url):
    files = json.loads(fetch_text(source_url))
    return [item for item in files if item.get("download_url") and item.get("name", "").endswith(".txt")]


def build_license_doc(file_info, obligation_keys, now):
    raw_text = fetch_text(file_info["download_url"])
    metadata, license_text = parse_frontmatter(raw_text)
    license_id = metadata.get("spdx-id", "").strip()
    title = metadata.get("title", "").strip()
    if not license_id:
        raise RuntimeError(f"{file_info['name']} has no spdx-id")
    if not title:
        raise RuntimeError(f"{file_info['name']} has no title")

    slug = file_info["name"].removesuffix(".txt")
    obligation_key_list, unknown_tags = map_obligations(metadata, obligation_keys)
    key = str(uuid.uuid4())

    return {
        "_key": key,
        "_id": f"licenses/{key}",
        "Created": now,
        "Deleted": False,
        "Updated": now,
        "description": metadata.get("description", ""),
        "isDeprecatedLicenseId": False,
        "licenseId": license_id,
        "meta": {
            "IsLicenseChart": True,
            "approvalState": "approved",
            "family": family_from_conditions(metadata.get("conditions", [])),
            "fsfApproved": False,
            "licenseType": "open source",
            "licenseUrl": f"{CHOOSEALICENSE_URL}{slug}",
            "obligationsKeyList": obligation_key_list,
            "osiApproved": False,
            "reviewDate": now[:10],
            "reviewState": "reviewed",
            "sourceUrl": file_info["download_url"],
        },
        "name": title,
        "source": "choosealicense",
        "text": license_text,
    }, unknown_tags


def main():
    parser = argparse.ArgumentParser(description="Fetch new ChooseALicense licenses and append them to the license seed file.")
    parser.add_argument("--licenses-file", default=DEFAULT_LICENSES_FILE)
    parser.add_argument("--obligations-file", default=DEFAULT_OBLIGATIONS_FILE)
    parser.add_argument("--source-url", default=DEFAULT_SOURCE_URL)
    parser.add_argument("--dry-run", action="store_true")
    args = parser.parse_args()

    licenses_path = Path(args.licenses_file)
    obligations_path = Path(args.obligations_file)
    existing_ids = {obj.get("licenseId") for obj in load_json_objects(licenses_path) if obj.get("licenseId")}
    obligation_keys = load_obligation_keys(obligations_path)
    now = datetime.now(timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")

    new_docs = []
    skipped = []
    errors = []
    unknown_tags = set()

    for file_info in load_choosealicense_files(args.source_url):
        try:
            doc, doc_unknown_tags = build_license_doc(file_info, obligation_keys, now)
        except Exception as exc:
            errors.append(f"{file_info.get('name', '<unknown>')}: {exc}")
            continue

        if doc["licenseId"] in existing_ids:
            skipped.append(doc["licenseId"])
            continue

        new_docs.append(doc)
        existing_ids.add(doc["licenseId"])
        unknown_tags.update(doc_unknown_tags)

    if unknown_tags:
        print("Unknown ChooseALicense tags or missing obligation mappings:", file=sys.stderr)
        for tag in sorted(unknown_tags):
            print(f"  - {tag}", file=sys.stderr)
        return 2

    if errors:
        print("Errors while building licenses:", file=sys.stderr)
        for error in errors:
            print(f"  - {error}", file=sys.stderr)
        return 2

    if not args.dry_run:
        append_json_objects(licenses_path, new_docs)

    print(f"ChooseALicense files loaded: {len(new_docs) + len(skipped)}")
    print(f"Existing licenses skipped: {len(skipped)}")
    print(f"New licenses found: {len(new_docs)}")
    if new_docs:
        print("New license IDs:")
        for doc in new_docs:
            print(f"  - {doc['licenseId']} | {doc['name']}")
    if args.dry_run:
        print("Dry run only. No file was changed.")
    else:
        print(f"Updated seed file: {licenses_path}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
