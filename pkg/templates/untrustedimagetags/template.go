package untrustedimagetags

import (
	"fmt"
	"strings"

	"golang.stackrox.io/kube-linter/pkg/check"
	"golang.stackrox.io/kube-linter/pkg/config"
	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/objectkinds"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/util"
	v1 "k8s.io/api/core/v1"
	"golang.stackrox.io/kube-linter/pkg/templates/untrustedimagetags/internal/params"
)

const (
	templateKey = "untrusted-image-tags"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Untrusted image tags",
		Key:         templateKey,
		Description: "Flag applications running containers with untrusted image tags",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {

				// checking whether the container image has any image tag out of the list of not-allowed image tags passed as the paramater

				res := false

				for _, imageTag := range p.NotAllowedTags {
					res = strings.Contains(container.Image, imageTag)
				}
				if res {
					return []diagnostic.Diagnostic{{Message: fmt.Sprintf("container %q image doesn't have a valid image tag", container.Name)}}
				}
				return nil
			}), nil
		}),
	})
}
