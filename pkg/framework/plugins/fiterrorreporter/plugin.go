package fiterrorreporter

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"os"
	"text/tabwriter"
)

const Name = "FitErrorReporter"

type FitErrorReporter struct {
	client kubernetes.Interface
}

func New(client kubernetes.Interface, _ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
	return &FitErrorReporter{
		client: client,
	}, nil
}

func (pl *FitErrorReporter) Name() string {
	return Name
}

func (pl *FitErrorReporter) PostFilter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, filteredNodeStatusMap framework.NodeToStatusMap) (*framework.PostFilterResult, *framework.Status) {

	w := tabwriter.NewWriter(os.Stdout, 0, 8, 2, ' ', tabwriter.AlignRight)

	for nodeName, st := range filteredNodeStatusMap {
		_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", pod.Name, nodeName, st.Code(), st.Plugin(), st.Message())
		if err != nil {
			fmt.Println(w, "ERROR")
		}
	}
	err := w.Flush()
	if err != nil {
		fmt.Println(w, "ERROR")
	}
	return nil, nil
}
