package sort

import (
	"errors"
	"fmt"
	"reflect" // Needed for deep slice comparison
	"testing"
)

func TestTopologicalSortReachable(t *testing.T) {
	// Define test cases using a table
	testCases := []struct {
		name        string              // Name for the test case
		graph       map[string][]string // Input graph
		startNodeID string              // Input start node
		wantResult  []string            // Expected successful result (nil if error expected)
		wantErr     error               // Expected error (nil if success expected) - use errors.Is for matching
	}{
		{
			name: "Basic DAG from B",
			graph: map[string][]string{
				"A": {"C"},
				"B": {"C", "D"},
				"C": {"E"},
				"D": {"F"},
				"E": {"F"},
				"F": {},
				"G": {"A", "B"},
			},
			startNodeID: "B",
			// Note: Topological sort can have multiple valid orders.
			// We expect one valid order here. A more robust test might
			// verify the dependency constraints instead of exact order.
			// Let's assume the implementation might produce this order:
			wantResult: []string{"B", "D", "C", "E", "F"},
			// Alternative valid result: []string{"B", "C", "D", "E", "F"}
			// If your implementation consistently gives the alternative, update wantResult.
			wantErr: nil,
		},
		{
			name: "Disconnected component from 1",
			graph: map[string][]string{
				"1": {"3"},
				"2": {"3"},
				"3": {"4"},
				"4": {},
				"5": {"6"},
				"6": {},
			},
			startNodeID: "1",
			wantResult:  []string{"1", "3", "4"},
			wantErr:     nil,
		},
		{
			name: "Graph with cycle",
			graph: map[string][]string{
				"X": {"Y"},
				"Y": {"Z"},
				"Z": {"X"}, // Cycle Z -> X
				"W": {"X"},
			},
			startNodeID: "W",
			wantResult:  nil,        // Expecting an error, so result should be nil
			wantErr:     CycleError, // Expect the specific CycleError type
		},
		{
			name: "Start node not found",
			graph: map[string][]string{
				"A": {"B"},
			},
			startNodeID: "C",
			wantResult:  nil,
			wantErr:     fmt.Errorf("start node 'C' not found in the graph"), // Expect specific error message/type
		},
		{
			name: "Start node exists but no outgoing edges",
			graph: map[string][]string{
				"A": {},
				"B": {"A"},
			},
			startNodeID: "A",
			wantResult:  []string{"A"},
			wantErr:     nil,
		},
		{
			name: "Start node is isolated",
			graph: map[string][]string{
				"X": {"Y"},
				"Z": {}, // Isolated node
			},
			startNodeID: "Z",
			wantResult:  []string{"Z"},
			wantErr:     nil,
		},
		{
			name:        "Empty graph",
			graph:       map[string][]string{},
			startNodeID: "A",
			wantResult:  nil,
			wantErr:     fmt.Errorf("start node 'A' not found in the graph"),
		},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		// Run each case as a subtest
		t.Run(tc.name, func(t *testing.T) {
			allNodes := make(map[string]string) // 使用空结构体节省空间
			for node, neighbors := range tc.graph {
				allNodes[node] = node
				for _, neighbor := range neighbors {
					allNodes[neighbor] = neighbor
				}
			}
			gotResult, gotErr := TopologicalSortReachable(allNodes, tc.graph, tc.startNodeID)

			// --- Assertions ---

			// 1. Check for unexpected errors
			if tc.wantErr == nil && gotErr != nil {
				t.Fatalf("TopologicalSortReachable() returned an unexpected error: %v", gotErr)
			}

			// 2. Check if expected error occurred and has the right type/value
			if tc.wantErr != nil {
				if gotErr == nil {
					t.Fatalf("TopologicalSortReachable() expected error '%v', but got nil", tc.wantErr)
				}
				// Use errors.Is for checking wrapped errors like CycleError
				// Use direct comparison for specific fmt.Errorf instances (if messages must match exactly)
				if errors.Is(tc.wantErr, CycleError) { // Check if it's CycleError or wraps it
					if !errors.Is(gotErr, CycleError) {
						t.Errorf("TopologicalSortReachable() expected error type '%T', but got type '%T' (%v)", CycleError, gotErr, gotErr)
					}
				} else if tc.wantErr.Error() != gotErr.Error() { // For other specific errors, compare messages
					t.Errorf("TopologicalSortReachable() expected error '%v', but got '%v'", tc.wantErr, gotErr)
				}
				// Also check if result is nil on error
				if gotResult != nil {
					t.Errorf("TopologicalSortReachable() expected nil result on error, but got %v", gotResult)
				}
			}

			// 3. Check the result slice if no error was expected
			if tc.wantErr == nil {
				// Use reflect.DeepEqual for comparing slice contents
				if !reflect.DeepEqual(gotResult, tc.wantResult) {
					// Handle the potential ambiguity for the first test case
					if tc.name == "Basic DAG from B" {
						alternativeWant := []string{"B", "C", "D", "E", "F"}
						if !reflect.DeepEqual(gotResult, alternativeWant) {
							t.Errorf("TopologicalSortReachable() = %v, want %v (or %v)", gotResult, tc.wantResult, alternativeWant)
						}
					} else {
						t.Errorf("TopologicalSortReachable() = %v, want %v", gotResult, tc.wantResult)
					}
				}
			}
		})
	}
}

