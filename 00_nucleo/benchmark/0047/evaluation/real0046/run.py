#!/usr/bin/env python3
"""Offline source-only retrieval evaluation. Run from any cwd: python3 run.py baseline|candidate.
Reads immutable module ZIPs, never imports/builds/runs external project code.
Temporary Go adapter imports only Bateia public domain/parser API.
"""
import argparse,datetime,hashlib,json,os,pathlib,subprocess,tempfile,zipfile
HERE=pathlib.Path(__file__).resolve().parent
ROOT=HERE.parents[4]
sha=lambda b:hashlib.sha256(b).hexdigest()
def read(n): return json.loads((HERE/n).read_text())
def main():
 ap=argparse.ArgumentParser(); ap.add_argument('label',choices=['baseline','candidate','final']); args=ap.parse_args()
 freeze=read('freeze.json')
 for name,digest in freeze['files'].items():
  assert sha((HERE/name).read_bytes())==digest, 'Frozen input changed: '+name
 manifest=read('manifest.json'); questions=read('questions.json'); budget=manifest['budgets']
 projects={}; scans={}
 for p in manifest['projects']:
  z=pathlib.Path(p['zip_path']); assert sha(z.read_bytes())==p['zip_sha256']
  prefix=p['module']+'@'+p['version']+'/'
  accepted=[]; skipped=[]; size=0
  with zipfile.ZipFile(z) as archive:
   for f in p['files']:
    data=archive.read(prefix+f['path']); assert sha(data)==f['sha256'] and len(data)==f['bytes']
    if len(data)>budget['max_file_bytes'] or len(accepted)>=budget['max_files'] or size+len(data)>budget['max_bytes']:
     skipped.append(f['path']); continue
    accepted.append({'path':f['path'],'content':data.decode(),'size':len(data)}); size+=len(data)
  projects[p['id']]=accepted
  scans[p['id']]={'files_available':len(p['files']),'files_passed_to_rank':len(accepted),'bytes_passed_to_rank':size,'skipped_budget_paths':skipped,'limited':bool(skipped)}
 payload=[{'id':q['id'],'query':q['query'],'files':projects[q['project']],'limit':budget['max_results']} for q in questions]
 source='''package main
import("encoding/json";"os";"bateia/01_core/domain";"bateia/01_core/parser")
type Input struct {ID string `json:"id"`; Query string `json:"query"`; Files []domain.ExtractedFile `json:"files"`; Limit int `json:"limit"`}
type Output struct {ID string `json:"id"`; Report domain.ImplementationSearchReport `json:"report"`}
func main(){var in []Input;if err:=json.NewDecoder(os.Stdin).Decode(&in);err!=nil{panic(err)};out:=[]Output{};for _,v:=range in{out=append(out,Output{v.ID,parser.RankGoImplementations(v.Query,v.Files,v.Limit)})};if err:=json.NewEncoder(os.Stdout).Encode(out);err!=nil{panic(err)}}
'''
 source_hashes={str(p.relative_to(ROOT)):sha(p.read_bytes()) for base in ['01_core/parser','01_core/domain'] for p in sorted((ROOT/base).glob('*.go')) if not p.name.endswith('_test.go')}
 with tempfile.TemporaryDirectory(prefix='bateia0046-real-') as d:
  adapter=pathlib.Path(d)/'main.go'; adapter.write_text(source)
  env=dict(os.environ,GOPROXY='off',GOSUMDB='off',GOCACHE='/tmp/bateia-0047-cache')
  result=subprocess.run(['go','run',str(adapter)],cwd=ROOT,env=env,input=json.dumps(payload),text=True,capture_output=True,check=True)
 reports={o['id']:o['report'] for o in json.loads(result.stdout)}
 rows=[]
 for q in questions:
  r=reports[q['id']]; paths=list(dict.fromkeys(h['path'] for h in (r['hits'] or []))); expected=set(e['path'] for e in q['expected']); examined={f['path'] for f in (r['examined'] or [])}; row={'id':q['id'],'project':q['project'],'query':q['query'],'negative':bool(q.get('negative')),'expected_paths':sorted(expected),'ranked_paths':paths,'hit_count':len(r['hits'] or []),'examined_count':len(examined),'missing_from_scan':sorted(expected-{f['path'] for f in projects[q['project']]}),'missing_from_rank':sorted(expected-set(paths)),'negative_indications_returned':bool(paths) if q.get('negative') else None}
  for k in [3,12]:
   selected=paths[:k]; found=expected & set(selected)
   row['at_'+str(k)]={'paths':selected,'labeled_expected_found':len(found),'recall_labeled':len(found)/len(expected) if expected else None,'precision_labeled':len(found)/len(selected) if selected and expected else None,'missed_expected':sorted(expected-set(selected))}
  rows.append(row)
 positives=[r for r in rows if not r['negative']]; negatives=[r for r in rows if r['negative']]
 aggregate={'positive_requirements':len(positives),'negative_requirements':len(negatives),'negative_with_indications':sum(r['negative_indications_returned'] for r in negatives)}
 for k in [3,12]:
  aggregate['at_'+str(k)]={'macro_recall_labeled':sum(r['at_'+str(k)]['recall_labeled'] for r in positives)/len(positives),'micro_precision_labeled':sum(r['at_'+str(k)]['labeled_expected_found'] for r in positives)/max(1,sum(len(r['at_'+str(k)]['paths']) for r in positives))}
 output={'label':args.label,'executed_at_utc':datetime.datetime.now(datetime.timezone.utc).isoformat(),'input_freeze_sha256':sha((HERE/'freeze.json').read_bytes()),'runner_sha256':sha(pathlib.Path(__file__).read_bytes()),'adapter_sha256':sha(source.encode()),'production_source_sha256':source_hashes,'go_version':subprocess.check_output(['go','version'],text=True).strip(),'budgets':budget,'scan':scans,'aggregate':aggregate,'requirements':rows,'raw_reports':reports,'limitations':['Rank public API only; scan budgets enforced by adapter before Rank. No provider scan, extraction reread, hash revalidation or CLI measured here.','at_3/at_12 mean first distinct paths in 24-hit rank report, simulated extraction slots, not actual extracted content.','Precision only against manually labeled expected files; other returned files not proven irrelevant.','Negative labels are manual absence in stated implementation scope, not universal absence.']}
 target=HERE/(args.label+'.json')
 if target.exists(): raise SystemExit('Refusing overwrite of '+str(target))
 target.write_text(json.dumps(output,indent=2,ensure_ascii=False)+'\n'); print(json.dumps(aggregate,indent=2)); print(str(target))
if __name__=='__main__':main()
