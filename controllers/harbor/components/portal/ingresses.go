package portal

import (
	"context"
	"net/url"
	"strings"

	"github.com/goharbor/harbor-operator/pkg/ingress"

	"github.com/pkg/errors"
	netv1 "k8s.io/api/networking/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	goharborv1alpha1 "github.com/goharbor/harbor-operator/api/v1alpha1"
	"github.com/goharbor/harbor-operator/pkg/factories/application"
)

func (p *Portal) GetIngresses(ctx context.Context) []*netv1.Ingress { // nolint:funlen
	operatorName := application.GetName(ctx)
	harborName := p.harbor.Name

	u, err := url.Parse(p.harbor.Spec.PublicURL)
	if err != nil {
		panic(errors.Wrap(err, "invalid url"))
	}

	host := strings.SplitN(u.Host, ":", 1) // nolint:mnd

	var tls []netv1.IngressTLS
	if u.Scheme == "https" {
		tls = []netv1.IngressTLS{
			{
				SecretName: p.harbor.Spec.TLSSecretName,
				Hosts: []string{
					host[0],
				},
			},
		}
	}

	// Add annotations for cert-manager awareness
	annotations := ingress.GenerateIngressCertAnnotations(p.harbor.Spec)

	return []*netv1.Ingress{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:      p.harbor.NormalizeComponentName(goharborv1alpha1.PortalName),
				Namespace: p.harbor.Namespace,
				Labels: map[string]string{
					"app":      goharborv1alpha1.PortalName,
					"harbor":   harborName,
					"operator": operatorName,
				},
				Annotations: annotations,
			},
			Spec: netv1.IngressSpec{
				TLS: tls,
				Rules: []netv1.IngressRule{
					{
						Host: host[0],
						IngressRuleValue: netv1.IngressRuleValue{
							HTTP: &netv1.HTTPIngressRuleValue{
								Paths: []netv1.HTTPIngressPath{
									{
										Path: "/",
										Backend: netv1.IngressBackend{
											ServiceName: p.harbor.NormalizeComponentName(goharborv1alpha1.PortalName),
											ServicePort: intstr.FromInt(PublicPort),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
