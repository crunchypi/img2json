package main

import (
	"fmt"
	"img2json/src/filters"
	"img2json/src/filters/color"
	"img2json/src/points"
	"os"
	"strconv"
	"strings"
)

var helpString = `
____________________________________________________________
Image to JSON converter with a bit of extra filtering stuff.
____________________________________________________________

This cli uses piping, meaning that argument positioning does
matter. The arguments are as follows (examples below):

    -help             Prints out this message.
    -ii      <path>   Specifies the 'Input Image' path.
    -ij      <path>   Specifies the 'Input JSON' path.
    -oi      <path>   Specifies the 'Output Image' path.
    -oj      <path>   Specifies the 'Output JSON' path.
    -byColor <...>    Accepts a value with the form:
                        int,int,int,int,int,int
                      .. which represents an RGB range.
                      So the first three ints are min RGB,
                      and the three latter are max RGB.
                      Everything outside of that range is
                      filtered out, i.e removed.
    -byRand  <int>    Filters a specified percentage (0..100)
                      of all currently stored pixels.
					
Examples:

	.. To input an image, filter 50% of randomly selected
	pixels (uniform selection) and store as a new image:
		> -ii ./my.png -byRand 50 -oi ./my2.png

	.. To Input an image, filter out anything that is not
	absolute darkness, then store it as a JSON:
		> -ii ./my.png -byColor 0,0,0,0,0,0 -oj ./my2.json

	.. To import a JSON file created with this CLI and
	save it as an image:
		> -ij ./my.json -oi > -oi ./my.png

`

// container for task funcs associated with cli switches, and
// a bool which specifies whether or not the aforementioned
// cli switches should expect a following value.
type argHelper struct {
	f          func(*points.Points, string) string
	acceptsArg bool
}

// returns a map of command-line switches and their associated
// 'argHelper'.
func cliSwitches() map[string]argHelper {
	return map[string]argHelper{
		"-help":    {help, false},   // # Help printout.
		"-ii":      {ii, true},      // # Image input.
		"-ij":      {ij, true},      // # JSON input.
		"-oi":      {oi, true},      // # Image output.
		"-oj":      {oj, true},      // # JSON output.
		"-byColor": {byColor, true}, // # Filter by color.
		"-byRand":  {byRand, true},  // # Filter by uniform rand.
	}
}

// Command-line interpreter.
func main() {
	// # Nice-2-have.
	if len(os.Args) == 1 {
		help(nil, "")
		return
	}

	// # State to pass between task funcs.
	p := points.Points{}

	// # Exhaust arguments. Strange looping because some switches
	// # expect a following argument (k:v), and having the ability
	// # to skip by a varying amount of steps allows for this.
	i := 1
	for {
		// # Exit clause.
		if i >= len(os.Args) {
			return
		}
		// # Skip next i by this val (default).
		step := 1
		helper, ok := cliSwitches()[os.Args[i]]
		if !ok {
			fmt.Printf("switch '%s' isn't recognized\n", os.Args[i])
			os.Exit(0)
		}

		// # Collects error message from task funcs.
		errMsg := ""
		switch helper.acceptsArg {
		case true:
			// # This switch expects a following val, so get
			// # that and increase the next step.
			if i+1 >= len(os.Args) {
				fmt.Printf("switch '%s' expects a value but got nothing\n", os.Args[i])
				os.Exit(0)
			}
			errMsg = helper.f(&p, os.Args[i+1])
			step++
		case false:
			errMsg = helper.f(&p, "")
		}

		if len(errMsg) != 0 {
			fmt.Printf("issue on '%s': %s\n", os.Args[i], errMsg)
		}

		i += step
	}
}

// for -help switch
func help(pointState *points.Points, arg string) string {
	fmt.Println(helpString)
	return ""
}

// for -ii switch.
func ii(pointState *points.Points, arg string) string {
	p, err := points.NewFromImageFile(arg)
	if err != nil {
		return fmt.Sprintf("file at path '%s' either doesn't exist or is not an img", arg)
	}

	*pointState = *p
	return ""
}

// for -ij switch.
func ij(pointState *points.Points, arg string) string {
	p, err := points.NewFromJSONFile(arg)
	if err != nil {
		return fmt.Sprintf("file at path '%s' either doesn't exist or isn't in json", arg)
	}

	*pointState = *p
	return ""
}

// for -oi switch.
func oi(pointState *points.Points, arg string) string {
	if err := pointState.SaveAsImage(arg); err != nil {
		return fmt.Sprintf("issue while saving to '%s': %s", arg, err.Error())
	}
	return ""
}

// for -ooj switch.
func oj(pointState *points.Points, arg string) string {
	if err := pointState.SaveAsJSON(arg); err != nil {
		return fmt.Sprintf("issue while saving to '%s': %s", arg, err.Error())
	}
	return ""
}

// for -byColor switch.
func byColor(pointState *points.Points, arg string) string {
	args := strings.Split(arg, ",")
	if len(args) != 6 {
		return fmt.Sprintf("expected 6 args, got %d", len(args))
	}
	argsInts := make([]uint8, 6)
	for i, arg := range args {
		num, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Sprintf("arg no. %d is not an integer", i)
		}
		if num < 0 || num > 255 {
			return fmt.Sprintf("arg no. %d is not in range 0..255", i)
		}
		argsInts[i] = uint8(num)
	}
	colors := color.ColorBounds{
		argsInts[0], argsInts[1], argsInts[2],
		argsInts[3], argsInts[4], argsInts[5],
	}
	filters.ByColor(pointState, colors)
	return ""
}

// for -byRand switch.
func byRand(pointState *points.Points, arg string) string {
	percent, err := strconv.Atoi(arg)
	if err != nil {
		return fmt.Sprintf("arg is not a integer")
	}
	if percent < 0 || percent > 100 {
		return fmt.Sprintf("arg is not in range 0..100")
	}
	filters.ByRand(pointState, float32(percent)/100)
	return ""
}
