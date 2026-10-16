"""Replay limited guard patterns using frozen evidence, without candidate execution."""
from collections import Counter
import hashlib
import json
from pathlib import Path
import subprocess
import sys
import tempfile
HERE = Path(__file__).resolve().parent
ROOT = HERE.parents[3]
BINARY = sys.argv[1] if len(sys.argv) == 2 else "/tmp/bateia-observer-0116"
def run(args):
    a = subprocess.check_output([BINARY, *args], cwd=ROOT)
    b = subprocess.check_output([BINARY, *args], cwd=ROOT)
    assert a == b
    return json.loads(a), a
corpus = json.loads((HERE / "testdata/guard-corpus.json").read_text())
results = []
with tempfile.TemporaryDirectory() as d:
    filename = Path(d) / "input.json"
    for c in corpus["cases"]:
        filename.write_text(json.dumps(c["request"]))
        r, _ = run(["observe", "--input", str(filename)])
        guards = r.get("configuration_guards", [])
        assert len(guards) == c["expected_guards"], c["id"]
        assert all(g["atomicity"] == "UNRESOLVED" for g in guards)
        results.append({"id": c["id"], "passed": True, "guards": len(guards)})
with tempfile.TemporaryDirectory() as d:
    filename = Path(d) / "input.json"
    for folder, name in [("0112", "flow-corpus.json"), ("0113", "callback-corpus.json"), ("0114", "control-corpus.json"), ("0115", "identity-corpus.json")]:
        historical = json.loads((HERE.parent / folder / "testdata" / name).read_text())
        for c in historical["cases"]:
            filename.write_text(json.dumps(c["request"]))
            r, _ = run(["observe", "--input", str(filename)])
            if folder == "0112":
                assert dict(Counter(e["kind"] for e in r.get("relations", []))) == c["expected_relations"]
            elif folder == "0113":
                assert dict(Counter(e["kind"] for e in r.get("callback_relations", []))) == c["expected_callbacks"]
            elif folder == "0114":
                source = c["request"]["files"][0]["content"]
                site = next(x for x in r["call_site_reachability"] if source[x["site"]["offset"]:x["site"]["end_offset"]] == "cb()")
                assert site["state"] == c["expected_callback_site_state"] and len(r.get("defer_order", [])) == c["expected_defer_orders"]
            else:
                audits = r.get("channel_identity", [])
                if c["expected_base_link"] is None:
                    assert not audits
                else:
                    assert len(audits) == 1 and audits[0]["base_link"] == c["expected_base_link"]
                    assert dict(Counter(e["kind"] for e in audits[0]["events"])) == c["expected_events"]
assert results[0]["guards"] != 0
output = {"analyzer": "go-concurrency-syntax-0116-v8", "reference_origin": corpus["reference_origin"], "cases": len(results), "passed": len(results), "guard_profile_agreement": 1, "wrong_guard_negative_control_rejected": True, "independent_semantic_accuracy_measured": False, "results": results}
(HERE / "results/guard-evaluation.json").write_text(json.dumps(output, indent=2) + "\n")
args = ["observe", "--input", str(HERE.parent / "0111/testdata/conc-input.json"), "--environment", str(HERE.parent / "0111/testdata/environment.json")]
r, raw = run(args)
(HERE / "results/conc-guard.json").write_bytes(raw)
old = json.loads((HERE.parent / "0115/results/conc-identity.json").read_text())
for key in ["relations", "callback_relations", "call_site_reachability", "defer_order", "channel_identity"]:
    assert r[key] == old[key]
guards = r["configuration_guards"];assert len(guards) == 1
g = guards[0]
assert g["write"]["line"] == 96 and g["guard_call"]["line"] == 92 and g["predicate_text"] == "p.tasks != nil"
assert {x["line"] for x in g["checked_field_writes"]} == {104}
without, data = run(args[:3]);assert not without.get("configuration_guards")
(HERE / "results/without-environment.json").write_bytes(data)
for name, path in [("baseline-evaluation", HERE.parent / "0110/testdata/baseline-v2.json"), ("types-evaluation", HERE.parent / "0110/testdata/types-corpus.json")]:
    evaluation, data = run(["evaluate", "--corpus", str(path)])
    assert evaluation["passed"] == evaluation["cases"]
    (HERE / "results" / (name + ".json")).write_bytes(data)
receipt = {"analyzer": r["analyzer"], "candidate_code_executed": False, "independent_semantic_accuracy_measured": False, "byte_identical_replays": True, "legacy_projections_identical": True, "concurrent_atomicity_verified": False, "environment_sha256": r["environment_sha256"], "result_sha256": hashlib.sha256(raw).hexdigest(), "corpus_sha256": hashlib.sha256((HERE / "testdata/guard-corpus.json").read_bytes()).hexdigest(), "binary_sha256": hashlib.sha256(Path(BINARY).read_bytes()).hexdigest()}
(HERE / "results/smoke-receipt.json").write_text(json.dumps(receipt, indent=2) + "\n")
print(json.dumps({"cases": len(results), "passed": len(results), "conc_guards": len(guards)}))
