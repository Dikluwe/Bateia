"""Compare frozen summaries and audit diagnostics; never infer semantic success."""
import argparse,json
from pathlib import Path
p=argparse.ArgumentParser();p.add_argument('--input',required=True);p.add_argument('--output',required=True);a=p.parse_args()
r=Path(a.input);corpus=json.loads(Path('00_nucleo/benchmark/0053/corpus.json').read_text());projects={v['id']:v['repo'] for v in corpus['projects']}
old=json.loads(Path('00_nucleo/benchmark/0055/collected-r2/summary.json').read_text());new=json.loads((r/'summary.json').read_text());out={'baseline':'0055/collected-r2','local_metrics_unchanged':old['metrics']==new['metrics'],'discovery_before':old['discovery'],'discovery_after':new['discovery'],'diagnostics':[]}
for i,c in enumerate(corpus['cases'],1):
 path=r/('case-%02d'%i)/'discovery.stdout'
 if not path.exists():continue
 raw=json.loads(path.read_text());attempts=raw.get('attempts') or [];failed=bool(raw.get('error')) or any(x.get('error') for x in attempts);expected=projects[c['project_id']]
 for n,t in enumerate(attempts):
  if not t.get('coverage'):continue
  ds=t.get('candidate_diagnostics') or [];target=next((d for d in ds if d['repository'].lower()==expected.lower()),None)
  out['diagnostics'].append({'case':c['id'],'attempt':n,'status':'unknown' if failed else 'observed','coverage':t['coverage'],'expected_repo':expected,'expected_in_examined':None if failed else target is not None,'expected_diagnostic':target,'candidates':[{'repository':d['repository'],'status':d['status'],'reason':d['reason'],'action_found':d['action_found'],'object_found':d['object_found'],'truncated_snippets':d['truncated_snippets'],'omitted_snippets':d['omitted_snippets']} for d in ds]})
Path(a.output).write_text(json.dumps(out,indent=2,ensure_ascii=False)+'\n')
