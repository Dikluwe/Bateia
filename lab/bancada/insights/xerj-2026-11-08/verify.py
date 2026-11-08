"""Verify frozen bytes, card witnesses, proposal references and offline replay."""
import argparse
import hashlib
import json
from pathlib import Path
import subprocess
import tempfile


def sha(data):
    return hashlib.sha256(data).hexdigest()


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--binary", type=Path, required=True)
    args = parser.parse_args()
    here = Path(__file__).resolve().parent
    bundles = json.loads((here / "exploration-bundles.json").read_text())
    report = json.loads((here / "exploration.json").read_text())
    sources = json.loads((here / "sources.json").read_text())
    discovery = json.loads((here / "discovery.json").read_text())
    assert sha((here / "discovery.json").read_bytes()) == sources["discovery_sha256"]
    items = {}
    for bundle in bundles:
        ref = bundle["repository"]
        identity = ref["owner"] + "/" + ref["repo"]
        for item in bundle["items"]:
            digest = sha(item["content"].encode())
            assert item["content_hash"].removeprefix("sha256:") == digest
            assert item["ref"] == ref
            items[(identity, ref["commit"], item["path"])] = item
    for source in sources["sources"]:
        for record in source["files"]:
            item = items[("xerj-org/xerj", source["repository"]["commit"], record["path"])]
            assert sha(item["content"].encode()) == record["sha256"]
    assert bundles[2:] == discovery["bundles"]
    ids = set()
    for card in report["cards"]:
        ref = card["reference"]
        source = card["source"]
        item = items[(source["repository"], source["commit"], ref["path"])]
        assert ref["content_hash"].removeprefix("sha256:") == sha(item["content"].encode())
        excerpt = "\n".join(item["content"].split("\n")[ref["start_line"] - 1:ref["end_line"]])
        assert ref["excerpt"] == excerpt
        ids.add(card["id"])
    for task in report["tasks"]:
        assert set(task["card_refs"]) <= ids
    proposals = json.loads((here / "proposals.json").read_text())
    for proposal in proposals:
        assert set(proposal["card_refs"]) <= ids
        for ref in proposal["evidence"]:
            item = items[(ref["repository"], ref["commit"], ref["path"])]
            assert ref["sha256"] == sha(item["content"].encode())
            assert ref["excerpt"] == "\n".join(item["content"].split("\n")[ref["start_line"] - 1:ref["end_line"]])
    command = [str(args.binary.resolve()), "idea", "insights", "explore", "--question", report["question"],
               "--bundles", str(here / "exploration-bundles.json"), "--max-passages", "16"]
    for familiar in report["familiar_approaches"]:
        command += ["--familiar", familiar]
    with tempfile.TemporaryDirectory(prefix="bateia-xerj-replay-") as directory:
        output = Path(directory) / "replay.json"
        subprocess.run(command + ["--output", str(output)], cwd=directory, check=True)
        assert json.loads(output.read_text()) == report
    checks = {"bundle_hashes": True, "target_files": sum(len(s["files"]) for s in sources["sources"]),
              "discovery_preserved": True, "cards_checked": len(ids), "proposals_checked": len(proposals),
              "offline_replay_identical": True, "semantic_validity_or_effect_verified": False}
    (here / "replay-check.json").write_text(json.dumps(checks, indent=2) + "\n")
    print(json.dumps(checks))


if __name__ == "__main__":
    main()
