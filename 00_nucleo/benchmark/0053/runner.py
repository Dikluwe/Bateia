#!/usr/bin/env python3
"""Frozen, corpus-blind measurement runner. See runner-contract.md."""
import argparse
import hashlib
import json
import os
from pathlib import Path, PurePosixPath
import signal
import subprocess
import sys
import time


def sha256(data):
    return hashlib.sha256(data).hexdigest()


def encoded(value):
    return json.dumps(value, ensure_ascii=False, sort_keys=True, separators=(",", ":")).encode()


def write_json(path, value):
    Path(path).write_text(json.dumps(value, ensure_ascii=False, indent=2, sort_keys=True) + "\n", encoding="utf-8")


def sequence(obj, key, required=False):
    if required and key not in obj:
        raise ValueError("missing field: " + key)
    value = obj.get(key)
    if value is None:
        return []
    if not isinstance(value, list) or any(not isinstance(x, dict) for x in value):
        raise ValueError("expected object array: " + key)
    return value


def safe_path(path):
    if not isinstance(path, str) or not path or "\\" in path:
        raise ValueError("invalid relative path")
    p = PurePosixPath(path)
    if p.is_absolute() or ".." in p.parts or p.as_posix() != path or path == ".":
        raise ValueError("unsafe or noncanonical relative path: " + path)
    return path


def validate_digest(path, digest, snapshot):
    safe_path(path)
    if path not in snapshot or not isinstance(digest, str) or digest != sha256(snapshot[path]):
        raise ValueError("unvalidated evidence hash: " + path)


def validate_search(raw, snapshot):
    if not isinstance(raw, dict):
        raise ValueError("implementation report must be an object")
    if raw.get("error"):
        raise ValueError("implementation search error: " + str(raw["error"]))
    if not isinstance(raw.get("status"), str):
        raise ValueError("implementation report missing status")
    if raw.get("status") in ("error", "failed", "unknown", "not_evaluated"):
        raise ValueError("implementation search status: " + raw["status"])
    for item in sequence(raw, "examined", required=True):
        validate_digest(item.get("path"), item.get("sha256"), snapshot)
    for hit in sequence(raw, "hits", required=True):
        validate_digest(hit.get("path"), hit.get("sha256"), snapshot)
        if not isinstance(hit.get("symbol"), str):
            raise ValueError("hit missing symbol")


def discovery_observation(project, discovery):
    if discovery is None:
        return {"status": "not_evaluated", "reason": "discovery not requested", "expected_repo_top10": None}
    try:
        if not isinstance(discovery, dict):
            raise ValueError("discovery must be an object")
        candidates = sequence(discovery, "candidates", required=True)
        attempts = sequence(discovery, "attempts", required=True)
        error = discovery.get("error")
        attempt_errors = [str(a["error"]) for a in attempts if a.get("error")]
        if attempt_errors:
            error = "; ".join(([str(error)] if error else []) + attempt_errors)
        if error:
            return {"status": "unknown", "reason": str(error), "candidates": candidates, "attempts": attempts, "expected_repo_top10": None}
        expected = project["repo"].casefold()
        rank = next((i + 1 for i, c in enumerate(candidates[:10]) if (str(c.get("owner", "")) + "/" + str(c.get("repo", ""))).casefold() == expected), None)
        return {"status": "evaluated", "candidates": candidates, "attempts": attempts, "expected_repo_top10": rank is not None, "expected_repo_rank": rank, "commit_identity": "not_evaluated"}
    except (ValueError, KeyError, TypeError) as exc:
        return {"status": "unknown", "reason": str(exc), "expected_repo_top10": None}


