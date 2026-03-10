package controller

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	examplev1 "example.com/my-project/api/v1"
)

// ReportReconciler reconciles a Report object
type ReportReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=example.example.com,resources=reports,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.example.com,resources=reports/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=example.example.com,resources=reports/finalizers,verbs=update

// +kubebuilder:rbac:groups="",resources=events,verbs=create

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *ReportReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	report := examplev1.Report{}
	if err := r.Get(ctx, req.NamespacedName, &report); err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	email := examplev1.Email{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      report.Name,
			Namespace: report.Namespace,
		},
		Spec: examplev1.EmailSpec{
			ToAddress:   "kovacsricsi@gmail.com",
			FromName:    "Demo Application",
			FromAddress: "ricsi.kovacs@inspirnation.eu",
		},
	}

	switch {
	case report.Generation == 1:
		email.Name += "-created"
		email.Spec.Subject = fmt.Sprintf("Report has been created: %s/%s", req.Namespace, req.Name)
	case report.DeletionTimestamp.IsZero():
		email.Name += "-updated"
		email.Spec.Subject = fmt.Sprintf("Report has been updated: %s/%s", req.Namespace, req.Name)
	case !controllerutil.ContainsFinalizer(&report, "example.example.com/finalizer"):
		return ctrl.Result{}, nil
	default:
		email.Name += "-deleted"
		email.Spec.Subject = fmt.Sprintf("Report has been deleted: %s/%s", req.Namespace, req.Name)

		controllerutil.RemoveFinalizer(&report, "example.example.com/finalizer")
		if err := r.Update(ctx, &report); err != nil {
			return ctrl.Result{}, err
		}
	}

	if err := r.Create(ctx, &email); err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ReportReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1.Report{}).
		Named("report").
		Complete(r)
}
