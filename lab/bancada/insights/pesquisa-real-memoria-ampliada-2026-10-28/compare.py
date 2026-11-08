"""Audit the two real reports without changing either report or memory."""
import hashlib
import json
from collections import Counter
from pathlib import Path

root = Path(__file__).resolve().parent
repo = root.parents[3]


def sha(content):
    return hashlib.sha256(content).hexdigest()


def load(name):
    return json.loads((root / name).read_text())


def documents(report):
    result = []
    groups = [("primary", b) for b in report["bundles"]]
    groups += [("primary", s["bundle"]) for s in report.get("site_contexts", [])]
    groups += [("followup", f["bundle"]) for f in report.get("followups", []) if "bundle" in f]
    for group, bundle in groups:
        for item in bundle["items"]:
            assert item["content_hash"] == "sha256:" + sha(item["content"].encode())
            ref = item["ref"]
            result.append({"route": group, "repository": ref["owner"] + "/" + ref["repo"],
                           "commit": ref["commit"], "path": item["path"],
                           "content_hash": item["content_hash"], "bytes": len(item["content"].encode())})
    return result


a, b = load("round-1.json"), load("round-2.json")
ra, rb = load("round-1-receipt.json"), load("round-2-receipt.json")
assert rb["memory_input_sha256"] == ra["result_sha256"] == sha((root / "round-1.json").read_bytes())
assert rb["result_sha256"] == sha((root / "round-2.json").read_bytes())
assert a["exploration_memory"] == b["plan"]["request"]["memory"]
assert b["exploration_memory"]["generation"] == a["exploration_memory"]["generation"] + 1
baseline = load("baseline-sha256.json")
changed = [p for p, h in baseline.items() if not (repo / p).is_file() or sha((repo / p).read_bytes()) != h]
assert not changed, changed
docs_a, docs_b = documents(a), documents(b)
known = {d["repository"] for d in docs_a}
entries_a = {e["id"]: e for e in a["exploration_memory"]["entries"]}
entries_b = {e["id"]: e for e in b["exploration_memory"]["entries"]}
memory_changes = []
for key, entry in entries_b.items():
    old = entries_a.get(key)
    if old != entry:
        memory_changes.append({"repository": entry["identity"]["external_id"],
                               "before": old["state"] if old else "ABSENT",
                               "after": entry["state"], "attempts_before": old["read_attempts"] if old else 0,
                               "attempts_after": entry["read_attempts"], "last_status": entry.get("last_status"),
                               "lead": entry.get("lead")})
summary = {"synthetic_data": False, "memory_input_verified_byte_for_byte": True,
           "memory_object_unchanged_between_rounds": True, "product_files_checked": len(baseline),
           "product_files_changed": changed, "round_1": {}, "round_2": {},
           "new_documents_round_2": [d for d in docs_b if d["repository"] not in known],
           "repeated_read_repositories": sorted(known & {d["repository"] for d in docs_b}),
           "memory_changes": memory_changes}
for number, report, receipt, docs in [(1, a, ra, docs_a), (2, b, rb, docs_b)]:
    states = Counter(e["state"] for e in report["exploration_memory"]["entries"])
    summary[f"round_{number}"] = {
        "status": report["status"], "elapsed_seconds": receipt["elapsed_seconds"],
        "catalog_calls": report["calls"], "candidates": len(report["candidates"]),
        "selected": [c["product"]["identity"]["external_id"] for c in report["candidates"] if c["selected"]],
        "documents": docs, "memory_states": dict(states), "memory_generation": report["exploration_memory"]["generation"],
        "readmes": report["readmes"],
        "followups": [{k: f[k] for k in ["entry_id", "lead_id", "identity", "status"]} for f in report.get("followups", [])],
        "catalog_statuses": dict(Counter(v["status"] for v in report["attempts"])),
        "novelty_state": report["novelty_state"], "effect_state": report["effect_state"],
    }
(root / "comparison.json").write_text(json.dumps(summary, ensure_ascii=False, indent=2) + "\n")
print(json.dumps({k: v for k, v in summary.items() if k not in ["round_1", "round_2", "memory_changes"]}, ensure_ascii=False, indent=2))
for number in [1, 2]:
    print("ROUND", number, json.dumps(summary[f"round_{number}"], ensure_ascii=False))
