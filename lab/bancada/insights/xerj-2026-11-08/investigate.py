"""Freeze a bounded XERJ source selection, collect with Bateia, then explore offline.

Does not execute XERJ, install it, send feedback or publish anything. The caller
supplies two public checkouts and a discovery result. Source bytes must match
git objects at the expected pins. Completeness is only over the selected files.
"""
import argparse
import hashlib
import json
from pathlib import Path
import subprocess
import tempfile


PINS = {
    "main": "1557de2f5be981d9db441969201c92a8af6d2de3",
    "corpus-hub": "708e974516a4ba1ba9bfa312aff164041997cfd6",
}
FILES = {
    "main": [
        "LICENSE", "docs/ARCHITECTURE.md", "docs/XERJ_QUERY_1.md",
        "engine/crates/xerj-common/src/hub.rs",
        "engine/crates/xerj-hub/src/build_catalogue.rs",
        "engine/crates/xerj-hub/src/pack.rs",
        "engine/crates/xerj-hub/src/lockfile.rs",
        "engine/crates/xerj-hub/src/feedback.rs",
        "engine/crates/xerj-hub/src/router.rs",
    ],
    "corpus-hub": [
        "LICENSE", "tools/xerj-code/hub/README.md",
        "tools/xerj-code/hub/PROGRAM-100.md",
        "tools/xerj-code/hub/CONTRIBUTING.md",
        "tools/xerj-code/hub/TEMPLATE.json",
        "tools/xerj-code/hub/backlog/g7-wave1-2026-10-graded.json",
        "tools/xerj-code/hub/candidates/authors/mindemory--proxenos.json",
        "tools/xerj-code/hub/candidates/owner_overrides.py",
    ],
}
QUESTION = (
    "How can XERJ corpus evidence provenance verification distinguish retrieval "
    "relevance from tested behavior, expose applicability conditions and support "
    "reproducible counterexamples without confusing signatures or feedback with correctness?"
)
FAMILIAR = ["BM25", "vector search", "commit pinning", "signed packs"]


def digest(data):
    return hashlib.sha256(data).hexdigest()


def write_json(path, value):
    path.write_text(json.dumps(value, ensure_ascii=False, indent=2) + "\n")


def run(command, log, cwd):
    result = subprocess.run(command, cwd=cwd, capture_output=True, text=True)
    log.write_text(result.stdout + result.stderr)
    if result.returncode:
        raise RuntimeError(f"command failed ({result.returncode}); see {log}")


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--binary", type=Path, required=True)
    parser.add_argument("--main", type=Path, required=True)
    parser.add_argument("--hub", type=Path, required=True)
    parser.add_argument("--discovery", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()
    binary = str(args.binary.resolve())
    out = args.output.resolve()
    out.mkdir(parents=True, exist_ok=True)
    discovery_path = args.discovery.resolve()
    discovery = json.loads(discovery_path.read_text())
    plans, sources = [], []
    with tempfile.TemporaryDirectory(prefix="bateia-xerj-input-") as temporary:
        scratch = Path(temporary)
        for branch, checkout in [("main", args.main.resolve()), ("corpus-hub", args.hub.resolve())]:
            pin = PINS[branch]
            actual = subprocess.check_output(["git", "-C", str(checkout), "rev-parse", "HEAD"], text=True).strip()
            if actual != pin:
                raise ValueError(f"{branch}: expected {pin}, got {actual}")
            stage = scratch / branch
            stage.mkdir()
            records, targets = [], []
            for relative in FILES[branch]:
                data = (checkout / relative).read_bytes()
                original = subprocess.check_output(["git", "-C", str(checkout), "show", f"{pin}:{relative}"])
                if data != original:
                    raise ValueError(f"worktree differs from pinned object: {branch}/{relative}")
                dest = stage / relative
                dest.parent.mkdir(parents=True, exist_ok=True)
                dest.write_bytes(data)
                kind = "license" if relative == "LICENSE" else "documentation" if relative.endswith(".md") else "file"
                targets.append({"path": relative, "kind": kind, "reason": "explicit bounded source selection; no execution claim"})
                records.append({"path": relative, "sha256": digest(data), "bytes": len(data),
                                "url": f"https://github.com/xerj-org/xerj/blob/{pin}/{relative}"})
            ref = {"owner": "xerj-org", "repo": "xerj", "commit": pin}
            plans.append({"schema_version": "1", "repository": ref,
                          "candidate": {"owner": "xerj-org", "repo": "xerj", "commit_id": pin,
                                        "url": stage.as_uri()},
                          "decision": {"repository": f"xerj-org/xerj@{pin}"},
                          "targets": targets,
                          "budget": {"documents": 0, "manifests": 0, "code_files": 0,
                                     "requirement_files": 0, "tests_examples": 0,
                                     "total_bytes": 1048576, "tree_depth": 16},
                          "requirement_queries": []})
            sources.append({"branch": branch, "repository": ref, "scope": "EXPLICIT_FILE_SELECTION",
                            "source_bytes_verified_against_git_objects": True, "files": records})
        write_json(out / "collection-plan.json", plans)
        command = [binary, "idea", "evidence", "collect", "--plan", str(out / "collection-plan.json"),
                   "--output", str(out / "target-bundles.json")]
        run(command, out / "collect.log", scratch)
    target = json.loads((out / "target-bundles.json").read_text())
    for bundle, source in zip(target, sources):
        actual = {item["path"]: item for item in bundle["items"]}
        for record in source["files"]:
            item = actual.get(record["path"])
            if not item or digest(item["content"].encode()) != record["sha256"]:
                raise ValueError(f"collection did not preserve source: {record['path']}")
    # Discovery bundles are preserved verbatim, including incompleteness.
    combined = target + discovery["bundles"]
    write_json(out / "exploration-bundles.json", combined)
    write_json(out / "case.json", {"question": QUESTION, "familiar_approaches": FAMILIAR,
                                  "max_passages": 16, "target": "xerj-org/xerj",
                                  "interpretation_by": "AI in this session; not generated hypotheses from the binary"})
    write_json(out / "sources.json", {"scope": "Selected XERJ files plus README-only discovery",
                                     "binary_sha256": digest(Path(binary).read_bytes()),
                                     "discovery_sha256": digest(discovery_path.read_bytes()),
                                     "sources": sources,
                                     "limitations": ["No XERJ binary or source tests executed.",
                                                     "Complete bundles cover the staged selection, not the full repository.",
                                                     "Local provider metadata is consultative; git-object comparison is performed by this Lab harness.",
                                                     "Only the first 1024 lines per file are eligible for automatic passage selection."]})
    command = [binary, "idea", "insights", "explore", "--question", QUESTION,
               "--bundles", str(out / "exploration-bundles.json"), "--max-passages", "16",
               "--output", str(out / "exploration.json")]
    for familiar in FAMILIAR:
        command += ["--familiar", familiar]
    with tempfile.TemporaryDirectory(prefix="bateia-xerj-explore-") as cwd:
        run(command, out / "explore.log", cwd)
    print(f"Collected {sum(len(s['files']) for s in sources)} target files; exploration saved in {out}")


if __name__ == "__main__":
    main()
