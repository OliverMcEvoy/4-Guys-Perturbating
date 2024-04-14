package main

import (
	"fmt"
	"math"
	"math/cmplx"
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
	xLength := float64(10)
	zLength := float64(10)
	yLength := float64(10)

	points, _ := genPoints(float64(time.Now().Unix()))
	mats := plotPoints(scene, points)

	createGraph(scene, xLength, zLength, yLength)

	// Create and add lights to the scene
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	scene.Add(pointLight)

	// Create and add an axis helper to the scene
	scene.Add(helper.NewAxes(0.5))

	// Set background color to gray
	a.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	//startTime := float64(time.Now().Unix())
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

		_, waveFunc := genPoints(float64(time.Now().Unix()))

		vals := minMaxNormalize(waveFunc)
		for i := 0; i < len(points); i++ {
			mats[i].SetColor(GenerateColorOnGradient(vals[i]))
		}

		// Update GUI timers
		gui.Manager().TimerManager.ProcessTimers()

		// Control and update FPS
		rater.Wait()
	})
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
	var mats []*material.Standard
	for i := 0; i < len(points); i++ {
		points[i] = NormalizeVals(points[i])
		geom := geometry.NewCube(0.5)
		//TODO This is the maths function call is meant to go we if the output here and assign colour
		mat := material.NewStandard(math32.NewColor("Red"))
		mats = append(mats, mat)
		mesh := graphic.NewMesh(geom, mat)
		mesh.SetPosition(float32(points[i][0]), float32(points[i][2]), float32(points[i][1]))
		scene.Add(mesh)
		fmt.Println(float32(points[i][0]), float32(points[i][2]), float32(points[i][1]))
	}
	return mats
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
	red := 1 - float32(value)
	blue := float32(value)
	green := float32(0)

	return &math32.Color{R: red, G: green, B: blue}
}

type quantumSystem struct {
	quantumNumbers        [3]int
	energy                float64
	wellLength            float64
	reducedPlanckConstant float64
}

type particle struct {
	position [3]float64
	mass     float64
}

func initialiseQuantumSystem(q *quantumSystem, p *particle) {
	p.mass = 9.10938356e-31 // Electron mass
	p.position = [3]float64{0.0, 0.0, 0.0}

	q.reducedPlanckConstant = 1.0545718e-34 // Planck constant divided by 2π
	q.wellLength = 10e-9                    // Further increase well length
	q.quantumNumbers = [3]int{3, 3, 3}
}

func waveFunction(q *quantumSystem, p *particle, time float64) float64 {
	waveFunctionValue := 0.0

	// if any of the positions are outside the well, return 0
	for i := 0; i < 3; i++ {
		if p.position[i] < 0.0 || p.position[i] > q.wellLength {
			// fmt.Printf("Position outside well, position = %.6e\n", p.position[i])
			return waveFunctionValue
		}
	}

	time = time * 1e-9 // Convert time to nanoseconds
	q.energy = math.Pow(math.Pi*q.reducedPlanckConstant, 2) / (2 * p.mass * math.Pow(q.wellLength, 2))

	quantumNumberSquareSum := 0
	for i := 0; i < 3; i++ {
		quantumNumberSquareSum += q.quantumNumbers[i] * q.quantumNumbers[i]
	}

	q.energy *= float64(quantumNumberSquareSum)

	waveFunctionValue = 2.0 * math.Sqrt(8.0/math.Pow(q.wellLength, 3))
	for i := 0; i < 3; i++ {
		waveFunctionValue *= math.Sin(float64(q.quantumNumbers[i]) * math.Pi * p.position[i] / q.wellLength)
	}

	timeDependentFactor := cmplx.Exp(-1i * complex(q.energy*time/q.reducedPlanckConstant, 0))

	return cmplx.Abs(complex(waveFunctionValue, 0) * timeDependentFactor)
}

func genPoints(t float64) ([][]float64, []float64) {
	// define quantum system and particle
	q := quantumSystem{}
	p := particle{}

	initialiseQuantumSystem(&q, &p)

	gridLength := 10e-9 // Further increase spatial dimension
	divisions := 10     // Keep divisions for larger step sizes

	xStep := float64(q.wellLength) / float64(divisions)
	yStep := xStep
	zStep := xStep

	var waveFuncs []float64
	var coords [][]float64
	for x := 0.0; x < float64(gridLength)+xStep; x += xStep {
		for y := 0.0; y < float64(gridLength)+yStep; y += yStep {
			for z := 0.0; z < float64(gridLength)+zStep; z += zStep {
				p.position = [3]float64{x, y, z}
				waveFunctionValue := waveFunction(&q, &p, t)
				waveFuncs = append(waveFuncs, waveFunctionValue)
				coord := []float64{x, y, z}
				coords = append(coords, coord)
			}
		}
	}
	return coords, waveFuncs
}
func NormalizeVals(vals []float64) []float64 {

	var empty []float64

	for i := 0; i < len(vals); i++ {
		empty = append(empty, vals[i]*1000000000.0)
	}
	return empty
}

func minMaxNormalize(vals []float64) []float64 {
	// Find the minimum and maximum values
	minVal, maxVal := vals[0], vals[0]
	for _, val := range vals {
		if val < minVal {
			minVal = val
		}
	}
	for _, val := range vals {
		if val > maxVal {
			maxVal = val
		}
	}

	// Create a new slice for the normalized values
	normalizedVals := make([]float64, len(vals))

	// Normalize the values
	for i, val := range vals {
		normalizedVals[i] = (val - minVal) / (maxVal - minVal)
	}

	return normalizedVals
}
