"""Validate the frozen multi-site run and replay only its reading context offline."""
import argparse
import copy
import hashlib
import json
from pathlib import Path
import subprocess
import tempfile


def sha(value):
    return hashlib.sha256(json.dumps(value,ensure_ascii=False,separators=(',',':')).encode()).hexdigest()


def main():
    parser=argparse.ArgumentParser()
    parser.add_argument('--binary',required=True)
    args=parser.parse_args()
    root=Path(__file__).resolve().parent
    report=json.loads((root/'online-result.json').read_text())
    receipt=json.loads((root/'online-check.json').read_text())
    assert receipt['no_credentials'] and receipt['isolated_cwd_without_dotenv']
    assert not receipt['candidate_code_executed'] and receipt['exit_code']==0
    assert not (root/'online-stdout.log').read_bytes()
    assert report['schema_version']=='insight-site-discovery-v1'
    assert report['status']=='PARTIAL' and report['novelty_state']=='NOT_ASSESSED' and report['effect_state']=='UNVERIFIED'
    plan=report['plan']
    canonical_plan=copy.deepcopy(plan)
    canonical_plan['id']=''
    assert plan['id']=='site_plan_'+sha(canonical_plan)
    sources={cap['id']:cap for cap in plan['sources']}
    assert len(sources)==11 and 'gitee-repository-catalog' in sources and 'atomgit-repository-catalog' in sources
    queries={q['id'] for q in plan['discovery']['queries']}
    assert len(queries)==5 and len(report['attempts'])==55
    assert [(a['query_id'],a['connector_id']) for a in report['attempts']]==[(q['id'],s['id']) for q in plan['discovery']['queries'] for s in plan['sources']]
    omitted={'NOT_STARTED','AUTH_BLOCKED','UNSUPPORTED_LOCALE','UNSUPPORTED_MARKET','BUDGET_EXHAUSTED','CANDIDATE_BUDGET'}
    calls=sum(a['status'] not in omitted for a in report['attempts'])
    assert calls==report['calls']==12 and calls<=plan['request']['max_calls']
    atom=[a['status'] for a in report['attempts'] if a['connector_id']=='atomgit-repository-catalog']
    gitee=[a['status'] for a in report['attempts'] if a['connector_id']=='gitee-repository-catalog']
    assert atom==['AUTH_REQUIRED']+['AUTH_BLOCKED']*4
    assert gitee[0]=='TIMED_OUT'
    candidates={c['id']:c for c in report['candidates']}
    assert len(candidates)==3
    selected={key:c for key,c in candidates.items() if c['selected']}
    assert len(selected)<=plan['request']['discovery']['max_references']
    for key,candidate in candidates.items():
        product=candidate['product']
        identity=product['identity']
        assert identity['connector_id'] in sources
        assert key=='catalog_ref_'+sha(dict(Connector=identity['connector_id'],External=identity['external_id']))
        assert candidate['origin']=='CATALOG_METADATA_UNVERIFIED'
        assert set(candidate['query_ids'])<=queries
        for origin in product['origins']:
            assert origin['connector_id']==identity['connector_id']
            assert origin['connector_version']==sources[identity['connector_id']]['version']
            assert origin['query_id'] in queries
            matches=[a for a in report['attempts'] if a['connector_id']==origin['connector_id'] and a['query_id']==origin['query_id']]
            assert len(matches)==1 and matches[0]['response_hash']==origin['response_hash']
            assert matches[0]['market']==origin['observed_market']
    readings={card['id']:card for card in report['exploration']['cards']}
    assert len(report['bundles'])==len(readings)==1
    for task in report['tasks']:
        assert set(task['catalog_reference_ids'])==set(selected)
        assert set(task['reading_card_refs'])<=set(readings)
        assert not (set(task['catalog_reference_ids']) & set(task['reading_card_refs']))
    for readme in report['readmes']:
        assert readme['candidate_id'] in selected
        if readme['status']=='COLLECTED':
            assert readme['connector_id']=='github-repository-catalog'
            bundle=report['bundles'][0]
            assert not bundle['complete'] and readme['receipt']['source']==bundle['repository']
            for item in bundle['items']:
                assert item['content_hash']=='sha256:'+hashlib.sha256(item['content'].encode()).hexdigest()
        else:
            assert readme['status']=='CATALOG_ONLY'
    with tempfile.TemporaryDirectory(prefix='bateia-site-context-replay-') as directory:
        path=Path(directory)/'bundles.json'
        path.write_text(json.dumps(report['bundles'],ensure_ascii=False))
        request=plan['request']['discovery']
        command=[str(Path(args.binary).resolve()),'idea','insights','explore','--question',request['question'],'--bundles',str(path),'--max-passages',str(request['max_passages'])]
        for familiar in request['familiar_approaches']:command+=['--familiar',familiar]
        replay=subprocess.run(command,capture_output=True,timeout=30)
        assert replay.returncode==0,replay.stderr.decode()
        assert json.loads(replay.stdout)==report['exploration']
    print(json.dumps(dict(schema_version='1',result_sha256=hashlib.sha256((root/'online-result.json').read_bytes()).hexdigest(),sources=11,planned_cells=55,logical_calls=12,unique_candidates=3,readme_bundles=1,report_status='PARTIAL',atomgit_without_token='AUTH_REQUIRED_THEN_AUTH_BLOCKED',gitee_public_run='TIMED_OUT',source_scoped_ids_verified=True,provenance_refs_resolve=True,catalog_and_reading_refs_separated=True,exploration_replay_equal=True,credentials_used=False,candidate_code_executed=False,semantic_relevance_assessed=False,insight_utility_measured=False),ensure_ascii=False,indent=2))


if __name__=='__main__':
    main()