// Helper type for testing
type TestNode struct {
	ID string
}

// TestTopologicalSortAll tests the TopologicalSortAll function
func TestTopologicalSortAll(t *testing.T) {
	// Helper to get ID from TestNode
	getNodeID := func(n TestNode) string { return n.ID }

	// Define test cases
	testCases := []struct {
		name        string
		allNodes    map[string]TestNode
		graph       map[string][]string
		expectError error
		// We don't specify expectedResult slice because multiple valid sorts can exist.
		// Instead, we'll use isValidTopologicalSort.
	}{
		{
			name:        "Empty Graph",
			allNodes:    map[string]TestNode{},
			graph:       map[string][]string{},
			expectError: nil,
		},
		{
			name: "Single Node",
			allNodes: map[string]TestNode{
				"A": {ID: "A"},
			},
			graph: map[string][]string{
				"A": {},
			},
			expectError: nil,
		},
		{
			name: "Simple Linear Chain", // A -> B -> C
			allNodes: map[string]TestNode{
				"A": {ID: "A"},
				"B": {ID: "B"},
				"C": {ID: "C"},
			},
			graph: map[string][]string{
				"A": {"B"},
				"B": {"C"},
				"C": {},
			},
			expectError: nil,
		},
		{
			name: "Multiple Independent Nodes", // A, B
			allNodes: map[string]TestNode{
				"A": {ID: "A"},
				"B": {ID: "B"},
			},
			graph: map[string][]string{
				"A": {},
				"B": {},
			},
			expectError: nil,
		},
		{
			name: "DAG with Multiple Paths", // A->B, A->C, B->D, C->D
			allNodes: map[string]TestNode{
				"A": {ID: "A"},
				"B": {ID: "B"},
				"C": {ID: "C"},
				"D": {ID: "D"},
			},
			graph: map[string][]string{
				"A": {"B", "C"},
				"B": {"D"},
				"C": {"D"},
				"D": {},
			},
			expectError: nil,
		},
		{
			name: "Disconnected Components", // A->B, C->D
			allNodes: map[string]TestNode{
				"A": {ID: "A"},
				"B": {ID: "B"},
				"C": {ID: "C"},
				"D": {ID: "D"},
			},
			graph: map[string][]string{
				"A": {"B"},
				"B": {},
				"C": {"D"},
				"D": {},
			},
			expectError: nil,
		},
		{
			name: "More Complex DAG", // 5->{11,2}, 7->{11,8}, 3->{8,10}, 11->{2,9,10}, 8->9
			allNodes: map[string]TestNode{
				"2":  {ID: "2"},
				"3":  {ID: "3"},
				"5":  {ID: "5"},
				"7":  {ID: "7"},
				"8":  {ID: "8"},
				"9":  {ID: "9"},
				"10": {ID: "10"},
				"11": {ID: "11"},
			},
			graph: map[string][]string{
				"5":  {"11", "2"},
				"7":  {"11", "8"},
				"3":  {"8", "10"},
				"11": {"2", "9", "10"},
				"8":  {"9"},
				"2":  {},
				"9":  {},
				"10": {},
			},
			expectError: nil,
			// Possible sorts: [3, 5, 7, 8, 11, 2, 9, 10], [5, 7, 3, 11, 8, 10, 9, 2], etc.
		},
		{
			name: "Simple Cycle", // A -> B, B -> A
			allNodes: map[string]TestNode{
				"A": {ID: "A"},
				"B": {ID: "B"},
			},
			graph: map[string][]string{
				"A": {"B"},
				"B": {"A"},
			},
			expectError: CycleError,
		},
		{
			name: "Longer Cycle", // A -> B -> C -> A
			allNodes: map[string]TestNode{
				"A": {ID: "A"},
				"B": {ID: "B"},
				"C": {ID: "C"},
			},
			graph: map[string][]string{
				"A": {"B"},
				"B": {"C"},
				"C": {"A"},
			},
			expectError: CycleError,
		},
		{
			name: "Self Loop", // A -> A
			allNodes: map[string]TestNode{
				"A": {ID: "A"},
			},
			graph: map[string][]string{
				"A": {"A"},
			},
			expectError: CycleError,
		},
		{
			name: "Cycle within Larger Graph", // A->B, B->C, C->D, D->B (cycle BCD), E->A
			allNodes: map[string]TestNode{
				"A": {ID: "A"},
				"B": {ID: "B"},
				"C": {ID: "C"},
				"D": {ID: "D"},
				"E": {ID: "E"},
			},
			graph: map[string][]string{
				"A": {"B"},
				"B": {"C"},
				"C": {"D"},
				"D": {"B"}, // Cycle back to B
				"E": {"A"},
			},
			expectError: CycleError,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// --- Run the function under test ---
			result, err := TopologicalSortAll(tc.allNodes, tc.graph)
			// --- ---

			// Check error expectation
			if tc.expectError != nil {
				if err == nil {
					t.Errorf("Expected error '%v', but got nil", tc.expectError)
				} else if !errors.Is(err, tc.expectError) { // Use errors.Is for potential wrapping
					t.Errorf("Expected error '%v', but got '%v'", tc.expectError, err)
				}
				// If error was expected, don't check the result slice (should be nil)
				if result != nil {
					t.Errorf("Expected nil result on error, but got %v", result)
				}
				return // End test case for expected error
			}

			// If no error was expected, check if one occurred
			if err != nil {
				t.Fatalf("Expected no error, but got: %v", err)
			}

			// Validate the topological sort result
			if ok, reason := isValidTopologicalSort(tc.allNodes, tc.graph, result, getNodeID); !ok {
				// The helper function prints details, just mark failure here
				// Convert result to IDs for easier debugging output
				resultIDs := make([]string, len(result))
				for i, n := range result {
					resultIDs[i] = getNodeID(n)
				}
				t.Errorf("Result %v is not a valid topological sort for the given graph: %s", resultIDs, reason)
			} else {
				// Optional: Log the valid sort found
				resultIDs := make([]string, len(result))
				for i, n := range result {
					resultIDs[i] = getNodeID(n)
				}
				t.Logf("Valid topological sort found: %v", resultIDs)
			}
		})
	}
}

