package kubernetes

import (
	"context"

	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/ships"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func InitDeployment(ship ships.Ship) error {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: ship.Name,
			Labels: map[string]string{
				"ship": string(ship.Frn),
			},
		},
		Spec: appsv1.DeploymentSpec{
			// 	Replicas: ,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"ship": string(ship.Frn),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"ship": string(ship.Frn),
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "fleet-agent",
							Image: "tgs266/fleet-agent",
							Env: []corev1.EnvVar{
								{
									Name:  "SHIP",
									Value: string(ship.Name),
								},
							},
						},
					},
				},
			},
		},
	}
	_, err := clientSet.clientSet.AppsV1().Deployments(ship.Namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return initService(ship)
}

// this allows the system to call localhost:8080 on minikube (INSTALLED ON SAME SYSTEM)
func initService(ship ships.Ship) error {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: ship.Name,
			Labels: map[string]string{
				"ship": string(ship.Frn),
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Ports: []corev1.ServicePort{
				{
					Port:       8080,
					TargetPort: intstr.FromInt(8080),
				},
			},
		},
	}

	endpoint := &corev1.Endpoints{
		ObjectMeta: metav1.ObjectMeta{
			Name: ship.Name,
			Labels: map[string]string{
				"ship": string(ship.Frn),
			},
		},
		Subsets: []corev1.EndpointSubset{
			{
				Addresses: []corev1.EndpointAddress{
					{
						IP: "192.168.65.2", // IP of host.minikube.internal (found through minikube ssh)
					},
				},
				Ports: []corev1.EndpointPort{
					{
						Port: 8080,
					},
				},
			},
		},
	}

	_, err := clientSet.clientSet.CoreV1().Services(ship.Namespace).Create(context.TODO(), service, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	_, err = clientSet.clientSet.CoreV1().Endpoints(ship.Namespace).Create(context.TODO(), endpoint, metav1.CreateOptions{})
	return err
}
