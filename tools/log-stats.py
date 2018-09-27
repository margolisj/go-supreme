import json
import re
import jsbeautifier

log_file = "logs/logfile-1538060396.log"

with open(log_file, 'r') as f:
  content = f.readlines()

j = [json.loads(line) for line in content]
d = {}
unlabeled = []

for line in j:
  match = re.match("^(\d{1,2})\s" , line['msg'])
  if match:
    key = match.group(1)
    if key not in d:
      d[key] = []
    d[key].append(line)
  elif 'thread' in line:
    key = str(line['thread'])
    if key not in d:
      d[key] = []
    d[key].append(line)
  else:
    unlabeled.append(line)

for k, v in d.items():
  with open("./tempOutput/%s.json" % k, 'w') as f:
    f.write(jsbeautifier.beautify(str(json.dumps(v))))

with open("./tempOutput/unlabled.json", 'w') as f:
  f.write(jsbeautifier.beautify(str(json.dumps(unlabeled))))
# print(d)
