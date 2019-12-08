/* Problem 03-A
 *
 * The gravity assist was successful, and you're well on your way to the Venus refuelling station.
 * During the rush back on Earth, the fuel management system wasn't completely installed, so that's
 * next on the priority list.
 *
 * Opening the front panel reveals a jumble of wires. Specifically, two wires are connected to a
 * central port and extend outward on a grid. You trace the path each wire takes as it leaves the
 * central port, one wire per line of text (your puzzle input).
 *
 * The wires twist and turn, but the two wires occasionally cross paths. To fix the circuit, you
 * need to find the intersection point closest to the central port. Because the wires are on a
 * grid, use the Manhattan distance for this measurement. While the wires do technically cross
 * right at the central port where they both start, this point does not count, nor does a wire
 * count as crossing with itself.
 *
 * For example, if the first wire's path is R8,U5,L5,D3, then starting from the central port (o),
 * it goes right 8, up 5, left 5, and finally down 3:
 *
 * ...........
 * ...........
 * ...........
 * ....+----+.
 * ....|....|.
 * ....|....|.
 * ....|....|.
 * .........|.
 * .o-------+.
 * ...........
 *
 * Then, if the second wire's path is U7,R6,D4,L4, it goes up 7, right 6, down 4, and left 4:
 *
 * ...........
 * .+-----+...
 * .|.....|...
 * .|..+--X-+.
 * .|..|..|.|.
 * .|.-X--+.|.
 * .|..|....|.
 * .|.......|.
 * .o-------+.
 * ...........
 *
 * These wires cross at two locations (marked X), but the lower-left one is closer to the
 * central port: its distance is 3 + 3 = 6.
 *
 * Here are a few more examples:
 *
 * R75,D30,R83,U83,L12,D49,R71,U7,L72
 * U62,R66,U55,R34,D71,R55,D58,R83 = distance 159
 * R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51
 * U98,R91,D20,R16,D67,R40,U7,R15,U6,R7 = distance 135
 *
 * What is the Manhattan distance from the central port to the closest intersection?
 *
 * Answer: 227, point (227, 0)
 */

//TODO: Better approach?  https://www.hackerearth.com/practice/math/geometry/line-intersection-using-bentley-ottmann-algorithm/tutorial/

package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type route struct {
	direction string
	distance  int
}

type point struct {
	x int
	y int
}

// getWireRoutes reads the vector-based routes that define each wire from the comma delimited
// text data file. Each line in the file represents one route.  This function will handle N
// number of routes, but we only expect two.
func getWireRoutes(fileName string) ([][]route, error) {

	routes := []route{}
	wireRoutes := [][]route{}

	// Open data file containing a program
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("cannot load from %s: %v", fileName, err)
	}
	defer file.Close()

	// Read vector-based routes from text file
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		data := strings.Split(scanner.Text(), ",")
		for _, vector := range data {
			if len(vector) > 1 {
				distInt, err := strconv.Atoi(vector[1:])
				if err != nil {
					return nil, fmt.Errorf("could not convert '%s' to an integer", vector[1:])
				}
				routes = append(routes, route{
					direction: vector[:1],
					distance:  distInt,
				})
			}
		}
		wireRoutes = append(wireRoutes, routes)
		routes = nil
	}

	return wireRoutes, nil
}

// traceWires returns a slice containings map of all the points for each wire.  Each wire
// is one element of the slice and the slice element is a map composed of the points as keys
// and the value being the length of the path to get to that point
func traceWires(wireRoutes [][]route) []map[point]int {
	wires := []map[point]int{}
	for _, wireRoute := range wireRoutes {
		location := point{}
		totalSteps := 0
		stepsToLocation := map[point]int{}
		for _, path := range wireRoute {
			for step := 0; step < path.distance; step++ {
				switch path.direction {
				case "R":
					location.x++
				case "D":
					location.y--
				case "L":
					location.x--
				case "U":
					location.y++
				default:
					log.Fatalf("invalid direction of '%s' encountered", path.direction)
				}
				totalSteps++
				stepsToLocation[location] = totalSteps
			}
		}
		wires = append(wires, stepsToLocation)
	}
	return wires
}

// getIntersection utilizes the two maps of the points traced by each wire and returns
// the set of points that both maps have in common
func getIntersections(wires []map[point]int) []point {
	intersections := []point{}
	for key := range wires[0] {
		_, ok := wires[1][key]
		if ok {
			intersections = append(intersections, key)
		}
	}
	return intersections
}

func problem03A(fileName string) int {

	wireRoutes, err := getWireRoutes(fileName)
	if err != nil {
		log.Fatal(err)
	}

	intersections := getIntersections(traceWires(wireRoutes))

	manhattanDistance := math.MaxInt64
	for _, val := range intersections {
		newDistance := int(math.Abs(float64(val.x)) + math.Abs(float64(val.y)))
		if newDistance < manhattanDistance {
			manhattanDistance = newDistance
		}
	}

	return manhattanDistance
}
