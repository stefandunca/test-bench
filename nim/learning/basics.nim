echo "-------------- runtime ----------------------"

# Demo of compile time `const`
proc fillString(message: string): string =
  result = ""
  echo message
  for i in 0 .. 4:
    result.add($i)

proc constCompileTimeDemo() =
  const compileCount = fillString("Generating string - Compile time message")
  echo "@dd compileCount: ", compileCount
  let runtimeCount = fillString("Generating string - Runtime time message")
  echo "@dd runtimeCount: ", runtimeCount

proc caseAndStyleInsensitiveDemo() =
  let
    styleInsensitive = "Style insensitive variable"
    caseInsensitive = 123
  echo "@dd `styleInsensitive` used as `style_insensitive`: ", style_insensitive
  echo "@dd `caseInsensitive` used as `caseinsensitive`: ", caseinsensitive

proc oneLineProc: string = "One line proc"
proc simplestProc: bool = echo "Simplest proc"; false
proc discardedReturnedProc: string = discard "Discarded string"
proc autoReturnedProc: auto = 5.6
proc asertingProc =
  # won't raise exception if compiled with assertions off: `nim c --assertions:off`
  assert true, "Will be disabled in release mode"
  # will raise an exception also when compiled with assertions off: `nim c --assertions:off`
  doAssert true, "Triggered in release mode"
proc proceduresDemo() =
  echo oneLineProc()
  discard simplestProc()
  let discardedReturn = discardedReturnedProc()
  echo "Discarded return: " & repr(discardedReturn)
  echo "Auto return: ", autoReturnedProc()
  asertingProc()

proc containersDemo() =
  var arr: array[-5 .. 3, int] = [-3, -2, -1, 1, 3, 6, 9, 11, 12]
  echo "@dd `arr` is: ", arr[(arr.low + 2)..2]
  var seqTest: seq[float] = @[1.1, 2.2, 3.3, 4.4, 5.5]
  seqTest.add(6.6)
  echo "@dd `seqTest` is: ", seqTest[3..seqTest.high]
  let bToH: set[char] = {'B' .. 'h'}
  echo "@dd `bToH` is: ", bToH

proc fanciesDemo() =
  const choice = 19
  case choice:
    of 1 .. 3: echo "Choice is between 1 and 3"
    of 4: echo "Choice is 4"
    of 6, 9, 11: echo "Choice is special"
    of 12 .. 15, 20-1: echo "Choice is even more special"
    else: echo "Choice is unknown"

  block blockDemo:
    while true:
      echo "Attempting infinite loop"
      while true:
        echo "Attempting infinite inception loop"
        break blockDemo
  echo "breaking the block works"

iterator forever(): int =
  var i: int = 0
  while true:
    yield i
    i.inc

iterator foreverOne(): int =
  while true:
    yield 1

proc iteratorsDemo() =
  for i in forever():
    echo i, " to forever"
    if i > 3:
      break

  for i in foreverOne():
    echo i, " forever"
    break

constCompileTimeDemo()
caseAndStyleInsensitiveDemo()
proceduresDemo()
containersDemo()
fanciesDemo()
iteratorsDemo()
