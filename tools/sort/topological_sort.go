package sort

import (
	"errors"
	"fmt"
)

// NodeState represents the state of a node during DFS traversal
type NodeState int

const (
	UNVISITED NodeState = 0 // unvisited
	VISITING  NodeState = 1 // visiting (in the DFS stack)
	VISITED   NodeState = 2 // visited
)

// CycleError means that a cycle was detected in the graph during topological sorting
var CycleError = errors.New("graph contains a cycle")

// TopologicalSortReachable sorts the nodes reachable from startNodeID in topological order.
// allNodes: a map of all nodes (nodeID -> node)
// graph: an adjacency list representation of the graph (nodeID -> []neighborIDs)
// startNodeID: the ID of the starting node for traversal
// Returns: a slice of nodes in topological order reachable from startNodeID; or an error if a cycle is detected
func TopologicalSortReachable[T any](allNodes map[string]T, graph map[string][]string, startNodeID string) ([]T, error) {
	if _, ok := graph[startNodeID]; !ok {
		return nil, fmt.Errorf("start node '%s' not found in the graph", startNodeID)
	}
	if _, ok := allNodes[startNodeID]; !ok {
		return nil, fmt.Errorf("start node '%s' not found in the allNodes", startNodeID)
	}

	visitedStatus := make(map[string]NodeState)
	sortedNodes := make([]T, 0)

	// Init all nodes to UNVISITED
	for node := range allNodes {
		visitedStatus[node] = UNVISITED
	}

	// dfsVisit is a recursive function that performs DFS on the graph
	// It returns an error if a cycle is detected
	var dfsVisit func(node string) error
	dfsVisit = func(node string) error {
		// Mark the current node as VISITING
		visitedStatus[node] = VISITING

		// Iterate through all neighbors of the current node
		neighbors := graph[node]
		for _, neighbor := range neighbors {
			neighborStatus, statusExists := visitedStatus[neighbor]
			if !statusExists {
				return fmt.Errorf("The key '%s' is in the graph, but not in allNodes", neighbor)
			}

			switch neighborStatus {
			case VISITING:
				// If the neighbor is already being visited, a cycle is detected
				return fmt.Errorf("%w: detected edge from %s to %s", CycleError, node, neighbor)
			case UNVISITED:
				// If the neighbor has not been visited yet, recursively visit it
				if err := dfsVisit(neighbor); err != nil {
					return err
				}
			case VISITED:
				continue
			}
		}

		// All neighbors of the current node have been visited
		// Mark the current node as VISITED
		visitedStatus[node] = VISITED
		// Append the current node to the end of the result list (to be reversed later)
		sortedNodes = append(sortedNodes, allNodes[node])
		return nil
	}

	// Start DFS from the specified startNodeID
	if err := dfsVisit(startNodeID); err != nil {
		return nil, err
	}

	// Reverse sortedIDs to get the correct topological order
	for i, j := 0, len(sortedNodes)-1; i < j; i, j = i+1, j-1 {
		sortedNodes[i], sortedNodes[j] = sortedNodes[j], sortedNodes[i]
	}

	return sortedNodes, nil
}

// TopologicalSortAll sorts the nodes in topological order (DFS implementation).
// allNodes: a map of all nodes (nodeID -> node data of type T)
// graph: an adjacency list representation of the graph (nodeID -> []neighborIDs)
// Returns: a slice of nodes (of type T) in topological order; or an error if a cycle is detected
func TopologicalSortAll[T any](allNodes map[string]T, graph map[string][]string) ([]T, error) {
	// visited map tracks the state of each node (UNVISITED, VISITING, VISITED)
	visited := make(map[string]NodeState)
	// result slice will store the nodes in topological order (built in reverse)
	result := make([]T, 0, len(allNodes))

	// Define the recursive DFS function using a closure to capture visited, graph, allNodes, and result
	var dfs func(nodeID string) error
	dfs = func(nodeID string) error {
		// Check the current state of the node
		state := visited[nodeID]
		switch state {
		case VISITING:
			// If we encounter a node marked as VISITING, we have found a back edge, indicating a cycle.
			return CycleError
		case VISITED:
			// If the node is already VISITED, we have already processed this node and its descendants.
			return nil // No error, just skip
		case UNVISITED:
			// Mark the node as VISITING (currently in the recursion stack)
			visited[nodeID] = VISITING

			// Recursively visit all neighbors of the current node
			neighbors, exists := graph[nodeID]
			if exists { // Check if the node has outgoing edges listed in the graph
				for _, neighborID := range neighbors {
					// Important: Ensure neighbor exists in allNodes before recursing.
					// If the graph might contain edges to nodes not in allNodes, handle appropriately.
					// Assuming graph only refers to nodes present in allNodes based on typical usage.
					if _, nodeExists := allNodes[neighborID]; !nodeExists {
						return fmt.Errorf("neighbor node '%s' not found in allNodes", neighborID)
					}

					if err := dfs(neighborID); err != nil {
						return err
					}
				}
			}

			// After visiting all neighbors (and their subtrees), mark the node as VISITED.
			visited[nodeID] = VISITED
			// Add the node's data to the *end* of the result slice.
			// Since we add nodes *after* all their dependencies are processed,
			// this builds the topological sort in reverse order.
			result = append(result, allNodes[nodeID])
			return nil // No error for this path
		}
		return fmt.Errorf("Expected error. This should not happen.")
	}

	// Iterate through all nodes provided in allNodes map.
	// This ensures we handle disconnected components in the graph.
	for nodeID := range allNodes {
		if visited[nodeID] == UNVISITED {
			if err := dfs(nodeID); err != nil {
				return nil, err
			}
		}
	}

	// Reverse the slice to get the correct topological order.
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result, nil
}
