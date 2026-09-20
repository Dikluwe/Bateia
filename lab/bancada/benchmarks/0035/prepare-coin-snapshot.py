from pathlib import Path
import shutil,shlex,json,re,subprocess,hashlib
source=Path('/mnt/Laranja/Git/externos/coin');snapshot=Path('/tmp/bateia-0035-coin');build=Path('/tmp/coin-contribution-review-build')
snapshot.mkdir(exist_ok=True)
for name in ['src','include']:
 shutil.copytree(source/name,snapshot/name,dirs_exist_ok=True)
for name in ['CMakeLists.txt','README.md','README']:
 if (source/name).is_file():shutil.copy2(source/name,snapshot/name)
entries=[];seen=set();provenance=[]
for flags in sorted(build.glob('src/**/CMakeFiles/*.dir/flags.make')):
 text=flags.read_text();args=[]
 for field in ['CXX_DEFINES','CXX_INCLUDES','CXX_FLAGS']:
  match=re.search('^'+field+r' = (.*)$',text,re.M)
  if match:args+=shlex.split(match[1])
 make=flags.with_name('build.make')
 if not make.exists():continue
 files=re.findall(re.escape(str(source))+r'/[^\s:]+\.cpp',make.read_text())
 for file in files:
  if file in seen:continue
  seen.add(file)
  mapped=file.replace(str(source),str(snapshot))
  entries.append({'directory':str(build),'file':mapped,'arguments':['/usr/bin/c++']+[arg.replace(str(source),str(snapshot)) for arg in args]+['-c',mapped]})
 provenance.append({'path':str(flags),'sha256':hashlib.sha256(flags.read_bytes()).hexdigest()})
(snapshot/'compile_commands.json').write_text(json.dumps(entries,indent=2)+'\n')
receipt={'original_root':str(source),'snapshot_root':str(snapshot),'base_commit':subprocess.check_output(['git','-C',str(source),'rev-parse','HEAD'],text=True).strip(),'compilation_units':len(entries),'configuration_source':'arguments reconstruídos de flags.make e caminhos de build.make existentes; nenhum build executado','flags_inputs':provenance}
Path('00_nucleo/benchmark/0035/coin-inputs.json').write_text(json.dumps(receipt,indent=2)+'\n')
print(len(entries),'compilation units')
