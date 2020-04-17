package status

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/alex-kennedy/wikilinks/prep/pipeline"
	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

//statusColors are the colors of graph bubbles in the graph.
var statusColors = map[pipeline.TaskStatus]string{
	pipeline.DoneStatus:       "#81C784",
	pipeline.InProgressStatus: "#FFD54F",
	pipeline.NotDoneStatus:    "#90A4AE",
}

//PageData contains the information needed for rendering the status page.
type PageData struct {
	Image          template.HTML
	CSS            template.CSS
	PipelineStatus string
	TargetTask     string
}

//Status captures the functionality of the pipeline status website.
type Status struct {
	Pipeline *pipeline.Pipeline
	Template *template.Template
	CSS      *template.CSS
}

//RenderPage is the handler for the web status.
func (s *Status) RenderPage(w http.ResponseWriter, r *http.Request) {
	graph, err := s.RenderGraph()
	if err != nil {
		graph = "Error making graph!"
		log.Print(err)
	}
	var pipelineStatus string
	if s.Pipeline.PipelineError {
		pipelineStatus = "Error!"
	} else if s.Pipeline.PipelineRunning {
		pipelineStatus = "Running."
	} else {
		pipelineStatus = "Done."
	}
	targetTask := fmt.Sprintf("%T", s.Pipeline.RootTask)
	data := PageData{template.HTML(graph), *s.CSS, pipelineStatus, targetTask}
	s.Template.Execute(w, data)
}

//RenderGraph renders the current task status as graphviz SVG graph.
func (s *Status) RenderGraph() (string, error) {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			log.Fatal(err)
		}
		g.Close()
	}()

	//Create nodes
	nodes := make(map[string]*cgraph.Node)
	for _, task := range *s.Pipeline.SortedTasks {
		nodeName := fmt.Sprintf("%T", task)
		node, err := graph.CreateNode(fmt.Sprintf("%T", task))
		if err != nil {
			return "", err
		}
		nodes[nodeName] = node
		node.SetStyle("filled")
		node.SetFillColor(statusColors[s.Pipeline.TaskStatuses[nodeName]])
	}

	for _, task := range *s.Pipeline.SortedTasks {
		nodeName := fmt.Sprintf("%T", task)
		for _, dep := range task.Deps() {
			depName := fmt.Sprintf("%T", dep)
			_, err := graph.CreateEdge("", nodes[depName], nodes[nodeName])
			if err != nil {
				return "", err
			}
		}
	}

	var buf bytes.Buffer
	if err := g.Render(graph, graphviz.SVG, &buf); err != nil {
		log.Fatal(err)
	}
	return buf.String(), nil
}

//NewStatusSite initializes the site's required static pages.
func NewStatusSite(p *pipeline.Pipeline) *Status {
	pageTemplate := template.Must(template.ParseFiles("prep/status/status.template.html"))
	cssRaw, err := ioutil.ReadFile("prep/status/styles.css")
	if err != nil {
		panic(err)
	}
	css := template.CSS(cssRaw)
	return &Status{p, pageTemplate, &css}
}