func isValidTopologicalSort[T any](
	allNodes map[string]T,
	graph map[string][]string,
	result []T,
	idGetter func(T) string,
) (bool, string) {

	// --- Check 1: Length Mismatch ---
	// The result must contain the same number of nodes as expected.
	if len(result) != len(allNodes) {
		return false, fmt.Sprintf("length mismatch: expected %d nodes, but result has %d", len(allNodes), len(result))
	}

	// Handle the trivial case of an empty graph explicitly for clarity
	if len(allNodes) == 0 {
		if len(result) == 0 {
			return true, "" // Empty graph has an empty valid sort
		}
		// Already covered by length check above, but reinforces logic
		// return false, "length mismatch: expected 0 nodes for empty graph, but result is not empty"
	}

	// --- Check 2: Node Set Verification (Completeness & Uniqueness) ---
	// Build a map to track nodes seen in the result and their positions.
	nodePositions := make(map[string]int, len(result))
	seenIDs := make(map[string]bool, len(result))
	resultIDs := make([]string, len(result)) // Keep IDs for potential error messages

	for i, node := range result {
		nodeID := idGetter(node)
		resultIDs[i] = nodeID // Store for error messages

		// 2a: Check if the node from the result is one of the expected nodes.
		if _, expected := allNodes[nodeID]; !expected {
			return false, fmt.Sprintf("unexpected node: result contains node '%s' which was not in the expected set of nodes", nodeID)
		}

		// 2b: Check for duplicate nodes in the result.
		if seenIDs[nodeID] {
			return false, fmt.Sprintf("duplicate node: node '%s' appears more than once in the result", nodeID)
		}

		// Mark node as seen and record its position.
		seenIDs[nodeID] = true
		nodePositions[nodeID] = i
	}

	// 2c: Verify all expected nodes were present in the result.
	// This check is technically redundant if the length check passed and no duplicates were found,
	// but it adds clarity and catches potential logic errors elsewhere.
	if len(seenIDs) != len(allNodes) {
		// This case implies len(result) == len(allNodes) but len(seenIDs) < len(allNodes),
		// meaning duplicates *must* have existed, which should have been caught earlier.
		// Or, it implies a node in `result` wasn't in `allNodes`, also caught earlier.
		// It primarily guards against unexpected states or logic flaws.
		return false, fmt.Sprintf("missing nodes: expected %d unique nodes, but result only contains %d unique nodes from the expected set", len(allNodes), len(seenIDs))
	}

	// --- Check 3: Dependency Order Verification ---
	// Iterate through all defined edges in the graph.
	for uID, neighbors := range graph {
		// Get the position of the source node 'u'.
		// We only care about edges where the source node 'uID' was actually expected
		// to be part of the sort (i.e., it's a key in allNodes).
		if _, uIsExpected := allNodes[uID]; !uIsExpected {
			continue // Skip edges originating from nodes not in the defined set.
		}

		uPos, uOk := nodePositions[uID]
		if !uOk {
			// This indicates that an expected node 'uID' (which is part of an edge in the graph)
			// was NOT found in the 'result' slice. This contradicts Check 2c passing.
			// This signifies a logic error either here or in Check 2, or inconsistent input state.
			return false, fmt.Sprintf("internal validation error: expected node '%s' (source of edge) not found in result positions map, though it should be present", uID)
		}

		// Check each dependency (edge u -> v).
		for _, vID := range neighbors {
			// Similarly, only check dependencies where the target node 'vID' was expected.
			if _, vIsExpected := allNodes[vID]; !vIsExpected {
				continue // Skip edges pointing to nodes not in the defined set.
			}

			vPos, vOk := nodePositions[vID]
			if !vOk {
				// An expected node 'vID' (target of an edge) was not found in the result.
				return false, fmt.Sprintf("internal validation error: expected node '%s' (target of edge from '%s') not found in result positions map, though it should be present", vID, uID)
			}

			// The core topological sort check: u must come before v.
			if uPos >= vPos {
				return false, fmt.Sprintf("dependency violation: node '%s' must appear before '%s' (due to edge %s -> %s), but positions are %d and %d in result %v", uID, vID, uID, vID, uPos, vPos, resultIDs)
			}
		}
	}

	// --- Success ---
	// If all checks passed, the result is a valid topological sort.
	return true, ""
}
