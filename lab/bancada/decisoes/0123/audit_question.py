"""Check question instruction/option budgets separately from state truncation flags."""
import json
from pathlib import Path
from transformers import AutoTokenizer
from laya.common import _encode_question_text, render_options

HERE = Path(__file__).resolve().parent
pilot = json.loads((HERE.parent / '0122/pilot-requests.json').read_text())
checkpoint = Path(json.loads((HERE / 'results/checkpoint.json').read_text())['local_path'])
config = json.loads((checkpoint / 'rl_agent_config.json').read_text())
tokenizer = AutoTokenizer.from_pretrained(checkpoint / 'tokenizer', local_files_only=True)
definition = pilot['cases'][0]['request']['questions']['declaration_relation']
q = {'t': definition['type'], 'ins': definition['instructions'], 'crit': definition['criteria']}
instruction_ids = _encode_question_text(tokenizer, 'choice question: ' + q['ins'], add_special_tokens=False)
lengths = [len(_encode_question_text(tokenizer, ' ' + text, add_special_tokens=False))
           for text in render_options(q)]
option_cost = sum(1 + min(length, 48) for length in lengths)
budget = config['head_max_len'] - option_cost
assert budget >= 16 and len(instruction_ids) <= max(8, budget)
assert all(length <= 48 for length in lengths)
report = {'head_max_len': config['head_max_len'], 'max_len': config['max_len'],
          'instruction_tokens': len(instruction_ids), 'option_tokens': lengths,
          'instruction_budget': budget, 'no_instruction_or_option_text_clipping': True,
          'question_types_and_counts': {'choice': 3},
          'checkpoint_temperature_choice_3_5': config['temperature_by_options']['choice:3-5'],
          'checkpoint_warning': 'choice:11+ temperature was clamped at load; diagnostics use choice:3-5.',
          'domain_calibration_performed': False}
(HERE / 'results/question-budget.json').write_text(json.dumps(report, indent=2) + '\n')
print(json.dumps(report))
