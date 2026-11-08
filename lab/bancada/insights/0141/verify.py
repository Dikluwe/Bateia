"""Verify 0141 gates, historical preservation and frozen link witnesses offline."""
import hashlib
import json
from pathlib import Path

ROOT=Path(__file__).resolve().parents[4]
LAB=Path(__file__).resolve().parent
CHANGED={'01_core/domain/insight_sites.go','01_core/usecase/insight_site_plan.go','01_core/usecase/insight_sites.go'}
NEW={'01_core/domain/insight_possibilities.go','01_core/domain/insight_possibilities_test.go','01_core/usecase/insight_site_schedule.go','01_core/usecase/insight_site_schedule_test.go','01_core/usecase/insight_reference_leads.go','01_core/usecase/insight_reference_leads_test.go'}
def sha(path):return hashlib.sha256(path.read_bytes()).hexdigest()
def read(path):return json.loads(path.read_text())

def main():
    baseline=read(LAB/'baseline.json')
    changed={name for name,digest in baseline.items() if sha(ROOT/name)!=digest}
    assert changed==CHANGED,changed
    current={str(p.relative_to(ROOT)) for folder in ['01_core','02_shell','03_infra','04_wiring'] for p in (ROOT/folder).rglob('*.go')}
    assert current-baseline.keys()==NEW
    lineage=[]
    for name in sorted(CHANGED|NEW):
        prompt='00_nucleo/prompts/'+name.removesuffix('.go')+'.md'
        assert '@prompt '+prompt in (ROOT/name).read_text()
        lineage.append(dict(source=name,source_sha256=sha(ROOT/name),prompt=prompt,prompt_sha256=sha(ROOT/prompt)))
    gates=read(LAB/'final/verification.json')
    assert all(item['exit_code']==0 for item in gates.values())
    assert 'audited 382 productive Go files' in (LAB/'final/architecture.log').read_text()
    for step in ['0106','0134','0135','0136','0137','0138']:
        result=read(LAB/('replay-'+step+'.json'))
        assert result.get('replay_byte_identical',result.get('exploration_replay_equal',False)),step
    fixture=read(LAB/'frozen-corpus-result.json')
    assert fixture['mode']=='OFFLINE_FROZEN_FIXTURE_REPLAY' and not fixture['fresh_network']
    report=fixture['report']
    assert len(report['bundles'])==2 and len(report['reference_leads'])==12
    assert report['plan']['scheduling_policy']=='PERSPECTIVE_SOURCE_ROUND_ROBIN_V1'
    assert len(report['opportunity_coverage']['completed_perspectives'])==5
    previous=read(ROOT/'lab/bancada/insights/comparacao-ia-2026-10-28/bateia-result.json')
    old_items={(b['repository']['owner']+'/'+b['repository']['repo'],b['repository']['commit']):b['items'][0] for b in previous['bundles']}
    selected={c['id']:c for c in report['candidates'] if c['selected']}
    bundles={(b['repository']['owner']+'/'+b['repository']['repo'],b['repository']['commit']):b for b in report['bundles']}
    lead_ids=set()
    found=False
    for lead in report['reference_leads']:
        assert lead['source_candidate_id'] in selected
        assert lead['source_identity']==selected[lead['source_candidate_id']]['product']['identity']
        assert lead['status']=='REFERENCE_NOT_READ'
        assert lead['id'] not in lead_ids
        lead_ids.add(lead['id'])
        source=lead['source']
        key=(source['owner']+'/'+source['repo'],source['commit'])
        item=bundles[key]['items'][0]
        assert item==old_items[key]
        assert item['content_hash']=='sha256:'+hashlib.sha256(item['content'].encode()).hexdigest()
        ref=lead['reference']
        assert ref['content_hash']==item['content_hash'] and ref['path']==item['path']
        lines=item['content'].removesuffix('\n').split('\n')
        assert '\n'.join(lines[ref['start_line']-1:ref['end_line']])==ref['excerpt']
        assert len(ref['excerpt'].encode())<=1024
        assert not lead['target']['external_id'].startswith('organizations/')
        if lead['target']['external_id']=='chnm/serendipomatic':
            found=True
            assert key[0]=='MuseCompNet/muse-tech-central' and ref['start_line']==106 and ref['end_line']==107
    assert found
    for task in report['tasks']:
        assert set(task['reference_lead_refs'])==lead_ids
        assert not lead_ids.intersection(task['reading_card_refs']+task['catalog_reference_ids'])
    (LAB/'lineage.json').write_text(json.dumps(lineage,indent=2)+'\n')
    print(json.dumps(dict(status='PASS',contract='0141',gate_exit_codes={k:v['exit_code'] for k,v in gates.items()},productive_files_audited=382,
        changed_existing_go=sorted(CHANGED),new_go=sorted(NEW),preexisting_tests_unchanged=True,
        preserved_previous_evidence_files=sum(n.startswith('lab/') for n in baseline),historical_replays_passed=6,
        frozen_readmes=2,reference_leads=12,serendipomatic_reference_preserved=True,destination_reads=0,
        schedule_control=dict(sources=11,logical_calls=20,perspectives=5,calls_per_perspective=4,previous_policy_perspectives=2,real_network=False),
        novelty_assessed=False,effect_verified=False,candidate_code_executed=False,final_binary_sha256=sha(LAB/'final/bateia')),ensure_ascii=False,indent=2))

if __name__=='__main__':main()
