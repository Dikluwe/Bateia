"""Replay a declaration-scoped identity audit, without candidate execution."""
from collections import Counter
import hashlib
import json
from pathlib import Path
import subprocess
import sys
import tempfile
HERE = Path(__file__).resolve().parent
ROOT = HERE.parents[3]
BINARY = sys.argv[1] if len(sys.argv) == 2 else "/tmp/bateia-observer-0115"
def run(args):
    a = subprocess.check_output([BINARY, *args], cwd=ROOT)
    b = subprocess.check_output([BINARY, *args], cwd=ROOT)
    assert a == b
    return json.loads(a), a
corpus = json.loads((HERE / "testdata/identity-corpus.json").read_text())
results = []
with tempfile.TemporaryDirectory() as d:
    filename = Path(d) / "input.json"
    for c in corpus["cases"]:
        filename.write_text(json.dumps(c["request"]))
        r, _ = run(["observe", "--input", str(filename)])
        audits = r.get("channel_identity", [])
        if c["expected_base_link"] is None:
            assert not audits
        else:
            assert len(audits) == 1
            a = audits[0]
            assert a["base_link"] == c["expected_base_link"]
            assert a["state"] == "INSTANCE_IDENTITY_UNRESOLVED"
            assert dict(Counter(e["kind"] for e in a["events"])) == c["expected_events"], (c["id"], a["events"])
        results.append({"id": c["id"], "passed": True})
    for folder, name, key, actual_field in [("0112", "flow-corpus.json", "expected_relations", "relations"), ("0113", "callback-corpus.json", "expected_callbacks", "callback_relations")]:
        old = json.loads((HERE.parent / folder / "testdata" / name).read_text())
        for c in old["cases"]:
            filename.write_text(json.dumps(c["request"]))
            r, _ = run(["observe", "--input", str(filename)])
            assert dict(Counter(e["kind"] for e in r.get(actual_field, []))) == c[key]
    old = json.loads((HERE.parent / "0114/testdata/control-corpus.json").read_text())
    for c in old["cases"]:
        filename.write_text(json.dumps(c["request"]))
        r, _ = run(["observe", "--input", str(filename)])
        source = c["request"]["files"][0]["content"]
        site = next(s for s in r["call_site_reachability"] if source[s["site"]["offset"]:s["site"]["end_offset"]] == "cb()")
        assert site["state"] == c["expected_callback_site_state"] and len(r.get("defer_order", [])) == c["expected_defer_orders"]
output = {"analyzer": "go-concurrency-syntax-0115-v7", "reference_origin": corpus["reference_origin"], "cases": len(results), "passed": len(results), "identity_audit_profile_agreement": 1, "independent_semantic_accuracy_measured": False, "absence_of_events_certifies_stability": False, "results": results}
(HERE / "results/identity-evaluation.json").write_text(json.dumps(output, indent=2) + "\n")
args = ["observe", "--input", str(HERE.parent / "0111/testdata/conc-input.json"), "--environment", str(HERE.parent / "0111/testdata/environment.json")]
r, raw = run(args)
(HERE / "results/conc-identity.json").write_bytes(raw)
old = json.loads((HERE.parent / "0114/results/conc-control.json").read_text())
for key in ["relations", "callback_relations", "call_site_reachability", "defer_order"]:
    assert r[key] == old[key]
audit = r["channel_identity"][0]
assert audit["base_link"] == "DISTINCT_BASE_OBJECTS" and audit["state"] == "INSTANCE_IDENTITY_UNRESOLVED"
assert dict(Counter(e["kind"] for e in audit["events"])) == {"FIELD_WRITE": 1, "RECEIVER_CAPTURE": 2}
without, data = run(args[:3]);assert not without.get("channel_identity")
(HERE / "results/without-environment.json").write_bytes(data)
for name, path in [("baseline-evaluation", HERE.parent / "0110/testdata/baseline-v2.json"), ("types-evaluation", HERE.parent / "0110/testdata/types-corpus.json")]:
    evaluation, data = run(["evaluate", "--corpus", str(path)])
    assert evaluation["passed"] == evaluation["cases"]
    (HERE / "results" / (name + ".json")).write_bytes(data)
receipt = {"analyzer": r["analyzer"], "candidate_code_executed": False, "independent_semantic_accuracy_measured": False, "byte_identical_replays": True, "legacy_projections_identical": True, "instance_identity_verified": False, "environment_sha256": r["environment_sha256"], "result_sha256": hashlib.sha256(raw).hexdigest(), "corpus_sha256": hashlib.sha256((HERE / "testdata/identity-corpus.json").read_bytes()).hexdigest(), "binary_sha256": hashlib.sha256(Path(BINARY).read_bytes()).hexdigest()}
(HERE / "results/smoke-receipt.json").write_text(json.dumps(receipt, indent=2) + "\n")
print(json.dumps({"cases": len(results), "passed": len(results), "conc_identity_audits": len(r["channel_identity"])}))
