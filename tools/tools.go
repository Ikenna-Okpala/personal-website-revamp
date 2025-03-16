package main

import (
	"crypto/rand"
	"fmt"
)

var x = `
<div>
    <div id="counter">
      {{.Counter}}
    </div>
    <div>
      <button hx-target="#counter" hx-post="/decrease" hx-swap=”outerHTML”>
        Decrease -
      </button>
      <button hx-target="#counter" hx-post="/increase" hx-swap=”outerHTML”>
        Increase +
      </button>
    </div>
  </div>

`

func main() {

	password:= rand.Text()

  fmt.Println(password)
}