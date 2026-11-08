"""Check the frozen public trial and replay its exploration offline."""
import argparse
import hashlib
import json
from pathlib import Path
import re
import subprocess
import tempfile


def canonical(value):
    return json.dumps(value, ensure_ascii=False, separators=(',', ':')).encode()


def digest(value):
    return hashlib.sha256(canonical(value)).hexdigest()


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--binary', required=True)
    args = parser.parse_args()
    binary = str(Path(args.binary).resolve())
    root = Path(__file__).resolve().parent
    original = (root/'online-result.json').read_bytes()
    report = json.loads(original)
    run = json.loads((root/'online-check.json').read_text())
    assert run['public_no_github_token'] and run['isolated_cwd_without_dotenv']
    assert run['exit_code'] == 0 and not run['candidate_code_executed']
    assert not (root/'online-stdout.log').read_bytes()
    assert report['schema_version'] == 'insight-discovery-v1'
    assert report['novelty_state'] == 'NOT_ASSESSED' and report['effect_state'] == 'UNVERIFIED'
    plan = report['plan']
    assert plan['request_sha256'] == digest(plan['request'])
    assert plan['id'] == 'discovery_plan_'+digest(dict(Contract='insight-discovery-v1', Input=plan['request_sha256']))
    assert len(plan['queries']) == len(report['attempts']) == 5
    query_ids = set()
    for query in plan['queries']:
        assert query['derivation'] == 'LEXICAL_EXPLORATION_HEURISTIC'
        assert query['id'] == 'discovery_query_'+digest(dict(Plan=plan['id'], Perspective=query['perspective'], Query=query['query']))
        query_ids.add(query['id'])
    assert len(query_ids) == 5
    assert {attempt['query_id'] for attempt in report['attempts']} == query_ids
    assert all(a['status'] == 'SUCCEEDED' for a in report['attempts'])
    refs = {reference['id']: reference for reference in report['candidates']}
    selected = {key: ref for key, ref in refs.items() if ref['selected']}
    assert len(refs) == 4 and len(selected) == 2
    assert len({ref['repository']['owner'].lower() for ref in selected.values()}) == 2
    for ref in refs.values():
        assert set(ref['query_ids']) <= query_ids
        assert ref['id'] == 'reference_'+digest((ref['repository']['owner']+'/'+ref['repository']['repo']).lower())
    assert {receipt['reference_id'] for receipt in report['readmes']} == set(selected)
    assert sorted(receipt['status'] for receipt in report['readmes']) == ['COLLECTED', 'CONTENT_LIMIT']
    assert report['status'] == 'PARTIAL'
    receipts = {receipt['reference_id']: receipt for receipt in report['readmes']}
    bundle_items = {}
    for bundle in report['bundles']:
        source = bundle['repository']
        assert re.fullmatch('[0-9a-f]{40}|[0-9a-f]{64}', source['commit'])
        assert not bundle['complete'] and not bundle['tree_complete']
        assert len(bundle['items']) == 1
        reference_id = 'reference_'+digest((source['owner']+'/'+source['repo']).lower())
        assert reference_id in selected and not selected[reference_id]['revision_conflict']
        assert source['commit'] in selected[reference_id]['observed_commits']
        receipt = receipts[reference_id]
        assert receipt['status'] == 'COLLECTED' and receipt['scope'] == 'ROOT_README_ONLY' and receipt['source'] == source
        for item in bundle['items']:
            assert item['ref'] == source and item['path'] == receipt['path']
            assert 0 < len(item['content'].encode()) == receipt['bytes'] <= 32 << 10
            assert item['content_hash'] == 'sha256:'+hashlib.sha256(item['content'].encode()).hexdigest()
            bundle_items[(source['owner']+'/'+source['repo'], source['commit'], item['path'])] = item
    assert len(bundle_items) == 1
    exploration = report['exploration']
    assert exploration['status'] == 'CONTEXT_PREPARED'
    assert len(exploration['cards']) == 1
    for card in exploration['cards']:
        source, ref = card['source'], card['reference']
        item = bundle_items[(source['repository'], source['commit'], ref['path'])]
        lines = item['content'].removesuffix('\n').split('\n')
        assert ref['content_hash'] == item['content_hash']
        assert ref['excerpt'] == '\n'.join(lines[ref['start_line']-1:ref['end_line']])
    with tempfile.TemporaryDirectory(prefix='bateia-discovery-replay-') as directory:
        bundles_path = Path(directory)/'bundles.json'
        bundles_path.write_text(json.dumps(report['bundles'], ensure_ascii=False))
        command = [binary, 'idea', 'insights', 'explore', '--question', plan['request']['question'], '--bundles', str(bundles_path), '--max-passages', str(plan['request']['max_passages'])]
        for label in plan['request']['familiar_approaches']:
            command += ['--familiar', label]
        replay = subprocess.run(command, capture_output=True, timeout=30)
        assert replay.returncode == 0, replay.stderr.decode()
        assert json.loads(replay.stdout) == exploration
    print(json.dumps(dict(
        schema_version='1', public_result_sha256=hashlib.sha256(original).hexdigest(),
        queries_attempted=5, queries_succeeded=5, unique_candidates=4, selected_references=2,
        readmes_collected=1, readmes_rejected_by_content_limit=1, report_status='PARTIAL',
        plan_identity_verified=True, origin_refs_resolve=True,
        pinned_readme_hash_and_exact_excerpt_verified=True, exploration_replay_equal=True,
        fresh_online_search_replayed=False, semantic_diversity_assessed=False,
        novelty_assessed=False, effect_verified=False, candidate_code_executed=False,
    ), ensure_ascii=False, indent=2))


if __name__ == '__main__':
    main()
