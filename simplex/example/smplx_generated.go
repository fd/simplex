//line github.com/fd/w/simplex/example/example.smplx:1

// vim: set ft=go :
package main

//line smplx_generated.go:3

//line smplx_generated.go:2
import sx_runtime "github.com/fd/w/simplex/runtime"

//line github.com/fd/w/simplex/example/example.smplx:5

//line github.com/fd/w/simplex/example/example.smplx:4
type Location struct {
	name    string
	website string
	gps     GPS
}

//line github.com/fd/w/simplex/example/example.smplx:11

//line github.com/fd/w/simplex/example/example.smplx:10
type GPS struct {
	Lat, Lng float64
}

//line github.com/fd/w/simplex/example/example.smplx:15

//line github.com/fd/w/simplex/example/example.smplx:14
var Locations = LocationViewSource().Where(func(m interface {
//line github.com/fd/w/simplex/example/example.smplx:14
}) bool {
//line github.com/fd/w/simplex/example/example.smplx:14
return has_website(m.(Location)) }).
//line github.com/fd/w/simplex/example/example.smplx:14
Sort(func(m interface{}) interface{} {
//line github.com/fd/w/simplex/example/example.smplx:14
return locations_gps_ne(m.(Location)) })

//line github.com/fd/w/simplex/example/example.smplx:14
var LocationsWithWebsite = Locations.WithWebsite()
var Coordinates = GPSViewCollectedFrom(Locations, func(m interface{}) interface{} {
//line github.com/fd/w/simplex/example/example.smplx:16
return gps(m.(Location)) }).
//line github.com/fd/w/simplex/example/example.smplx:16
Sort(func(m interface{}) interface{} {
//line github.com/fd/w/simplex/example/example.smplx:16
return gps_ne(m.(GPS)) })

//line github.com/fd/w/simplex/example/example.smplx:16
//line github.com/fd/w/simplex/example/example.smplx:19

//line github.com/fd/w/simplex/example/example.smplx:18
func has_website(loc Location) bool {
	return loc.website != ""
}

//line github.com/fd/w/simplex/example/example.smplx:23

//line github.com/fd/w/simplex/example/example.smplx:22
func gps(loc Location) GPS {
	return loc.gps
}

//line github.com/fd/w/simplex/example/example.smplx:27

//line github.com/fd/w/simplex/example/example.smplx:26
func locations_gps_ne(loc Location) []float64 {
	return gps_ne(loc.gps)
}

//line github.com/fd/w/simplex/example/example.smplx:31

//line github.com/fd/w/simplex/example/example.smplx:30
func gps_ne(gps GPS) []float64 {
	return []float64{gps.Lat, gps.Lng}
}

//line github.com/fd/w/simplex/example/example.smplx:35

//line github.com/fd/w/simplex/example/example.smplx:34
func (v LocationView) WithWebsite() LocationView {
	return v.Where(func(m interface{}) bool {
//line github.com/fd/w/simplex/example/example.smplx:35
	return has_website(m.(Location)) })
//line github.com/fd/w/simplex/example/example.smplx:35
}

//line smplx_generated.go:4

//line smplx_generated.go:3
type LocationView struct{ view sx_runtime.View }

//line smplx_generated.go:5

//line smplx_generated.go:4
func LocationViewSource() LocationView                           { return LocationView{sx_runtime.Source("LocationView")} }
func (w LocationView) Where(f sx_runtime.WhereFunc) LocationView { return LocationView{w.view.Where(f)} }
func (w LocationView) Sort(f sx_runtime.SortFunc) LocationView   { return LocationView{w.view.Sort(f)} }
func (w LocationView) Group(f sx_runtime.GroupFunc) LocationView { return LocationView{w.view.Group(f)} }
func LocationViewCollectedFrom(input sx_runtime.ViewWrapper, f sx_runtime.CollectFunc) LocationView {
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
func GPSViewSource() GPSView                           { return GPSView{sx_runtime.Source("GPSView")} }
func (w GPSView) Where(f sx_runtime.WhereFunc) GPSView { return GPSView{w.view.Where(f)} }
func (w GPSView) Sort(f sx_runtime.SortFunc) GPSView   { return GPSView{w.view.Sort(f)} }
func (w GPSView) Group(f sx_runtime.GroupFunc) GPSView { return GPSView{w.view.Group(f)} }
func GPSViewCollectedFrom(input sx_runtime.ViewWrapper, f sx_runtime.CollectFunc) GPSView {
//line smplx_generated.go:15
	return GPSView{input.View().Collect(f)}
//line smplx_generated.go:15
}
func (w GPSView) View() sx_runtime.View { return w.view }
