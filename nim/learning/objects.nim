
echo "------------------- runtime ---------------------------"

type
    StackAllocatedOnly = object
        memberStr: string
        memberInt: int
        justFloat: float

let test = StackAllocatedOnly(memberStr: "test", memberInt: 1, justFloat: 2.1)

echo "@dd ", test
