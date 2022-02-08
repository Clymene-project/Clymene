package kafka

import (
	"github.com/Clymene-project/Clymene/model/labels"
	"github.com/Clymene-project/Clymene/model/relabel"
	util "github.com/Clymene-project/Clymene/pkg/lokiutil"
	"strings"

	"github.com/prometheus/common/model"
)

func format(lbs labels.Labels, cfg []*relabel.Config) model.LabelSet {
	if len(lbs) == 0 {
		return nil
	}
	processed := relabel.Process(lbs, cfg...)
	labelOut := model.LabelSet(util.LabelsToMetric(processed))
	for k := range labelOut {
		if strings.HasPrefix(string(k), "__") {
			delete(labelOut, k)
		}
	}
	return labelOut
}
