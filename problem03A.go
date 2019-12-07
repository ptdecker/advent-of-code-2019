/* Problem 03-A
 *
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

type orientation int

type route struct {
	direction string
	distance  int
}

type point struct {
	x int
	y int
}

type segment struct {
	start    point
	end      point
	vertical bool
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

// Converts sets of vector-based routes to a route defined by a series of points
func convertRoutesToPaths(wireRoutes [][]route) ([][]point, error) {
	paths := [][]point{}
	for _, routes := range wireRoutes {
		position := point{x: 0, y: 0}
		path := []point{position}
		for _, route := range routes {
			switch route.direction {
			case "R":
				position.x += route.distance
			case "D":
				position.y -= route.distance
			case "L":
				position.x -= route.distance
			case "U":
				position.y += route.distance
			default:
				return nil, fmt.Errorf("invalid direction of '%s' encountered", route.direction)
			}
			path = append(path, position)
		}
		paths = append(paths, path)
	}
	return paths, nil
}

// getSegments pulls together line segments from a wire path defined as a set of point.
// It also determines if that segment is vertically or horizongally oriented.
func getSegments(wirePaths [][]point) [][]segment {
	wireSegments := [][]segment{}
	for _, paths := range wirePaths {
		segments := []segment{}
		if len(paths) > 1 {
			var start point
			for index, path := range paths {
				if index > 0 {
					segments = append(segments, segment{
						start:    start,
						end:      path,
						vertical: (start.x == path.x),
					})
				}
				start = path
			}
			wireSegments = append(wireSegments, segments)
			segments = nil
		}
	}
	return wireSegments
}

// findIntersection finds the intersection of two line segments if one exists. The two
// segments must be tangent to each other, one horizontal, one vertical and the horizontal segment
// is passed first.
func findIntersection(horzSeg, vertSeg segment) (point, bool) {

	if horzSeg.start.y != horzSeg.end.y {
		log.Fatalf("horizontal segment %v is not horizontal", horzSeg)
	}
	horzAxis := horzSeg.start.y
	yMin := vertSeg.start.y
	yMax := vertSeg.end.y
	if yMin > yMax {
		yMin = vertSeg.end.y
		yMax = vertSeg.start.y
	}
	if yMax < horzAxis || yMin > horzAxis {
		return point{x: 0, y: 0}, false
	}

	if vertSeg.start.x != vertSeg.end.x {
		log.Fatalf("vertical segment %v is not vertical", vertSeg)
	}
	vertAxis := vertSeg.start.x
	xMin := horzSeg.start.x
	xMax := horzSeg.end.x
	if xMin > xMax {
		xMin = horzSeg.end.x
		xMax = horzSeg.start.x
	}
	if xMax < vertAxis || xMin > vertAxis {
		return point{x: 0, y: 0}, false
	}

	return point{x: vertAxis, y: horzAxis}, true
}

//getIntersections returns a set of points that represent all the intersections between two
//sets of line segments.  It does not return intersections between vertical or horizontal
//segments, i.e. it only returns intersections if the segments have opposite orientations.
func getIntersections(set1, set2 []segment) []point {
	points := []point{}
	for _, set1seg := range set1 {
		for _, set2seg := range set2 {
			if set1seg.vertical == set2seg.vertical { // ignore segments with same orientation
				continue
			}
			var intersection point
			var ok bool
			if set1seg.vertical {
				intersection, ok = findIntersection(set2seg, set1seg)
			} else {
				intersection, ok = findIntersection(set1seg, set2seg)
			}
			if ok {
				// fmt.Printf("%v might intersect %v at %v\n", set1seg, set2seg, intersection)
				points = append(points, intersection)
			} else {
				// fmt.Printf("%v does not intersect %v\n", set1seg, set2seg)
			}
		}
	}
	return points
}

// calcMinManhattanDistance calculates the Manhattan Distances for each point in a set between
// that point and the origin (0,0).  It returns the minimum distance and the point.  If the
// origin itself is passed in the set, it is ignored.  If no points are passed at all, then
// a distance of 0 and a point of (0,0) is returned.
func calcMinManhattanDistance(points []point) (int, point) {
	if len(points) == 0 {
		return 0, point{}
	}
	minDistance := math.MaxInt64
	minPoint := point{}
	for _, point := range points {
		if point.x != 0 || point.y != 0 {
			distance := int(math.Abs(float64(point.x)) + math.Abs(float64(point.y)))
			if distance < minDistance {
				minDistance = distance
				minPoint = point
			}
		}
	}
	return minDistance, minPoint
}

func problem03A(fileName string) int {

	wireRoutes, err := getWireRoutes(fileName)
	if err != nil {
		log.Fatal(err)
	}

	wirePaths, err := convertRoutesToPaths(wireRoutes)
	if err != nil {
		log.Fatal(err)
	}

	wireSegments := getSegments(wirePaths)

	if len(wireSegments) < 2 {
		log.Fatalln("less than two wire segments were defiend")
	}
	if len(wireSegments) > 2 {
		log.Fatalln("more then two wires were defined")
	}
	intersections := getIntersections(wireSegments[0], wireSegments[1])

	minDistance, _ := calcMinManhattanDistance(intersections)

	return minDistance
}
