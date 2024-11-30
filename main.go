package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func GetData(dataFile string) (start, end string, rooms []string, links []string, antNumbers int) {
	data, err := os.ReadFile(dataFile)
	if err != nil {
		log.Fatal(err)
	}

	if len(data) == 0 {
		log.Fatal("ERROR: No data in this file!!")
	}

	checkData := strings.Split(string(data), "\n")
	canStart, canEnd := false, false
	for _, flag := range checkData {
		if flag == "##start" {
			canStart = true
		}
		if flag == "##end" {
			canEnd = true
		}
	}
	if !canStart || !canEnd {
		log.Fatal("ERROR: Check if the ##start command and ##end command exist")
	}
	affStart, affEnd := false, false
	findStart, findEnd := false, false

	for i := len(checkData) - 1; i >= 0; i-- {

		tmp := strings.TrimSpace(checkData[i])

		if tmp == "" {
			checkData = checkData[:i]
		} else {
			break
		}

	}

	for i := 0; i < len(checkData)-1; i++ {

		tmp := strings.TrimSpace(checkData[i])

		if tmp == "" {
			checkData = checkData[i+1:]

			i--
		} else {
			break
		}

	}

	for i, line := range checkData {

		if len(line) == 0 && i != 0 {
			log.Fatal("ERROR: invalid data format, empty line data!")
		}

		if strings.HasPrefix(line, "#") {
			if line == "##start" {
				if findStart {
					log.Fatal("ERROR: Issue with the start command (##start)")
				}
				affStart = true
				findStart = true
			} else if line == "##end" && len(start) != 0 {
				if findEnd {
					log.Fatal("ERROR: Issue with the end command (##end)")
				}
				affEnd = true
				findEnd = true
			} else if line == "##end" && len(start) == 0 {
				log.Fatal("ERROR: Can't found the start room")
			} else if len(line) > 1 && strings.HasPrefix(line, "##") && line != "##start" && line != "##end" {
				log.Fatal("ERROR: Only ##start and ##end are allowed commands")
			}
			continue
		}

		findErr := strings.Fields(line)

		if antNumbers != 0 {
			if !strings.Contains(line, "-") && len(findErr) != 3 {
				log.Fatal("ERROR: Invalid data!")
			}
		}
		if affStart {
			if len(strings.Fields(line)) == 3 {
				start = strings.Fields(line)[0]
				affStart = false
			} else {
				log.Fatal("ERROR: Can't found the start room")
			}
		}
		// If this is the first non-comment line after ##end, capture the end room
		if affEnd {
			if len(strings.Fields(line)) == 3 {
				end = strings.Fields(line)[0]
				affEnd = false
			} else {
				log.Fatal("ERROR: Can't found the end room")
			}
		}

		if antNumbers == 0 {
			antNumbers, err = strconv.Atoi(line)
			if err != nil {
				log.Fatal("Error converting number of ants to int:", err)
			}
			if antNumbers <= 0 {
				log.Fatal("ERROR: invalid number of ants")
			}
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 3 {
			rooms = append(rooms, parts[0])
		}
		if strings.Contains(line, "-") {
			x := strings.Split(line, "-")
			if len(x) != 2 || x[0] == "" || x[1] == "" || x[0] == x[1] {
				log.Fatal("ERROR: Invalid data format")
			}

			validRoomName := func(name string) bool {
				for _, room := range rooms {
					if room == name {
						return true
					}
				}

				return false
			}

			if !validRoomName(x[0]) || !validRoomName(x[1]) {
				log.Fatal("ERROR: Issue with room names")
			}

			links = append(links, line)
		}
	}

	for i, v := range rooms {
		for j, vv := range rooms {
			if i != j && vv == v {
				log.Fatal("ERROR: Room duplicated room !!!")
			}
		}
	}

	for i, v := range links {
		for j, vv := range links {
			if i != j && vv == v {
				log.Fatal("ERROR: Room duplicated links !!!")
			}
		}
	}
	checkStartRoom, checkEndRoom := false, false
	for _, v := range links {
		check := strings.Split(v, "-")
		for _, room := range check {
			if room == start {
				checkStartRoom = true
			}
			if room == end {
				checkEndRoom = true
			}
		}

	}
	if !checkStartRoom || !checkEndRoom {
		log.Fatal("ERROR: Start/End Room not linked !")
	}

	for _, line := range checkData {
		fmt.Println(line)
	}

	fmt.Println()
	return start, end, rooms, links, antNumbers
}

func GraphMaker(Rooms, Links []string) map[string][]string {
	Graph := make(map[string][]string, len(Rooms))
	for _, room := range Rooms {
		Graph[room] = []string{}
	}
	for _, link := range Links {
		rooms := strings.Split(link, "-")
		if len(rooms) == 2 {
			room1, room2 := rooms[0], rooms[1]
			Graph[room1] = append(Graph[room1], room2)
			Graph[room2] = append(Graph[room2], room1)
		}
	}

	return Graph
}

// FindAllPaths finds all paths from start to end
func FindAllPaths(graph map[string][]string, start, end string) [][]string {
	var paths [][]string
	var currentPath []string
	visited := make(map[string]bool)

	var dfs func(string)
	dfs = func(node string) {
		if visited[node] {
			return
		}
		// Add the current node to the path
		currentPath = append(currentPath, node)
		visited[node] = true

		if node == end {
			// Found a valid path, append it
			pathCopy := make([]string, len(currentPath))
			copy(pathCopy, currentPath)
			paths = append(paths, pathCopy)
		} else {
			// Explore neighbors
			for _, neighbor := range graph[node] {
				dfs(neighbor)
			}
		}

		// Backtrack
		visited[node] = false
		currentPath = currentPath[:len(currentPath)-1]
	}

	dfs(start)

	// Sort paths by length, then lexicographically as a tiebreaker
	sort.Slice(paths, func(i, j int) bool {
		if len(paths[i]) != len(paths[j]) {
			return len(paths[i]) < len(paths[j]) // Shorter paths first
		}
		// If lengths are equal, sort lexicographically
		for x := 0; x < len(paths[i]) && x < len(paths[j]); x++ {
			if paths[i][x] != paths[j][x] {
				return paths[i][x] < paths[j][x]
			}
		}
		return false
	})

	return paths
}

func FilterPaths(paths [][]string, start, end string) [][]string {
	var filteredPaths [][]string
	visitedRooms := make(map[string]bool) // Track visited intermediate rooms

	for _, path := range paths {
		hasOverlap := false

		// Check intermediate rooms for overlap (exclude start and end)
		for i := 1; i < len(path)-1; i++ {
			room := path[i]
			if visitedRooms[room] {
				hasOverlap = true
				break
			}
		}

		// If no overlap, accept the path and mark its intermediate rooms as visited
		if !hasOverlap {
			filteredPaths = append(filteredPaths, path)
			for i := 1; i < len(path)-1; i++ {
				visitedRooms[path[i]] = true
			}
		}
	}

	return filteredPaths
}

// Path struct holds path details
type Path struct {
	Nodes    []string
	Capacity int // The max number of ants that can be sent through this path
}

// AssignAntsToPaths decides whether to use one path or both based on the number of ants
func AssignAntsToPaths(paths [][]string, numAnts int) [][]string {
	// Define the paths and their capacities
	// Assume each path has a capacity based on the number of rooms, you can adjust the logic here
	pathDetails := []Path{
		{Nodes: paths[0], Capacity: 5}, // Path 1, capacity 5 (arbitrary)
		{Nodes: paths[1], Capacity: 6}, // Path 2, capacity 6 (arbitrary)
	}

	var finalPaths [][]string

	// Step 1: Decide based on the number of ants if we use one or both paths
	if numAnts <= pathDetails[0].Capacity {
		// If ants are less than or equal to the capacity of the first path, use only the first path
		finalPaths = append(finalPaths, pathDetails[0].Nodes)
	} else {
		// If ants exceed the capacity of the first path, use both paths
		finalPaths = append(finalPaths, pathDetails[0].Nodes) // Use first path
		finalPaths = append(finalPaths, pathDetails[1].Nodes) // Use second path
	}

	return finalPaths
}


// PrintAntMovements function prints the movement of ants based on their path assignments
func PrintAntMovements(antAssignments [][]int) string {
	// Initialize the result string to store each step's movements
	var result []string

	// Get the maximum length of the paths to determine how many steps we have
	maxSteps := 0
	for _, path := range antAssignments {
		if len(path) > maxSteps {
			maxSteps = len(path)
		}
	}

	// Simulate the movements for each step
	for step := 0; step < maxSteps; step++ {
		// Collect the movements for this step
		var stepMovements []string

		// Loop over each ant's path and get its current position at this step
		for antID, path := range antAssignments {
			if step < len(path) {
				// Append the movement in the required format: L<ant_number>-<room_number>
				stepMovements = append(stepMovements, fmt.Sprintf("L%d-%d", antID+1, path[step]))
			}
		}

		// If there are movements at this step, add them to the result
		if len(stepMovements) > 0 {
			result = append(result, strings.Join(stepMovements, " "))
		}
	}

	// Return the final result as a formatted string with line breaks
	return strings.Join(result, "\n")
}


func main() {
	if len(os.Args) != 2 {
		log.Fatal("\nInvalid Arguments, \nUsage: go run . [file_name.txt]")
	}

	// Start, End, Rooms, Links, antNumbers := GetData(os.Args[1])
	Start, End, Rooms, Links, antNumbers := GetData(os.Args[1])

	Graph := GraphMaker(Rooms, Links)
	fmt.Println(Graph)

	Paths := FindAllPaths(Graph, Start, End)
	fmt.Println("All paths: ")
	for _, p := range Paths {
		fmt.Println(p)
	}
	fmt.Println("Length of paths: ", len(Paths))

	FilteredPaths := FilterPaths(Paths, Start, End)
	fmt.Println("Paths after filtering: ")
	for _, pf := range FilteredPaths {
		fmt.Println(pf)
	}
	fmt.Println("Length of paths after feltering: ", len(FilteredPaths))

	finalPath := AssignAntsToPaths(FilteredPaths, antNumbers)
	fmt.Println("Final Path: ", finalPath)



	PrintAntMovements(finalPath, antNumbers)
}
