package abidec

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
)

func ParseEventSignature(sig string) (abi.Event, error) {
	// Remove "event " prefix if present
	if strings.HasPrefix(sig, "event ") {
		sig = strings.TrimPrefix(sig, "event ")
	}

	// Create a dummy JSON ABI with just this event
	jsonABI := fmt.Sprintf(`[{"type":"event","name":"%s","inputs":%s}]`,
		getEventName(sig),
		getEventArgs(sig))

	// Parse the JSON ABI
	parsedABI, err := abi.JSON(strings.NewReader(jsonABI))
	if err != nil {
		return abi.Event{}, fmt.Errorf("failed to parse event ABI: %w", err)
	}

	// Get the event from the parsed ABI
	event, ok := parsedABI.Events[getEventName(sig)]
	if !ok {
		return abi.Event{}, fmt.Errorf("event not found in parsed ABI")
	}

	return event, nil
}

// getEventName extracts the event name from the signature
// e.g., "Transfer(address,address,uint256)" -> "Transfer"
func getEventName(sig string) string {
	if idx := strings.Index(sig, "("); idx != -1 {
		return sig[:idx]
	}
	return sig
}

// getEventArgs extracts and formats the event arguments as JSON array
// e.g., "Transfer(address indexed from, address indexed to, uint256 value)" ->
// [{"type":"address","name":"from","indexed":true},{"type":"address","name":"to","indexed":true},{"type":"uint256","name":"value","indexed":false}]
func getEventArgs(sig string) string {
	start := strings.Index(sig, "(")
	end := strings.LastIndex(sig, ")")
	if start == -1 || end == -1 {
		return "[]"
	}

	args := sig[start+1 : end]
	if args == "" {
		return "[]"
	}

	// Split the arguments
	parts := strings.Split(args, ",")
	jsonArgs := make([]string, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// Parse each argument
		words := strings.Fields(part)
		if len(words) < 2 {
			continue
		}

		isIndexed := false
		typeIdx := 0
		nameIdx := len(words) - 1

		// Check for 'indexed' keyword
		for i, word := range words {
			if word == "indexed" {
				isIndexed = true
				typeIdx = i - 1
				break
			}
		}
		if !isIndexed {
			typeIdx = len(words) - 2
		}

		// Create JSON object for the argument
		argJSON := fmt.Sprintf(`{"type":"%s","name":"%s","indexed":%t}`,
			words[typeIdx],
			words[nameIdx],
			isIndexed)
		jsonArgs = append(jsonArgs, argJSON)
	}

	return "[" + strings.Join(jsonArgs, ",") + "]"
}

func ParseLogIntoMap(event abi.Event, out map[string]any, log *types.Log) error {
	if event.ID.Cmp(log.Topics[0]) != 0 {
		return fmt.Errorf("event ID mismatch")
	}
	if len(log.Data) > 0 {
		if err := event.Inputs.UnpackIntoMap(out, log.Data); err != nil {
			return fmt.Errorf("failed to unpack data: %w", err)
		}
	}

	var indexed abi.Arguments
	for _, arg := range event.Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}

	if err := abi.ParseTopicsIntoMap(out, indexed, log.Topics[1:]); err != nil {
		return fmt.Errorf("failed to parse topics: %w", err)
	}

	return nil
}
