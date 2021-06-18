package untrustedimagetags

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"golang.stackrox.io/kube-linter/pkg/diagnostic"
	"golang.stackrox.io/kube-linter/pkg/lintcontext/mocks"
	"golang.stackrox.io/kube-linter/pkg/templates"
	"golang.stackrox.io/kube-linter/pkg/templates/untrustedimagetags/internal/params"

	appsv1 "k8s.io/api/apps/v1"
)

func TestReplicas(t *testing.T) {
	suite.Run(t, new(ReplicaTestSuite))
}

type ReplicaTestSuite struct {
	templates.TemplateTestSuite

	ctx *mocks.MockLintContext
}

func (s *ReplicaTestSuite) SetupTest() {
	s.Init(templateKey)
	s.ctx = mocks.NewMockContext()
}

func (s *ReplicaTestSuite) addDeploymentWithUntrustedImageTag(name string, image string) {
	s.ctx.AddMockDeployment(s.T(), name)
	s.ctx.ModifyDeployment(s.T(), name, func(deployment *appsv1.Deployment) {
		deployment.Spec.Template.Containers.Image = &image
	})
}

func (s *ReplicaTestSuite) TestUntrustedImageTags() {
	const (
		latestImageTagDepName        = "latest-image-tag"
	)
	s.addDeploymentWithReplicas(latestImageTagDepName, "nginx:latest")

	s.Validate(s.ctx, []templates.TestCase{
		{
			Param: params.Params{
				NotAllowedTags: ["latest"],
			},
			Diagnostics: map[string][]diagnostic.Diagnostic{
				twoReplicasDepName: {
					{Message: "object has container image tag 'latest' but it should be other than the list of allowed image tags"},
				},
			},
			ExpectInstantiationError: false,
		},
	})
}

//func (s *ReplicaTestSuite) TestAcceptableReplicas() {
//	const (
//		acceptableReplicasDepName = "acceptable-replicas"
//	)
//	s.addDeploymentWithReplicas(acceptableReplicasDepName, 3)
//
//	s.Validate(s.ctx, []templates.TestCase{
//		{
//			Param: params.Params{
//				MinReplicas: 3,
//			},
//			Diagnostics: map[string][]diagnostic.Diagnostic{
//				acceptableReplicasDepName: nil,
//			},
//			ExpectInstantiationError: false,
//		},
//	})
//}