def evaluate_case(case, project, raw, report, snapshot, *, errors=None, inventory_valid=True, discovery=None):
    """Pure measurement: snapshot is a mapping of relative paths to immutable bytes."""
    result = {"id": case["id"], "project_id": case["project_id"], "query": case["query"], "polarity": case["polarity"], "status": "unknown", "errors": list(errors or []), "inventory_valid": inventory_valid, "verified_behavior": "not_evaluated", "discovery": discovery_observation(project, discovery), "targets": [], "paths": []}
    if not inventory_valid:
        result["errors"].append("snapshot inventory/hash mismatch")
    if result["errors"]:
        return result
    try:
        validate_search(raw, snapshot)
        if not isinstance(report, dict):
            raise ValueError("CLI report must be an object")
        repos = sequence(report, "results", required=True)
        if report.get("failures"):
            raise ValueError("CLI reports execution failures: " + json.dumps(report["failures"], ensure_ascii=False))
        files, final_hits, coverage, requirements = {}, [], [], []
        for repo in repos:
            coverage.append(repo.get("coverage"))
            for f in sequence(repo, "files", required=True):
                path = f.get("path")
                if not isinstance(f.get("evidence"), dict):
                    raise ValueError("extracted file has no evidence")
                validate_digest(path, f["evidence"].get("sha256"), snapshot)
                if not isinstance(f.get("content"), str):
                    raise ValueError("extracted file has no content")
                if sha256(f["content"].encode("utf-8")) != sha256(snapshot[path]):
                    raise ValueError("extracted content bytes/hash differ from snapshot: " + path)
                files.setdefault(path, []).append(f)
            final = repo.get("implementation_search")
            if final is not None:
                validate_search(final, snapshot)
                final_hits.extend(sequence(final, "hits"))
                requirements.extend(sequence(final, "requirements"))
        raw_hits = sequence(raw, "hits")
        examined = {x["path"] for x in sequence(raw, "examined")}
        ranked_paths = list(dict.fromkeys(h["path"] for h in raw_hits))
        raw_keys = {(h["path"], h["symbol"]) for h in raw_hits}
        final_keys = {(h["path"], h["symbol"]) for h in final_hits}
        channels = {"lexical": [], "documentation": [], "direct_call": [], "other": []}
        for hit in raw_hits:
            channel = hit.get("retrieval") or "lexical"
            channels[channel if channel in channels else "other"].append(hit)
        for target in case["expected"]:
            path, symbol = target["path"], target["symbol"]
            validate_digest(path, next((f["sha256"] for f in project["files"] if f["path"] == path), None), snapshot)
            matches = [h for h in final_hits if (h["path"], h["symbol"]) == (path, symbol)]
            body = False
            for hit in matches:
                start, end = hit.get("line"), hit.get("end_line")
                if isinstance(start, int) and not isinstance(start, bool) and isinstance(end, int) and not isinstance(end, bool) and 1 <= start <= end <= len(snapshot[path].splitlines()):
                    body = body or any(f["content"].encode("utf-8") == snapshot[path] for f in files.get(path, []))
            result["targets"].append({"path": path, "symbol": symbol, "file_found": path in examined, "symbol_recovered": (path, symbol) in raw_keys, "path_top3": path in ranked_paths[:3], "path_top12": path in ranked_paths[:12], "file_extracted": path in files, "symbol_confirmed": (path, symbol) in final_keys, "body_available": body, "channels": sorted({(h.get("retrieval") or "lexical") for h in raw_hits if (h["path"], h["symbol"]) == (path, symbol)})})
        for path in dict.fromkeys(t["path"] for t in case["expected"]):
            result["paths"].append({"path": path, "file_found": path in examined, "path_top3": path in ranked_paths[:3], "path_top12": path in ranked_paths[:12], "file_extracted": path in files})
        all_requirements = sequence(raw, "requirements") + requirements
        claims = [r for r in all_requirements if str(r.get("status", "")).casefold() in ("satisfied", "verified", "proven", "verified_behavior", "behavior_verified")]
        result.update({"status": "evaluated", "raw_search": {k: raw.get(k) for k in ("status", "reason", "method", "limited", "issues", "budget", "files_read", "bytes_read")}, "raw_requirements": sequence(raw, "requirements"), "final_requirements": requirements, "coverage": coverage, "channels": channels, "raw_distinct_paths": ranked_paths, "raw_hit_count": len(raw_hits), "final_hit_count": len(final_hits), "negative_indications": bool(raw_hits or final_hits) if case["polarity"] == "negative" else None, "explicit_satisfaction_claims": claims, "undue_conclusions": bool(claims) if case["polarity"] == "negative" else None})
    except (ValueError, KeyError, TypeError, AttributeError, UnicodeError) as exc:
        result["errors"].append(str(exc))
        result["targets"], result["paths"] = [], []
    return result


PATH_METRICS = ("file_found", "path_top3", "path_top12", "file_extracted")
SYMBOL_METRICS = ("symbol_recovered", "symbol_confirmed", "body_available")


def fraction(values):
    values = list(values)
    numerator, denominator = sum(bool(x) for x in values), len(values)
    return {"numerator": numerator, "denominator": denominator, "ratio": numerator / denominator if denominator else None}


