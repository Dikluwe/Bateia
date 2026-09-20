"""Reevaluate preserved observations without repeating network or CLI collection."""
import argparse
import json
from pathlib import Path
import runner

def main():
    p=argparse.ArgumentParser()
    p.add_argument('--input',required=True,type=Path)
    p.add_argument('--output',required=True,type=Path)
    p.add_argument('--corpus',type=Path,default=Path(__file__).with_name('corpus.json'))
    p.add_argument('--freeze',type=Path,default=Path(__file__).with_name('freeze.json'))
    a=p.parse_args()
    assert not a.output.exists(), 'output must be new'
    assert runner.sha256(a.corpus.read_bytes())==json.loads(a.freeze.read_text())['corpus_sha256']
    artifacts=json.loads((a.input/'artifacts.json').read_text())
    for name,digest in artifacts.items():
        name=runner.safe_path(name)
        assert runner.sha256((a.input/name).read_bytes())==digest, name
    corpus=json.loads(a.corpus.read_text())
    projects={p['id']:p for p in corpus['projects']}
    results=[]
    for i,case in enumerate(corpus['cases'],1):
        d=a.input/f'case-{i:02}'
        prior=json.loads((d/'evaluation.json').read_text())
        project=projects[case['project_id']]
        snapshot,inventory=runner.checked_inventory(project)
        def read(name):
            f=d/name
            return json.loads(f.read_text()) if f.exists() else None
        result=runner.evaluate_case(case,project,read('raw.stdout'),read('cli.json'),snapshot,
            errors=prior['errors'],inventory_valid=inventory['valid'] and prior['inventory_valid'],discovery=read('discovery.stdout'))
        results.append(result)
    a.output.mkdir(parents=True)
    runner.write_json(a.output/'cases.json',results)
    runner.write_json(a.output/'summary.json',runner.summarize(results))
    runner.write_json(a.output/'receipt.json',{'source':str(a.input),'artifacts_sha256':runner.sha256((a.input/'artifacts.json').read_bytes()),'corpus_sha256':runner.sha256(a.corpus.read_bytes()),'runner_sha256':runner.sha256(Path(runner.__file__).read_bytes()),'network_repeated':False})

if __name__=='__main__':main()
