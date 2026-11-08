"""Check the recorded pilot; no network or candidate execution."""
import collections
import hashlib
import json
from pathlib import Path

root=Path(__file__).resolve().parent
repo=root.parents[3]
def read(name):return json.loads((root/name).read_text())
def sha(path):return hashlib.sha256(path.read_bytes()).hexdigest()
freeze=read('direct-freeze.json')
assert sha(root/freeze['file'])==freeze['sha256']
run=read('bateia-run.json')
assert run['direct_arm_sha256']==freeze['sha256']
assert run['credentials_used'] is False and run['candidate_code_executed'] is False
assert sha(root/'bateia-result.json')==run['result_sha256']
report=read('bateia-result.json')
assert report['schema_version']=='insight-site-discovery-v1'
assert report['novelty_state']=='NOT_ASSESSED' and report['effect_state']=='UNVERIFIED'
assert report['calls']<=20
queries={q['id']:q['perspective'] for q in report['plan']['discovery']['queries']}
assert len(queries)==5
sources={s['id'] for s in report['plan']['sources']}
assert len(report['attempts'])==len(sources)*5
omitted={'NOT_STARTED','AUTH_BLOCKED','UNSUPPORTED_LOCALE','UNSUPPORTED_MARKET','BUDGET_EXHAUSTED','CANDIDATE_BUDGET'}
executed=[a for a in report['attempts'] if a['status'] not in omitted]
assert len(executed)==report['calls']
candidates={c['id']:c for c in report['candidates']}
selected={key:c for key,c in candidates.items() if c['selected']}
assert len(selected)<=6
all_cards={}
for card in report['exploration']['cards']:all_cards[card['id']]=card
bundles=list(report['bundles'])
for scoped in report.get('site_contexts',[]):
    assert scoped['candidate_id'] in selected
    assert scoped['identity']==selected[scoped['candidate_id']]['product']['identity']
    bundles.append(scoped['bundle'])
    for card in scoped['cards']:
        assert card['id'] not in all_cards
        all_cards[card['id']]=card
        assert card['source']['repository']==scoped['identity']['external_id']
        assert card['source']['commit']==scoped['bundle']['repository']['commit']
assert len(all_cards)<=6
items={}
for bundle in bundles:
    ref=bundle['repository']
    assert len(ref['commit']) in (40,64)
    assert not bundle['complete']
    for item in bundle['items']:
        assert item['ref']==ref
        assert item['content_hash']=='sha256:'+hashlib.sha256(item['content'].encode()).hexdigest()
        items[(ref['owner']+'/'+ref['repo'],ref['commit'],item['path'],item['content_hash'])]=item
for card in all_cards.values():
    key=(card['source']['repository'],card['source']['commit'],card['reference']['path'],card['reference']['content_hash'])
    assert key in items
    ref=card['reference']
    content=items[key]['content'].replace('\r\n','\n').replace('\r','\n').split('\n')
    assert '\n'.join(content[ref['start_line']-1:ref['end_line']])==ref['excerpt']
for task in report['tasks']:
    assert set(task['reading_card_refs'])==set(all_cards)
    assert set(task['catalog_reference_ids'])==set(selected)
for name,digest in read('baseline-sha256.json').items():
    assert sha(repo/name)==digest,name
summary=dict(status='PASS',sources_planned=len(sources),planned_cells=len(report['attempts']),logical_catalog_calls=report['calls'],
    perspectives_attempted=dict(collections.Counter(queries[a['query_id']] for a in executed)),
    attempts_by_status=dict(collections.Counter(a['status'] for a in report['attempts'])),
    candidates=len(candidates),selected=len(selected),readme_bundles=len(bundles),reading_cards=len(all_cards),
    readmes_by_status=dict(collections.Counter(r['status'] for r in report['readmes'])),
    bytes_collected=sum(len(item['content'].encode()) for b in bundles for item in b['items']),
    direct_arm_frozen_before_bateia=True,product_and_nucleus_unchanged=True,
    network_or_auth_failures_are_not_relevance_judgments=True,
    measured_insight_utility=False,independent_evaluator=False)
print(json.dumps(summary,ensure_ascii=False,indent=2))