def summarize(cases):
    """Pure stable aggregate, deliberately excluding timing and receipts."""
    positives = [c for c in cases if c["polarity"] == "positive" and c["status"] == "evaluated"]
    negatives = [c for c in cases if c["polarity"] == "negative" and c["status"] == "evaluated"]
    disc = [c["discovery"] for c in cases if c["discovery"]["status"] == "evaluated"]
    metrics = {}
    for key in PATH_METRICS + SYMBOL_METRICS:
        field = "paths" if key in PATH_METRICS else "targets"
        metrics[key] = {"micro": fraction(t[key] for c in positives for t in c[field]), "positive_cases_any": fraction(any(t[key] for t in c[field]) for c in positives), "positive_cases_all": fraction(all(t[key] for t in c[field]) for c in positives)}
    return {"schema": "0053-results-v1", "total_cases": len(cases), "positive_evaluable": len(positives), "negative_evaluable": len(negatives), "unknown_cases": [c["id"] for c in cases if c["status"] == "unknown"], "metrics": metrics, "negative_indications": fraction(c["negative_indications"] for c in negatives), "negative_undue_conclusions": fraction(c["undue_conclusions"] for c in negatives), "discovery": {"expected_repo_top10": fraction(d["expected_repo_top10"] for d in disc), "unknown_cases": [c["id"] for c in cases if c["discovery"]["status"] == "unknown"], "not_evaluated_cases": [c["id"] for c in cases if c["discovery"]["status"] == "not_evaluated"]}, "verified_behavior": "not_evaluated", "integrated_remote_journey": "not_evaluated"}


def inventory(project):
    root = Path(project["local_path"])
    if not root.is_absolute() or not root.is_dir() or root.is_symlink():
        raise ValueError("snapshot must be an absolute nonsymlink directory")
    expected = {safe_path(f["path"]): f["sha256"] for f in project["files"]}
    if len(expected) != len(project["files"]):
        raise ValueError("duplicate inventory paths")
    snapshot, issues = {}, []
    def walk_error(exc):
        raise exc
    for current, dirs, names in os.walk(root, followlinks=False, onerror=walk_error):
        for name in list(dirs) + names:
            p = Path(current) / name
            if p.is_symlink():
                issues.append("symlink: " + p.relative_to(root).as_posix())
            elif p.is_file():
                snapshot[p.relative_to(root).as_posix()] = p.read_bytes()
            elif not p.is_dir():
                issues.append("nonregular entry: " + p.relative_to(root).as_posix())
    hashes = {p: sha256(b) for p, b in sorted(snapshot.items())}
    missing = sorted(set(expected) - set(hashes))
    extra = sorted(set(hashes) - set(expected))
    changed = sorted(p for p in set(expected) & set(hashes) if expected[p] != hashes[p])
    record = {"valid": not (missing or extra or changed or issues), "missing": missing, "extra": extra, "changed": changed, "issues": issues, "files": hashes, "inventory_sha256": sha256(encoded(hashes))}
    return snapshot, record


def checked_inventory(project):
    try:
        return inventory(project)
    except (OSError, ValueError, KeyError, TypeError) as exc:
        return {}, {"valid": False, "error": str(exc)}


def run_command(args, timeout, directory, name, json_file=None):
    """Keep command group bounded and preserve bytes even on error."""
    start = time.monotonic()
    receipt = {"args": args, "args_sha256": sha256(encoded(args)), "timeout_seconds": timeout}
    stdout, stderr, error = b"", b"", None
    try:
        process = subprocess.Popen(args, stdout=subprocess.PIPE, stderr=subprocess.PIPE, start_new_session=True)
        try:
            stdout, stderr = process.communicate(timeout=timeout)
        except subprocess.TimeoutExpired:
            os.killpg(process.pid, signal.SIGKILL)
            stdout, stderr = process.communicate()
            error = "timeout"
        receipt["returncode"] = process.returncode
        if process.returncode and error is None:
            error = "nonzero exit: " + str(process.returncode)
    except OSError as exc:
        error = str(exc)
    receipt["elapsed_seconds"] = round(time.monotonic() - start, 6)
    for label, data in (("stdout", stdout), ("stderr", stderr)):
        p = directory / (name + "." + label)
        p.write_bytes(data)
        receipt[label] = {"path": str(p), "sha256": sha256(data), "bytes": len(data)}
    payload = None
    data = stdout
    if json_file is not None:
        try:
            data = json_file.read_bytes()
            receipt["output"] = {"path": str(json_file), "sha256": sha256(data), "bytes": len(data)}
        except OSError as exc:
            error = error or str(exc)
    if error is None:
        try:
            payload = json.loads(data)
            if not isinstance(payload, dict):
                raise ValueError("JSON root must be an object")
        except (ValueError, UnicodeError) as exc:
            error = "invalid JSON: " + str(exc)
    receipt["error"] = error
    write_json(directory / (name + ".receipt.json"), receipt)
    return payload, receipt


def validate_corpus(corpus):
    projects = {p["id"]: p for p in corpus["projects"]}
    if len(projects) != len(corpus["projects"]):
        raise ValueError("duplicate project id")
    ids = set()
    for case in corpus["cases"]:
        if case["id"] in ids:
            raise ValueError("duplicate case id")
        ids.add(case["id"])
        project = projects[case["project_id"]]
        expected = case["expected"]
        if case["polarity"] not in ("positive", "negative") or not isinstance(case["query"], str) or not case["query"].strip():
            raise ValueError("invalid case polarity/query")
        if bool(expected) != (case["polarity"] == "positive"):
            raise ValueError("positive requires targets; negative requires empty expected")
        seen = set()
        for target in expected:
            safe_path(target["path"])
            key = (target["path"], target["symbol"])
            if key in seen or not isinstance(target["symbol"], str) or not target["symbol"]:
                raise ValueError("invalid or duplicate target")
            seen.add(key)
            if target["path"] not in {f["path"] for f in project["files"]}:
                raise ValueError("target outside inventory")
    return projects


