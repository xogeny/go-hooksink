# A Go based GitHub WebHook Listener

This is a simple package I created for my own purposes but I'm happy
to share it with other people.  At this point, I consider this
**alpha** code.  But my hope is to refine it into a solid chunk of
reusable code for easily creating Go based servers for handling
webhooks.  I have several applications in mind for this and I suspect
it will change significantly as I go.

## Example

The goal of this Go package is to allow users to easily create a
GitHub webhook server and attach handlers to it.  Different handlers
can, for example, be attached to different URLs by specifying
different paths.  In this example, a simple logging handler is
attached to `/log`:

```
package main

import "log"

import hs "github.com/xogeny/go-hooksink"

type Logger struct {}

func (l Logger) Push(msg hs.HubMessage) {
	log.Printf("PUSH: %v", msg);
}

func main() {
	h := hs.NewHookSink();
	h.Add("/log", Logger{});
	h.Start();
o}
```

## Acknowledgments

When I went looking for existing Go language libraries for creating
servers to handle GitHub webhook calls, I came across
[dockerhub-webhook-listener](https://github.com/cpuguy83/dockerhub-webhook-listener)
by Brian Goff (@cpuguy83).  I originally planned to reuse the code in
his repository (which was mostly generic, but a little bit Docker
specific) and refactor it into a more generic library.  But I realized
that once I started doing that there were so many things I wanted to
change that I felt like it would be better to start from scratch.

But this work was inspired considerably by his design and his approach
still supports several things that I currently don't (TLS, key
checking, etc).
