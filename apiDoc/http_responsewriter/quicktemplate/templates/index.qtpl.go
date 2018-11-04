//这个文件是由qtc从“hello.qtpl”自动生成的。
//有关详细信息，请参阅https://github.com/valyala/quicktemplate。
// Hello模板，实现了Partial的方法。

//line index.qtpl:3
package templates

//line index.qtpl:3
import (
	qtio422016 "io"
	qt422016 "github.com/valyala/quicktemplate"
)

//line index.qtpl:3
var (
	_ = qtio422016.Copy
	_ = qt422016.AcquireByteBuffer
)

//line index.qtpl:4
type Index struct{}

//line index.qtpl:7
func (i *Index) StreamBody(qw422016 *qt422016.Writer) {
	//line index.qtpl:7
	qw422016.N().S(`
	<h1>Index Page</h1>
	<div>
		This is our index page's body.
	</div>
`)
//line index.qtpl:12
}

//line index.qtpl:12
func (i *Index) WriteBody(qq422016 qtio422016.Writer) {
	//line index.qtpl:12
	qw422016 := qt422016.AcquireWriter(qq422016)
	//line index.qtpl:12
	i.StreamBody(qw422016)
	//line index.qtpl:12
	qt422016.ReleaseWriter(qw422016)
//line index.qtpl:12
}

//line index.qtpl:12
func (i *Index) Body() string {
	//line index.qtpl:12
	qb422016 := qt422016.AcquireByteBuffer()
	//line index.qtpl:12
	i.WriteBody(qb422016)
	//line index.qtpl:12
	qs422016 := string(qb422016.B)
	//line index.qtpl:12
	qt422016.ReleaseByteBuffer(qb422016)
	//line index.qtpl:12
	return qs422016
//line index.qtpl:12
}