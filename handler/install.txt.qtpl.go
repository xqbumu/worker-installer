// Code generated by qtc from "install.txt.qtpl". DO NOT EDIT.
// See https://github.com/valyala/quicktemplate for details.

//line handler/install.txt.qtpl:1
package handler

//line handler/install.txt.qtpl:1
import (
	qtio422016 "io"

	qt422016 "github.com/valyala/quicktemplate"
)

//line handler/install.txt.qtpl:1
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line handler/install.txt.qtpl:1
func StreamText(qw422016 *qt422016.Writer, r Result) {
//line handler/install.txt.qtpl:1
	qw422016.N().S(`
repository: https://github.com/`)
//line handler/install.txt.qtpl:2
	qw422016.E().S(r.User)
//line handler/install.txt.qtpl:2
	qw422016.N().S(`/`)
//line handler/install.txt.qtpl:2
	qw422016.E().S(r.Program)
//line handler/install.txt.qtpl:2
	qw422016.N().S(`
user: `)
//line handler/install.txt.qtpl:3
	qw422016.E().S(r.User)
//line handler/install.txt.qtpl:3
	qw422016.N().S(`
program: `)
//line handler/install.txt.qtpl:4
	qw422016.E().S(r.Program)
//line handler/install.txt.qtpl:4
	if len(r.AsProgram) > 0 {
//line handler/install.txt.qtpl:4
		qw422016.N().S(`
as: `)
//line handler/install.txt.qtpl:5
		qw422016.E().S(r.AsProgram)
//line handler/install.txt.qtpl:5
	}
//line handler/install.txt.qtpl:5
	qw422016.N().S(`
release: `)
//line handler/install.txt.qtpl:6
	qw422016.E().S(r.ResolvedRelease)
//line handler/install.txt.qtpl:6
	qw422016.N().S(`
move-into-path: `)
//line handler/install.txt.qtpl:7
	qw422016.E().V(r.MoveToPath)
//line handler/install.txt.qtpl:7
	qw422016.N().S(`
sudo-move: `)
//line handler/install.txt.qtpl:8
	qw422016.E().V(r.SudoMove)
//line handler/install.txt.qtpl:8
	qw422016.N().S(`
used-search: `)
//line handler/install.txt.qtpl:9
	qw422016.E().V(r.Search)
//line handler/install.txt.qtpl:9
	qw422016.N().S(`
asset-select: "`)
//line handler/install.txt.qtpl:10
	if len(r.Select) > 0 {
//line handler/install.txt.qtpl:10
		qw422016.N().S(` `)
//line handler/install.txt.qtpl:10
		qw422016.E().S(r.Select)
//line handler/install.txt.qtpl:10
		qw422016.N().S(` `)
//line handler/install.txt.qtpl:10
	}
//line handler/install.txt.qtpl:10
	qw422016.N().S(`"

release assets:
`)
//line handler/install.txt.qtpl:13
	for _, i := range r.Assets {
//line handler/install.txt.qtpl:13
		qw422016.N().S(`  `)
//line handler/install.txt.qtpl:13
		qw422016.E().S(i.Key())
//line handler/install.txt.qtpl:13
		qw422016.N().S(`
    url:    `)
//line handler/install.txt.qtpl:14
		qw422016.E().S(i.URL)
//line handler/install.txt.qtpl:14
		qw422016.N().S(` `)
//line handler/install.txt.qtpl:14
		if len(i.SHA256) != 0 {
//line handler/install.txt.qtpl:14
			qw422016.N().S(`
    sha256: `)
//line handler/install.txt.qtpl:15
			qw422016.E().S(i.SHA256)
//line handler/install.txt.qtpl:15
		}
//line handler/install.txt.qtpl:15
		qw422016.N().S(`
`)
//line handler/install.txt.qtpl:16
	}
//line handler/install.txt.qtpl:16
	qw422016.N().S(`
has-m1-asset: `)
//line handler/install.txt.qtpl:17
	qw422016.E().V(r.M1Asset)
//line handler/install.txt.qtpl:17
	qw422016.N().S(`

to see shell script, append ?type=script
for more information on this server, visit:
  github.com/xqbumu/worker-installer
`)
//line handler/install.txt.qtpl:22
}

//line handler/install.txt.qtpl:22
func WriteText(qq422016 qtio422016.Writer, r Result) {
//line handler/install.txt.qtpl:22
	qw422016 := qt422016.AcquireWriter(qq422016)
//line handler/install.txt.qtpl:22
	StreamText(qw422016, r)
//line handler/install.txt.qtpl:22
	qt422016.ReleaseWriter(qw422016)
//line handler/install.txt.qtpl:22
}

//line handler/install.txt.qtpl:22
func Text(r Result) string {
//line handler/install.txt.qtpl:22
	qb422016 := qt422016.AcquireByteBuffer()
//line handler/install.txt.qtpl:22
	WriteText(qb422016, r)
//line handler/install.txt.qtpl:22
	qs422016 := string(qb422016.B)
//line handler/install.txt.qtpl:22
	qt422016.ReleaseByteBuffer(qb422016)
//line handler/install.txt.qtpl:22
	return qs422016
//line handler/install.txt.qtpl:22
}
