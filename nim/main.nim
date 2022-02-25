#[ Nim try out and test bench
  - Number examples: 1_000_000, 1.0e9, 
]#

import std/sugar
from std/strutils import parseInt

echo "OS: ", system.hostOS, "; CPU: ", system.hostCPU
echo "Give me an interval. Start: "
let st = parseInt(readLine(stdin))
echo "End: "
let en = parseInt(readLine(stdin))

assert st < en and en < 10000

let all = collect(newSeq):
    for i in countup(st, en):
        if (i mod 2) == 0:
            i

echo "Range", all