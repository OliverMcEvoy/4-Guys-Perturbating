package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/util"
	"github.com/g3n/engine/util/helper"
	"github.com/g3n/engine/window"
	"gonum.org/v1/gonum/stat"
)

func main() {
	// Create application and scene
	a := app.App()
	scene := core.NewNode()
	rater := util.NewFrameRater(60)

	// Set the scene to be managed by the gui manager
	gui.Manager().Set(scene)

	// Create perspective camera
	cam := camera.New(1)
	cam.SetPosition(0, 0, 3)
	scene.Add(cam)

	// Set up orbit control for the camera
	camera.NewOrbitControl(cam)

	// Set up callback to update viewport and camera aspect ratio when the window is resized
	onResize := func(evname string, ev interface{}) {
		// Get framebuffer size and update viewport accordingly
		width, height := a.GetSize()
		a.Gls().Viewport(0, 0, int32(width), int32(height))
		// Update the camera's aspect ratio
		cam.SetAspect(float32(width) / float32(height))
	}
	a.Subscribe(window.OnWindowSize, onResize)
	onResize("", nil)

	// axis length
	//TODO function call to set maxs
	xLength := float64(15)
	zLength := float64(15)
	yLength := float64(15)

	//if max is 10 min is -5

	createGraph(scene, xLength, zLength, yLength)
	points := generateRandomCoords(25, 0, xLength, 0, yLength, 0, zLength)

	mats := plotPoints(scene, points)

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	scene.Add(pointLight)

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(0.5))

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	a.Run(func(rend *renderer.Renderer, deltaTime time.Duration) {
		// Start measuring this frame
		rater.Start()

		// Clear the color, depth, and stencil buffers
		a.Gls().Clear(gls.COLOR_BUFFER_BIT | gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT) // TODO maybe do inside renderer, and allow customization

		// Render scene
		err := rend.Render(scene, cam)
		if err != nil {
			panic(err)
		}

		// Update positions and velocities
		for i := 0; i < len(points); i++ {
			// Update position
			points[i][0][0] += points[i][1][0]
			points[i][0][1] += points[i][1][1]
			points[i][0][2] += points[i][1][2]
			// Check for collision with walls and update velocity
			if points[i][0][0] < 0 || points[i][0][0] > xLength {
				points[i][1][0] *= -1
			}
			if points[i][0][1] < 0 || points[i][0][1] > yLength {
				points[i][1][1] *= -1
			}
			if points[i][0][2] < 0 || points[i][0][2] > zLength {
				points[i][1][2] *= -1
			}
		}

		var vals []float64
		for i := 0; i < len(points); i++ {
			t := float64(time.Now().Unix())
			val := calculateWaveFunction(points[i][0][0], points[i][0][1], points[i][0][2], t)
			vals = append(vals, val)
		}
		vals = NormalizeVals(vals)
		for i := 0; i < len(mats); i++ {
			mats[i].SetColor(GenerateColorOnGradient((0 + vals[i]*(1))))
		}

		// Update GUI timers
		gui.Manager().TimerManager.ProcessTimers()

		// Control and update FPS
		rater.Wait()
	})
}

func calculateWaveFunction(x, y, z, t float64) float64 {
	var result float64

	result = math.Sqrt(5) * math.Sin(math.Pi*(x+t)/5) * math.Sin(math.Pi*(y+t)/10) * math.Sin(math.Pi*(z+t)/5) / 25

	return float64(result)
}

func createGraph(scene *core.Node, xLength, zLength, yLength float64) {
	// x axis
	geomX := geometry.NewBox(float32(xLength), 0.05, 0.05)
	matX := material.NewStandard(math32.NewColor("DarkBlue"))
	meshX := graphic.NewMesh(geomX, matX)
	meshX.SetPosition(float32(xLength/2), 0, 0)
	scene.Add(meshX)
	// z axis
	geomZ := geometry.NewBox(0.05, float32(zLength), 0.05)
	matZ := material.NewStandard(math32.NewColor("DarkBlue"))
	meshZ := graphic.NewMesh(geomZ, matZ)
	meshZ.SetPosition(0, float32(zLength/2), 0)
	scene.Add(meshZ)
	// y axis
	geomY := geometry.NewBox(0.05, 0.05, float32(yLength))
	matY := material.NewStandard(math32.NewColor("DarkBlue"))
	meshY := graphic.NewMesh(geomY, matY)
	meshY.SetPosition(0, 0, float32(yLength/2))
	scene.Add(meshY)
}

func plotPoints(scene *core.Node, points [][]float64) []*material.Standard {
	var Array []*material.Standard
	for i := 0; i < len(points); i++ {
		geom := geometry.NewCube(0.1)
		//TODO This is the maths function call is meant to go we if the output here and assign colour
		mat := material.NewStandard(math32.NewColor("DarkBlue"))
		Array = append(Array, mat)
		mesh := graphic.NewMesh(geom, mat)
		mesh.SetPosition(float32(points[i][0]), float32(points[i][2]), float32(points[i][1]))
		scene.Add(mesh)
	}
	return Array
}

// generateRandomCoords generates random arrays of float64 within specified bounds
func generateRandomCoords(numCoords int, minX, maxX, minY, maxY, minZ, maxZ float64) [][]float64 {
	// Create the 2D array to store the random arrays
	result := make([][]float64, numCoords)
	r := rand.New(rand.NewSource(38))

	// Generate random values for each array
	for i := 0; i < numCoords; i++ {
		result[i] = make([]float64, 3)
		result[i][0] = minX + r.Float64()*(maxX-minX) // Random value for x
		result[i][1] = minY + r.Float64()*(maxY-minY) // Random value for y
		result[i][2] = minZ + r.Float64()*(maxZ-minZ) // Random value for z
	}
	return result
}

// GenerateColorOnGradient generates a color on a gradient from red to blue based on the input value (0 to 1)
func GenerateColorOnGradient(value float64) *math32.Color {
	// Ensure the value is within the range of 0 to 1
	if value < 0 {
		value = 0
	} else if value > 1 {
		value = 1
	}

	// Interpolate between red and blue based on the input value
	red := float32(0)
	blue := float32(value)
	green := 0.7 - float32(value)/2

	return &math32.Color{R: red, G: green, B: blue}
}
func NormalizeVals(vals []float64) []float64 {
	valsMean := stat.Mean(vals, nil)
	valsStdDev := stat.StdDev(vals, nil)

	normalizedVals := make([]float64, len(vals))
	for i, val := range vals {
		normalizedVals[i] = (val - valsMean) / valsStdDev
	}

	return normalizedVals
}
func generateRandomVelocities(numPoints int, minVelocity, maxVelocity float64) [][]float64 {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Create the 2D slice to store the velocities
	velocities := make([][]float64, numPoints)

	// Generate random velocities for each point
	for i := 0; i < numPoints; i++ {
		velocities[i] = make([]float64, 3)
		velocities[i][0] = minVelocity + rand.Float64()*(maxVelocity-minVelocity) // Random velocity for x
		velocities[i][1] = minVelocity + rand.Float64()*(maxVelocity-minVelocity) // Random velocity for y
		velocities[i][2] = minVelocity + rand.Float64()*(maxVelocity-minVelocity) // Random velocity for z
	}

	return velocities
}