def main(argv=None):
    parser = argparse.ArgumentParser(description=__doc__)
    for flag in ("corpus", "freeze", "binary", "probe", "output"):
        parser.add_argument("--" + flag, required=True)
    parser.add_argument("--discovery", action="store_true")
    opts = parser.parse_args(argv)
    output = Path(opts.output).absolute()
    if output.exists():
        parser.error("output directory must be new")
    try:
        corpus_bytes = Path(opts.corpus).read_bytes()
        freeze_bytes = Path(opts.freeze).read_bytes()
        freeze = json.loads(freeze_bytes)
        if freeze.get("corpus_sha256") != sha256(corpus_bytes):
            raise ValueError("corpus SHA256 does not match freeze.corpus_sha256")
        corpus = json.loads(corpus_bytes)
        projects = validate_corpus(corpus)
        binary, probe = str(Path(opts.binary).resolve()), str(Path(opts.probe).resolve())
        hashes = {"binary": sha256(Path(binary).read_bytes()), "probe": sha256(Path(probe).read_bytes()), "runner": sha256(Path(__file__).read_bytes()), "corpus": sha256(corpus_bytes), "freeze": sha256(freeze_bytes)}
        for project in projects.values():
            root = Path(project["local_path"]).resolve()
            if output == root or root in output.parents:
                raise ValueError("output must be outside project snapshots")
        output.mkdir(parents=True, exist_ok=False)
    except (OSError, ValueError, KeyError, TypeError) as exc:
        parser.error(str(exc))
    receipts = {"schema": "0053-receipts-v1", "inputs_sha256": hashes, "arguments": vars(opts), "arguments_sha256": sha256(encoded(vars(opts))), "cases": []}
    cases = []
    for index, case in enumerate(corpus["cases"]):
        directory = output / ("case-%02d" % (index + 1))
        directory.mkdir()
        project = projects[case["project_id"]]
        snapshot, before = checked_inventory(project)
        commands, errors, raw, report, discovery = [], [], None, None, None
        if before["valid"]:
            for label, args, timeout, json_file in (("raw", [probe, "local", project["local_path"], case["query"]], 60, None), ("cli", [binary, "analyze", "-path", project["local_path"], "-query", case["query"], "-format", "json", "-output", str(directory / "cli.json")], 120, directory / "cli.json")):
                payload, receipt = run_command(args, timeout, directory, label, json_file)
                commands.append(receipt)
                if label == "raw":
                    raw = payload
                else:
                    report = payload
                if receipt["error"]:
                    errors.append(label + ": " + receipt["error"])
            if opts.discovery and case.get("discovery"):
                discovery, receipt = run_command([probe, "discovery", case["query"]], 45, directory, "discovery")
                commands.append(receipt)
                if receipt["error"]:
                    discovery = {"candidates": [], "attempts": [], "error": receipt["error"]}
        elif opts.discovery and case.get("discovery"):
            discovery = {"candidates": [], "attempts": [], "error": "snapshot invalid; case not run"}
        _, after = checked_inventory(project)
        try:
            immutable = sha256(Path(binary).read_bytes()) == hashes["binary"] and sha256(Path(probe).read_bytes()) == hashes["probe"] and sha256(Path(opts.corpus).read_bytes()) == hashes["corpus"] and sha256(Path(opts.freeze).read_bytes()) == hashes["freeze"]
        except OSError:
            immutable = False
        if not immutable:
            errors.append("executable/corpus/freeze changed during run")
        measured = evaluate_case(case, project, raw, report, snapshot, errors=errors, inventory_valid=before["valid"] and after["valid"], discovery=discovery)
        cases.append(measured)
        write_json(directory / "evaluation.json", measured)
        receipts["cases"].append({"id": case["id"], "before": before, "after": after, "commands": commands})
        write_json(output / "cases.json", cases)
        write_json(output / "summary.json", summarize(cases))
        write_json(output / "receipts.json", receipts)
        print(case["id"] + ": " + measured["status"], flush=True)
    artifacts = {p.relative_to(output).as_posix(): sha256(p.read_bytes()) for p in sorted(output.rglob("*")) if p.is_file()}
    write_json(output / "artifacts.json", artifacts)
    return 0


if __name__ == "__main__":
    sys.exit(main())
