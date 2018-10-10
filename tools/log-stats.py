import json
import re
import jsbeautifier
import sys
import os

if len(sys.argv) < 3:
  print("Not enough arguements provided")
  sys.exit(1)

log_file = sys.argv[1]
directory = sys.argv[2]
print("Write %s data to %s" % (log_file, directory))

# Read target file
with open(log_file, 'r') as f:
  content = f.readlines()

# Load to json
j = [json.loads(line) for line in content]
d = {}
unlabeled = []

# Group by taskID
for line in j:
  if "taskID" in line:
    key = line["taskID"]
    if key not in d:
      d[key] = []
    d[key].append(line)
  else:
    unlabeled.append(line)

# Create output folder
os.makedirs(os.path.dirname("./%s/" % directory), exist_ok=True)

# Write json files and beautify
for k, v in d.items():
  with open("./%s/%s.json" % (directory, k), 'w') as f:
    f.write(jsbeautifier.beautify(str(json.dumps(v))))

with open("./%s/unlabled.json" % directory, 'w') as f:
  f.write(jsbeautifier.beautify(str(json.dumps(unlabeled))))
