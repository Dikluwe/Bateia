"""Execute pinned real weights on known diagnostics, without training or threshold fitting."""
from copy import deepcopy
import hashlib
import importlib.metadata
from itertools import permutations
import json
from pathlib import Path
import platform
import statistics
import time
import torch
import laya

HERE = Path(__file__).resolve().parent
RESULTS = HERE / 'results'
pilot_path = HERE.parent / '0122/pilot-requests.json'
pilot = json.loads(pilot_path.read_text())
checkpoint = json.loads((RESULTS / 'checkpoint.json').read_text())
torch.set_num_threads(8)
torch.set_num_interop_threads(1)
torch.use_deterministic_algorithms(True)
loaded_at = time.perf_counter()
agent = laya.load(checkpoint['local_path'], device='cpu')
load_seconds = time.perf_counter() - loaded_at
records = []
order_keys = list(pilot['cases'][0]['request']['questions']['declaration_relation']['criteria'])
orders = list(permutations(order_keys))
assert tuple(order_keys) == orders[0]
with (RESULTS / 'predictions.jsonl').open('w') as output:
    for case in pilot['cases']:
        request = case['request']
        for index, order in enumerate(orders + [orders[0]]):
            question = deepcopy(request['questions'])
            original = question['declaration_relation']['criteria']
            question['declaration_relation']['criteria'] = {key: original[key] for key in order}
            start = time.perf_counter()
            result = agent.predict(request['state'], question)
            elapsed = time.perf_counter() - start
            usage = result['usage']
            assert usage['truncated'] is False and usage['state_tokens_dropped'] == 0
            assert not usage.get('options'), 'option budget collapsed options'
            answer = result['answers']['declaration_relation']
            assert answer['choice'] in original
            row = {'case_id': case['id'], 'variant': 'repeat_original' if index == 6 else 'order_' + str(index),
                   'option_order': list(order), 'expected_author_label': case['expected'],
                   'agrees_with_author_label': answer['choice'] == case['expected'],
                   'elapsed_seconds': elapsed, 'first_forward': not records,
                   'request_sha256': hashlib.sha256(json.dumps({'state': request['state'], 'questions': question}, separators=(',', ':')).encode()).hexdigest(),
                   'raw_result': result}
            output.write(json.dumps(row, allow_nan=False) + '\n')
            output.flush()
            records.append(row)
baseline = [r for r in records if r['variant'] == 'order_0']
comparison = []
for case in pilot['cases']:
    rows = [r for r in records if r['case_id'] == case['id']]
    first, repeated = rows[0]['raw_result']['answers'], rows[-1]['raw_result']['answers']
    comparison.append({'case_id': case['id'], 'author_label': case['expected'],
                       'baseline_choice': rows[0]['raw_result']['answers']['declaration_relation']['choice'],
                       'baseline_answer_confidence': rows[0]['raw_result']['answers']['declaration_relation'].get('answer_confidence'),
                       'baseline_agreement': rows[0]['agrees_with_author_label'],
                       'choices_across_orders': [r['raw_result']['answers']['declaration_relation']['choice'] for r in rows[:6]],
                       'distinct_choices_across_orders': len({r['raw_result']['answers']['declaration_relation']['choice'] for r in rows[:6]}),
                       'identical_repeat_answers': first == repeated})
summary = {'schema_version': 'laya-known-diagnostics-v1',
           'software_commit': '3cf26cbcb18725dbc2d127bb8bb2c4c43243ae63',
           'checkpoint_revision': checkpoint['revision'], 'checkpoint': 'ENGLISH_ROOT',
           'device': 'cpu', 'dtype': str(agent.dtype),
           'torch_threads': torch.get_num_threads(), 'deterministic_algorithms': True,
           'python': platform.python_version(), 'platform': platform.platform(),
           'packages': {name: importlib.metadata.version(name) for name in ('laya', 'torch', 'transformers', 'tokenizers', 'safetensors', 'huggingface_hub')},
           'reference_origin': 'AUTHOR_LABELED_KNOWN_DIAGNOSTIC_NOT_HOLDOUT',
           'independent_semantic_accuracy_measured': False,
           'training_performed': False, 'threshold_fitted': False,
           'min_confidence_applied': False, 'model_inference_executed': True,
           'unique_cases': 4, 'runs': len(records), 'option_orders_per_case': 6,
           'baseline_agrees_with_author': sum(r['agrees_with_author_label'] for r in baseline),
           'baseline_cases': len(baseline), 'load_seconds': load_seconds,
           'warm_exploratory_sdk_median_seconds': statistics.median(r['elapsed_seconds'] for r in records if not r['first_forward']),
           'no_state_or_option_truncation': True,
           'pilot_sha256': hashlib.sha256(pilot_path.read_bytes()).hexdigest(),
           'predictions_sha256': hashlib.sha256((RESULTS / 'predictions.jsonl').read_bytes()).hexdigest(),
           'cases': comparison,
           'limitations': ['Four known diagnostic cases, not a representative or independent test.',
                           'Permutations and repetition are correlated runs, not 28 examples.',
                           'No reference case labeled undetermined; abstention quality not measured.',
                           'SDK execution; no HTTP integration benchmark.']}
(RESULTS / 'diagnostics.json').write_text(json.dumps(summary, indent=2) + '\n')
print(json.dumps({'baseline_agreement': summary['baseline_agrees_with_author'],
                  'baseline_cases': 4, 'runs': len(records), 'cases': comparison}))
