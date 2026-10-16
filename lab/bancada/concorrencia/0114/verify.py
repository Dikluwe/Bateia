"""Verify limited structured-control projections using frozen inputs, without running candidates."""
from collections import Counter
import hashlib
import json
from pathlib import Path
import subprocess
import sys
import tempfile

HERE = Path(__file__).resolve().parent
ROOT = HERE.parents[3]
BINARY = sys.argv[1] if len(sys.argv) == 2 else "/tmp/bateia-observer-0114"

def run(args):
    a = subprocess.check_output([BINARY, *args], cwd=ROOT)
    b = subprocess.check_output([BINARY, *args], cwd=ROOT)
    assert a == b
    return json.loads(a), a

corpus = json.loads((HERE / "testdata/control-corpus.json").read_text())
results = []
with tempfile.TemporaryDirectory() as directory:
    filename = Path(directory) / "input.json"
    for c in corpus["cases"]:
        filename.write_text(json.dumps(c["request"]))
        r, _ = run(["observe", "--input", str(filename)])
        source = c["request"]["files"][0]["content"]
        site = next(s for s in r["call_site_reachability"] if source[s["site"]["offset"]:s["site"]["end_offset"]] == "cb()")
        assert site["state"] == c["expected_callback_site_state"]
        assert len(r.get("defer_order", [])) == c["expected_defer_orders"]
        results.append({"id": c["id"], "passed": True, "callback_site_state": site["state"], "defer_orders": len(r.get("defer_order", []))})
    assert results[0]["callback_site_state"] != "POSSIBLY_REACHABLE"
    for folder, key, actual_field in [("0112", "expected_relations", "relations"), ("0113", "expected_callbacks", "callback_relations")]:
        name = "flow-corpus.json" if folder == "0112" else "callback-corpus.json"
        old = json.loads((HERE.parent / folder / "testdata" / name).read_text())
        for c in old["cases"]:
            filename.write_text(json.dumps(c["request"]))
            r, _ = run(["observe", "--input", str(filename)])
            assert dict(Counter(e["kind"] for e in r.get(actual_field, []))) == c[key]
output = {"analyzer": "go-concurrency-syntax-0114-v6", "reference_origin": corpus["reference_origin"], "cases": len(results), "passed": len(results), "control_profile_agreement": 1, "independent_semantic_accuracy_measured": False, "wrong_reachability_negative_control_rejected": True, "legacy_flow_cases_passed": 13, "legacy_callback_cases_passed": 15, "results": results}
(HERE / "results/control-evaluation.json").write_text(json.dumps(output, indent=2) + "\n")
args = ["observe", "--input", str(HERE.parent / "0111/testdata/conc-input.json"), "--environment", str(HERE.parent / "0111/testdata/environment.json")]
r, raw = run(args)
(HERE / "results/conc-control.json").write_bytes(raw)
old = json.loads((HERE.parent / "0113/results/conc-callbacks.json").read_text())
assert r["relations"] == old["relations"] and r["callback_relations"] == old["callback_relations"]
worker = [e for e in r["defer_order"] if e["defer"]["path"] == "pool/pool.go" and e["defer"]["line"] == 153]
assert {e["site"]["line"] for e in worker} == {156, 160}
without, data = run(args[:3])
assert not without.get("call_site_reachability") and not without.get("defer_order")
(HERE / "results/without-environment.json").write_bytes(data)
for name, path in [("baseline-evaluation", HERE.parent / "0110/testdata/baseline-v2.json"), ("types-evaluation", HERE.parent / "0110/testdata/types-corpus.json")]:
    evaluation, data = run(["evaluate", "--corpus", str(path)])
    assert evaluation["passed"] == evaluation["cases"]
    (HERE / "results" / (name + ".json")).write_bytes(data)
receipt = {"analyzer": r["analyzer"], "candidate_code_executed": False, "independent_semantic_accuracy_measured": False, "byte_identical_replays": True, "legacy_relations_identical": True, "legacy_callbacks_identical": True, "site_states": dict(Counter(s["state"] for s in r["call_site_reachability"])), "defer_orders": len(r["defer_order"]), "environment_sha256": r["environment_sha256"], "result_sha256": hashlib.sha256(raw).hexdigest(), "corpus_sha256": hashlib.sha256((HERE / "testdata/control-corpus.json").read_bytes()).hexdigest(), "binary_sha256": hashlib.sha256(Path(BINARY).read_bytes()).hexdigest()}
(HERE / "results/smoke-receipt.json").write_text(json.dumps(receipt, indent=2) + "\n")
print(json.dumps({"cases": len(results), "passed": len(results), "conc_sites": len(r["call_site_reachability"]), "defer_orders": len(r["defer_order"])}))
