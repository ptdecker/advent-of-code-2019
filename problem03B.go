/* Problem 03-B
 *
 * It turns out that this circuit is very timing-sensitive; you actually need to minimize the signal
 * delay.
 *
 * To do this, calculate the number of steps each wire takes to reach each intersection; choose the
 * intersection where the sum of both wires' steps is lowest. If a wire visits a position on the grid
 * multiple times, use the steps value from the first time it visits that position when calculating
 * the total value of a specific intersection.
 *
 * The number of steps a wire takes is the total number of grid squares the wire has entered to get
 * to that location, including the intersection being considered. Again consider the example from
 * above:
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
 * In the above example, the intersection closest to the central port is reached after 8+5+5+2 = 20
 * steps by the first wire and 7+6+4+3 = 20 steps by the second wire for a total of 20+20 = 40 steps.
 *
 * However, the top-right intersection is better: the first wire takes only 8+5+2 = 15 and the
 * second wire takes only 7+6+2 = 15, a total of 15+15 = 30 steps.
 *
 * Here are the best steps for the extra examples from above:
 *
 * R75,D30,R83,U83,L12,D49,R71,U7,L72
 * U62,R66,U55,R34,D71,R55,D58,R83 = 610 steps
 * R98,U47,R26,D63,R33,U87,L62,D20,R33,U53,R51
 * U98,R91,D20,R16,D67,R40,U7,R15,U6,R7 = 410 steps
 *
 * What is the fewest combined steps the wires must take to reach an intersection?
 *
 * Answer: 20286
 *
 * !!! NOTE:  Credit for this answer goes totally to u/DukeNewcombe and u/matoxxx from this
 *            thread https://www.reddit.com/r/adventofcode/comments/e653xt/2019_day_3_part_2nodejs_program_gives_correct/
 *            which got me on the right track.  But, the solution is terribly slow.
 */

package main

import (
	"log"
)

func problem03B(fileName string) int {

	wireRoutes, err := getWireRoutes(fileName)
	if err != nil {
		log.Fatal(err)
	}

	history := []map[point]int{}
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
		history = append(history, stepsToLocation)
	}

	intersections := []point{}
	for key := range history[0] {
		_, ok := history[1][key]
		if ok {
			intersections = append(intersections, key)
		}
	}

	minKey := intersections[0]
	minDistance := history[0][minKey] + history[1][minKey]
	for _, val := range intersections {
		totalSteps := history[0][val] + history[1][val]
		if totalSteps < minDistance {
			minDistance = totalSteps
			minKey = val
		}
	}

	return minDistance
}