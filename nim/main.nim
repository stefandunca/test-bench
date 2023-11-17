#[ Nim try out and test bench
  - Number examples: 1_000_000, 1.0e9, 
]#

import std/sugar
import json

# Generate JsonNode from a string inline
var j = parseJson("""
  {
    "name": "John",
    "age": 30,
    "cars": [
      { "name":"Ford", "models":[ "Fiesta", "Focus", "Mustang" ] },
      { "name":"BMW", "models":[ "320", "X3", "X5" ] },
      { "name":"Fiat", "models":[ "500", "Panda" ] }
    ]
  }
""")

proc print_two(x: string, y: string) =
  echo "x = \"" & x & "\""
  echo "y = \"" & y & "\""

print_two($j["name"], $j["cars"])