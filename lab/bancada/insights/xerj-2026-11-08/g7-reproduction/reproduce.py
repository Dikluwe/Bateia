"""Exercise XERJ's original publication gates on a synthetic corpus, offline.

Runs the official catalogue builder, signature verifier, check-catalogue,
registry validator and full site generator. Never deploys, sends feedback,
starts a server or runs the code represented by a corpus snapshot.
"""
import argparse
import hashlib
import importlib.util
import io
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import tarfile
import tempfile


SLUG = "mini-g7"
PIN_A = "a" * 40
PIN_B = "b" * 40


def sha(data):
    return hashlib.sha256(data).hexdigest()


def dump(path, data):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(data, indent=2) + "\n")


def module(path, name):
    spec = importlib.util.spec_from_file_location(name, path)
    value = importlib.util.module_from_spec(spec)
    sys.modules[name] = value
    spec.loader.exec_module(value)
    return value


def main():
    ap = argparse.ArgumentParser(description=__doc__)
    ap.add_argument("--hub", type=Path, required=True)
    ap.add_argument("--main", type=Path, required=True)
    ap.add_argument("--binary", type=Path, required=True)
    ap.add_argument("--output", type=Path, required=True)
    a = ap.parse_args()
    hub, binary, output = a.hub.resolve(), a.binary.resolve(), a.output.resolve()
    main_source = a.main.resolve()
    for checkout, expected in [(hub, "708e974516a4ba1ba9bfa312aff164041997cfd6"),
                               (main_source, "1557de2f5be981d9db441969201c92a8af6d2de3")]:
        actual = subprocess.check_output(["git", "-C", str(checkout), "rev-parse", "HEAD"], text=True).strip()
        if actual != expected:
            raise ValueError(f"source pin mismatch: {checkout}: {actual}")
    if sha(binary.read_bytes()) != "8fa845b4b347ac1fa3d9674169f5082e2ca465ac2085362454d2dbdd9d98502c":
        raise ValueError("binary differs from the official rc.101 artifact used in this reproduction")
    output.mkdir(parents=True, exist_ok=True)
    records = []
    env = os.environ.copy()
    env.update({"XERJ_DISABLE_FEEDBACK": "true", "XERJ_HUB_OFFLINE": "1", "HUB_STAMP": "2026-11-08", "PYTHONDONTWRITEBYTECODE": "1"})

    def run(name, command, cwd, expect=0):
        r = subprocess.run([str(x) for x in command], cwd=cwd, env=env, capture_output=True, text=True)
        (output / (name + ".log")).write_text(r.stdout + r.stderr)
        records.append({"name": name, "command": [str(x) for x in command], "exit_code": r.returncode})
        if expect is not None and r.returncode != expect:
            raise RuntimeError(f"{name}: exit {r.returncode}, expected {expect}; see log")
        return r

    with tempfile.TemporaryDirectory(prefix="xerj-g7-pipeline-") as directory:
        root = Path(directory)
        site = root / "tools/corpus-hub-site"
        shutil.copytree(hub / "tools/corpus-hub-site", site, ignore=shutil.ignore_patterns("out", "__pycache__"))
        registry = root / "tools/xerj-code/hub"
        registry.mkdir(parents=True)
        (root / "tools/packs").mkdir()
        recipe_dir = root / "tools/packs" / SLUG
        recipe_dir.mkdir()
        (recipe_dir / "recipe.toml").write_text(
            '[recipe]\nname = "mini-g7"\nformat = 1\n'
            '[[sources]]\nslug = "fixture"\nkind = "dir"\npath = "records"\nlicence = "MIT"\n')
        (registry / "WAVES.md").write_text("# Synthetic fixture; no wave score\n")
        for relative in ["tools/xerj-code/hub/validate_hub.py", "tools/hub-packs/hub_packs.py", "tools/hub-packs/ed25519.sh"]:
            dest = root / relative
            dest.parent.mkdir(parents=True, exist_ok=True)
            shutil.copyfile(hub / relative, dest)
        hp = module(root / "tools/hub-packs/hub_packs.py", "xerj_pipeline_fixture")
        current_gate = root / "main-pipeline/hub_packs.py"
        current_gate.parent.mkdir()
        shutil.copyfile(main_source / "tools/hub-packs/hub_packs.py", current_gate)
        repo_url = "https://github.com/example/mini-g7"
        manifest_path = registry / (SLUG + ".json")
        dump(manifest_path, {"corpus": SLUG, "cloned_at": "2026-10-01T00:00:00Z", "repos": [{
            "repo": SLUG, "url": repo_url, "sha": PIN_A, "files": 1, "bytes": 24, "licence": "MIT",
            "review": {"spdx": "MIT", "use": "adapt-with-attribution", "by": "fixture", "at": "2026-10-01"}}]})
        backlog = registry / "backlog/backlog-100.json"
        dump(backlog, {"statuses": ["live"], "cats": ["NET"], "lanes": ["A"], "refreshes": ["weekly"],
                       "corpora": [{"slug": SLUG, "cat": "NET", "lane": "A", "refresh": "weekly", "status": "live",
                                    "domain": "an agent working on synthetic fixtures queries this for a demo function",
                                    "sources": [{"url": repo_url}]}]})
        grade_path = registry / "backlog/g7-fixture-graded.json"
        dump(grade_path, {"graded_at": "2026-10-01", "corpora": {SLUG: {
            "median": 5.0, "queries": [{"q": f"fixture query {n}", "verdict": "pass"} for n in range(5)]}}})
        grade_hash = sha(grade_path.read_bytes())
        shutil.copyfile(grade_path, output / "old-g7-suite.json")
        seed = root / "throwaway.key"
        seed.write_text(os.urandom(32).hex())
        sign = root / "tools/hub-packs/ed25519.sh"
        public = run("derive-key", ["bash", sign, "derive", seed], root).stdout.strip()
        env.update({"XERJ_HUB_PUBKEY": public, "XERJ_HUB_HOME": str(root / "client"), "XERJ_CODE_HOME": str(root / "code")})
        packs = root / "packs"
        packs.mkdir()

        def pack(pin):
            pm = {"format_version": 1, "slug": SLUG, "pin_sha": pin, "pin_date": "2026-10-01",
                  "built_at": "2026-11-08T00:00:00Z", "builder_version": "synthetic fixture",
                  "index_prefix": f"xc-{SLUG}-b1", "indices": [{"name": f"xc-{SLUG}-b1-docs", "docs": 1, "bytes": 24}],
                  "snapshot_name": "fixture", "licence_spdx": "MIT", "licence_class": "adapt-with-attribution",
                  "sources": [{"repo": SLUG, "url": repo_url, "sha": pin, "licence": "MIT"}], "defs_count": 1, "files_count": 1}
            members = {"pack.json": json.dumps(pm).encode(),
                       "defs.jsonl": (json.dumps({"symbol": "demo", "kind": "function", "lang": "c", "file": "demo.c",
                                                  "line": 1, "text": "int demo(void)", "value": None}) + "\n").encode(),
                       "files.jsonl": (json.dumps({"path": "demo.c", "lang": "c", "bytes": 24, "licence": "MIT"}) + "\n").encode(),
                       "snapshot/opaque": b"not an executable or a restored index: " + pin.encode()}
            members["SHA256SUMS"] = "".join(f"{sha(data)}  {name}\n" for name, data in sorted(members.items())).encode()
            stream = io.BytesIO()
            with tarfile.open(fileobj=stream, mode="w") as archive:
                for name, data in sorted(members.items()):
                    info = tarfile.TarInfo(name)
                    info.size, info.mode, info.mtime = len(data), 0o644, 0
                    archive.addfile(info, io.BytesIO(data))
            dest = packs / f"{SLUG}-{pin[:8]}.hubpack.tar.zst"
            compressed = subprocess.run(["zstd", "-q", "-c"], input=stream.getvalue(), capture_output=True, check=True).stdout
            dest.write_bytes(compressed)
            # This older branch's helper verifies but does not sign. Sign the
            # synthetic pack with OpenSSL using the same RFC8410 key encoding.
            der = root / "throwaway.der"
            der.write_bytes(bytes.fromhex("302e020100300506032b657004220420" + seed.read_text().strip()))
            der.chmod(0o600)
            signed = subprocess.run(["openssl", "pkeyutl", "-sign", "-inkey", str(der), "-keyform", "DER",
                                     "-rawin", "-in", str(dest)], capture_output=True, check=True).stdout.hex()
            Path(str(dest) + ".sig").write_text(signed)
            return dest

        pack(PIN_A)
        pack_b = pack(PIN_B)

        def stage(name):
            catalogue = root / (name + ".json")
            run(name + "-validate", [sys.executable, registry / "validate_hub.py"], root)
            run(name + "-catalogue", [binary, "hub", "build-catalogue", "--registry", registry,
                                      "--packs", packs, "--out", catalogue, "--sign-key", seed], root)
            run(name + "-signature", ["bash", sign, "verify", public, catalogue, str(catalogue) + ".sig"], root)
            # corpus-hub still expects format 1, while rc.101 emits format 2.
            # Record that genuine blocker; never patch the original validator.
            old_gate = run(name + "-old-branch-gate", [sys.executable, root / "tools/hub-packs/hub_packs.py", "check-catalogue",
                                             "--registry", registry, "--backlog", backlog, "--catalogue", catalogue,
                                             "--repo", "xerj-org/xerj"], root, expect=1)
            main_gate = run(name + "-main-publication-gate", [sys.executable, current_gate, "check-catalogue",
                                             "--registry", registry, "--backlog", backlog, "--catalogue", catalogue,
                                             "--repo", "xerj-org/xerj"], root, expect=1)
            assert "format_version is 2, expected 1" in old_gate.stderr
            assert "format_version is 2, expected 1" in main_gate.stderr
            run(name + "-site", [sys.executable, site / "gen.py"], root)
            entry = json.loads(catalogue.read_text())["corpora"][0]
            html = (site / "out/corpus/mini-g7.html").read_text()
            shutil.copyfile(catalogue, output / (name + "-catalogue.json"))
            shutil.copyfile(str(catalogue) + ".sig", output / (name + "-catalogue.json.sig"))
            shutil.copyfile(manifest_path, output / (name + "-manifest.json"))
            (output / (name + "-page.html")).write_text(html)
            shutil.copyfile(site / "out/index.html", output / (name + "-index.html"))
            result = {"pin_sha": entry["pin_sha"], "quality_grade": entry.get("quality_grade"),
                      "page_has_grade_5": "5.0/5 relevant" in html, "page_has_new_pin": PIN_B[:12] in html,
                      "g7_suite_sha256": sha(grade_path.read_bytes()) if grade_path.exists() else None}
            return result

        results = {"before": stage("before")}
        old = manifest_path.read_text()
        moved = hp.rewrite_manifest(old, [hp.Move(SLUG, repo_url, PIN_A, PIN_B)], "2026-11-08")
        manifest_path.write_text(moved)
        results["after_repin"] = stage("after-repin")
        assert results["after_repin"]["pin_sha"] == PIN_B
        assert results["after_repin"]["quality_grade"] == "G7 5.0/5"
        assert results["after_repin"]["page_has_grade_5"]
        assert results["after_repin"]["g7_suite_sha256"] == grade_hash
        grade_path.unlink()
        results["without_suite"] = stage("without-suite")
        assert results["without_suite"]["quality_grade"] is None
        assert not results["without_suite"]["page_has_grade_5"]
        # The original gate catches a mismatched pin, proving it is actually exercised.
        wrong = json.loads((output / "after-repin-catalogue.json").read_text())
        # Isolate the grade/pin guard from the independent schema mismatch.
        # This is a synthetic format-1 input, NOT an unmodified builder output.
        wrong["format_version"] = 1
        for entry in wrong["corpora"]:
            entry.pop("upstream", None)
        dump(root / "format1-control.json", wrong)
        run("format1-gate-control", [sys.executable, current_gate, "check-catalogue",
                                    "--registry", registry, "--backlog", backlog, "--catalogue", root / "format1-control.json",
                                    "--repo", "xerj-org/xerj"], root)
        wrong["corpora"][0]["pin_sha"] = PIN_A
        dump(root / "wrong-pin.json", wrong)
        mismatch = run("negative-wrong-pin", [sys.executable, current_gate, "check-catalogue",
                                    "--registry", registry, "--backlog", backlog, "--catalogue", root / "wrong-pin.json",
                                    "--repo", "xerj-org/xerj"], root, expect=1)
        assert "pin_sha" in mismatch.stderr and "manifest" in mismatch.stderr
        original = pack_b.read_bytes()
        pack_b.write_bytes(original + b"tampered")
        tampered = run("negative-signature", [binary, "hub", "build-catalogue", "--registry", registry,
                                    "--packs", packs, "--out", root / "bad.json", "--sign-key", seed], root, expect=1)
        assert "signature" in (tampered.stdout + tampered.stderr)
        sources = {}
        for relative in ["tools/hub-packs/hub_packs.py", "tools/hub-packs/ed25519.sh",
                         "tools/xerj-code/hub/validate_hub.py"]:
            sources[relative] = sha((hub / relative).read_bytes())
        for source in sorted((hub / "tools/corpus-hub-site").glob("*.py")):
            sources[str(source.relative_to(hub))] = sha(source.read_bytes())
        dump(output / "result.json", {"fixture": "Synthetic pins and opaque snapshots; no index restore or source execution",
                                       "binary_sha256": sha(binary.read_bytes()), "binary_version": run("version", [binary, "--version"], root).stdout.strip(),
                                       "hub_commit": subprocess.check_output(["git", "-C", str(hub), "rev-parse", "HEAD"], text=True).strip(),
                                       "main_commit": subprocess.check_output(["git", "-C", str(main_source), "rev-parse", "HEAD"], text=True).strip(),
                                       "main_gate_sha256": sha(current_gate.read_bytes()),
                                       "fixture_public_key": public,
                                       "source_sha256": sources, "results": results, "commands": records,
                                       "reproduced": True, "both_publication_gates_block_format_2": True,
                                       "deployed": False, "live_incident_claimed": False})
    print("Reproduced stale G7 in rc.101 signed catalogue and generated page. Both publication gates block format 2; format-1 control passes. Negative controls passed.")


if __name__ == "__main__":
    main()
