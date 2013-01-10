//line github.com/fd/w/simplex/example/example.smplx:1

// vim: set ft=go :
package main

//line smplx_generated.go:3

//line smplx_generated.go:2
import sx_runtime "github.com/fd/w/simplex/runtime"

//line github.com/fd/w/simplex/example/example.smplx:5

//line github.com/fd/w/simplex/example/example.smplx:4
type Location struct {
	// declare as a view type

	name    string
	website string
	gps     GPS
}

//line github.com/fd/w/simplex/example/example.smplx:13

//line github.com/fd/w/simplex/example/example.smplx:12
type GPS struct {
//line github.com/fd/w/simplex/example/example.smplx:15
	Lat, Lng float64
}

//line github.com/fd/w/simplex/example/example.smplx:19

//line github.com/fd/w/simplex/example/example.smplx:18
var Locations = LocationView{}.
//line smplx_generated.go:4
Source().Select(func(m interface{}) bool {
//line smplx_generated.go:3
return has_website(m.(Location)) }).

//line github.com/fd/w/simplex/example/example.smplx:4

//line github.com/fd/w/simplex/example/example.smplx:18
Sort(func(m interface{}) interface{} {
//line github.com/fd/w/simplex/example/example.smplx:18
	return locations_gps_ne(m.(Location))
//line github.com/fd/w/simplex/example/example.smplx:4
},

//line github.com/fd/w/simplex/example/example.smplx:18
)
var LocationsWithWebsite = Locations.WithWebsite()
var Coordinates = GPSView{}.
//line smplx_generated.go:15
CollectedFrom(Locations, func(m interface{}) interface{} {
//line github.com/fd/w/simplex/example/example.smplx:20
	return gps(m.(Location))
//line github.com/fd/w/simplex/example/example.smplx:4
},

//line github.com/fd/w/simplex/example/example.smplx:20
).Sort(func(m interface{}) interface{} {
//line github.com/fd/w/simplex/example/example.smplx:20
	return gps_ne(m.(GPS))
//line github.com/fd/w/simplex/example/example.smplx:12
},

//line github.com/fd/w/simplex/example/example.smplx:20
)

//line github.com/fd/w/simplex/example/example.smplx:23

//line github.com/fd/w/simplex/example/example.smplx:22
func has_website(loc Location) bool {
	return loc.website != ""
}

//line github.com/fd/w/simplex/example/example.smplx:27

//line github.com/fd/w/simplex/example/example.smplx:26
func gps(loc Location) GPS {
	return loc.gps
}

//line github.com/fd/w/simplex/example/example.smplx:31

//line github.com/fd/w/simplex/example/example.smplx:30
func locations_gps_ne(loc Location) []float64 {
	return gps_ne(loc.gps)
}

//line github.com/fd/w/simplex/example/example.smplx:35

//line github.com/fd/w/simplex/example/example.smplx:34
func gps_ne(gps GPS) []float64 {
	return []float64{gps.Lat, gps.Lng}
}

//line github.com/fd/w/simplex/example/example.smplx:39

//line github.com/fd/w/simplex/example/example.smplx:38
func (v LocationView) WithWebsite() LocationView {
	return v.Select(func(m interface{}) bool {
//line github.com/fd/w/simplex/example/example.smplx:39
		return has_website(m.(Location))
//line github.com/fd/w/simplex/example/example.smplx:4
	},
//line github.com/fd/w/simplex/example/example.smplx:39
	)
}

//line smplx_generated.go:4

//line smplx_generated.go:3
type LocationView struct{ view sx_runtime.View }

//line smplx_generated.go:5

//line smplx_generated.go:4
func (w LocationView) Source() LocationView { return LocationView{sx_runtime.Source("LocationView")} }
func (w LocationView) Select(f sx_runtime.SelectFunc) LocationView {
//line smplx_generated.go:5
	return LocationView{w.view.Select(f)}
//line smplx_generated.go:5
}
func (w LocationView) Sort(f sx_runtime.SortFunc) LocationView   { return LocationView{w.view.Sort(f)} }
func (w LocationView) Group(f sx_runtime.GroupFunc) LocationView { return LocationView{w.view.Group(f)} }
func (w LocationView) CollectedFrom(input sx_runtime.ViewWrapper, f sx_runtime.CollectFunc) LocationView {
//line smplx_generated.go:8
	return LocationView{input.View().Collect(f)}
//line smplx_generated.go:8
}
func (w LocationView) View() sx_runtime.View { return w.view }

//line smplx_generated.go:11

//line smplx_generated.go:10
type GPSView struct{ view sx_runtime.View }

//line smplx_generated.go:12

//line smplx_generated.go:11
func (w GPSView) Source() GPSView                        { return GPSView{sx_runtime.Source("GPSView")} }
func (w GPSView) Select(f sx_runtime.SelectFunc) GPSView { return GPSView{w.view.Select(f)} }
func (w GPSView) Sort(f sx_runtime.SortFunc) GPSView     { return GPSView{w.view.Sort(f)} }
func (w GPSView) Group(f sx_runtime.GroupFunc) GPSView   { return GPSView{w.view.Group(f)} }
func (w GPSView) CollectedFrom(input sx_runtime.ViewWrapper, f sx_runtime.CollectFunc) GPSView {
//line smplx_generated.go:15
	return GPSView{input.View().Collect(f)}
//line smplx_generated.go:15
}
func (w GPSView) View() sx_runtime.View { return w.view }
