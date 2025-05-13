package utils

// AnnotationMapper service loading abstraction for annotations that have different group name
type AnnotationMapper interface {
	Find(annotations map[string]string, key string) (string, bool)
	AddPrefix(annotations map[string]string) map[string]string
}

type ResourceGroupAnnotationsMapper struct {
	groups []string
}

func NewResourceGroupAnnotationsMapper(groups ...string) *ResourceGroupAnnotationsMapper {
	return &ResourceGroupAnnotationsMapper{groups: groups}
}

func (g *ResourceGroupAnnotationsMapper) Find(annotations map[string]string, key string) (string, bool) {
	for _, v := range g.groups {
		if value, found := annotations[v+"/"+key]; found {
			return value, true
		}
	}
	return "", false
}

func (g *ResourceGroupAnnotationsMapper) AddPrefix(annotations map[string]string) map[string]string {
	labeledAnnotations := make(map[string]string)
	for _, group := range g.groups {
		for k, v := range annotations {
			labeledAnnotations[group+"/"+k] = v
		}
	}
	return labeledAnnotations
}
