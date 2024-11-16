# abidec

ABI Decoder From Solidity To Golang

## Problem

Imagine you have a Solidity function and struct like this:

```js
struct Task {
	string title;
	string description;
	address reporter;
	address assignee;
	uint256[2] deadline;
}

function encode(string calldata title, string calldata description, address reporter, address assignee, uint256 startDate, uint256 endDate) external pure returns(bytes memory) {
    uint256[2] memory deadline = [startDate, endDate];
   
    Task memory task = Task(title, description, reporter, assignee, deadline);
   
    return abi.encode(task);
}
```

And you need to decode ABI encoded bytes to a struct, like that:
```js
function decode(bytes calldata params) external pure returns (Task memory) {
    (Task memory task) = abi.decode(params, (Task));
    return task;
}
```

**But in Golang.**

## Usage

Just paste your struct definition and use it in `abidec.NewAbiDecoder` options.

```go
decoder := abidec.NewAbiDecoder(
    abidec.WithStruct(`
		struct Task {
			string title;
			string description;
			address reporter;
			address assignee
			uint256[2] deadline;
		}
    `),
)
```

And then you can decode ABI encoded bytes to values:

```go
// params = 
// [1]  0000000000000000000000000000000000000000000000000000000000000020
// [2]  00000000000000000000000000000000000000000000000000000000000000c0
// [3]  0000000000000000000000000000000000000000000000000000000000000010
// [4]  000000000000000000000000acdad15d8f07d8df258fe11332b752785d6b1d22
// [5]  0000000000000000000000008b1383709d1e80a291de5d67993252dfc52c3700
// [6]  00000000000000000000000000000000000000000000000000000000672f269b
// [7]  00000000000000000000000000000000000000000000000000000000672f34ab
// [8]  0000000000000000000000000000000000000000000000000000000000000006
// [9]  526f636b657400000000000000000000000000000000000000000000000000000
// [10] 0000000000000000000000000000000000000000000000000000000000019
// [11] 43726561746520526f636b657420546f20546865204d6f6f6e00000000000000
paramsBytes := hex.DecodeString(params)

type Task struct {
	Title       string
	Description string
	Reporter    [20]byte // or common.Address from go-ethereum
	Assignee    [20]byte
	Deadline    [2]*big.Int	
}

var task Task
decoder.DecodeStruct(paramsBytes, &task)
```

Pretty print the values in JSON format
```json
{
    "assignee": "0x8B1383709D1e80A291DE5d67993252dFC52C3700",
    "deadline": [
        1731143323,
        1731146923
    ],
    "description": "Create Rocket To The Moon",
    "reporter": "0xaCDaD15d8F07D8Df258fe11332b752785d6b1d22",
    "title": "Rocket"
}
```

