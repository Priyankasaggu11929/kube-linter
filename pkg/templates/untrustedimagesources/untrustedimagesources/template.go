package untrustedimagesurces

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
	"golang.stackrox.io/kube-linter/pkg/templates/untrustedimagesources/internal/params"
)

const (
	templateKey = "untrusted-image-sources"
)

func init() {
	templates.Register(check.Template{
		HumanName:   "Untrusted image sources",
		Key:         templateKey,
		Description: "Flag applications running containers with untrusted container image source",
		SupportedObjectKinds: config.ObjectKindsDesc{
			ObjectKinds: []string{objectkinds.DeploymentLike},
		},
		Parameters:             params.ParamDescs,
		ParseAndValidateParams: params.ParseAndValidate,
		Instantiate: params.WrapInstantiateFunc(func(p params.Params) (check.Func, error) {
			return util.PerContainerCheck(func(container *v1.Container) []diagnostic.Diagnostic {

				// checking whether the container image has a trusted image source or not

				res := false

				for _, imageSource := range p.NotAllowedImageSources {
					res = strings.Contains(container.Image, imageSource)
				}
				if res {
					return []diagnostic.Diagnostic{{Message: fmt.Sprintf("container %q image doesn't have a valid image source", container.Name)}}
				}
				return nil
			}), nil
		}),
	})
}
