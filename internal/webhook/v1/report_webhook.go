/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"context"
	"fmt"
	"strconv"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	examplev1 "example.com/my-project/api/v1"
)

// nolint:unused
// log is for logging in this package.
var reportlog = logf.Log.WithName("report-resource")

// SetupReportWebhookWithManager registers the webhook for Report in the manager.
func SetupReportWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).For(&examplev1.Report{}).
		WithValidator(&ReportCustomValidator{}).
		WithDefaulter(&ReportCustomDefaulter{}).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-example-example-com-v1-report,mutating=true,failurePolicy=fail,sideEffects=None,groups=example.example.com,resources=reports,verbs=create;update,versions=v1,name=mreport-v1.kb.io,admissionReviewVersions=v1

// ReportCustomDefaulter struct is responsible for setting default values on the custom resource of the
// Kind Report when those are created or updated.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as it is used only for temporary operations and does not need to be deeply copied.
type ReportCustomDefaulter struct {
	// TODO(user): Add more fields as needed for defaulting
}

var _ webhook.CustomDefaulter = &ReportCustomDefaulter{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the Kind Report.
func (d *ReportCustomDefaulter) Default(_ context.Context, obj runtime.Object) error {
	report, ok := obj.(*examplev1.Report)
	if !ok {
		return fmt.Errorf("expected an Report object but got %T", obj)
	}

	if report.Labels == nil {
		report.Labels = make(map[string]string)
	}

	report.Labels["example.example.com/priority"] = strconv.Itoa(report.Spec.Priority)
	report.Labels["example.example.com/deadline"] = strconv.Itoa(int(report.Spec.Deadline.Unix()))

	return nil
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// NOTE: The 'path' attribute must follow a specific pattern and should not be modified directly here.
// Modifying the path for an invalid path can cause API server errors; failing to locate the webhook.
// +kubebuilder:webhook:path=/validate-example-example-com-v1-report,mutating=false,failurePolicy=fail,sideEffects=None,groups=example.example.com,resources=reports,verbs=create;update,versions=v1,name=vreport-v1.kb.io,admissionReviewVersions=v1

// ReportCustomValidator struct is responsible for validating the Report resource
// when it is created, updated, or deleted.
//
// NOTE: The +kubebuilder:object:generate=false marker prevents controller-gen from generating DeepCopy methods,
// as this struct is used only for temporary operations and does not need to be deeply copied.
type ReportCustomValidator struct {
	// TODO(user): Add more fields as needed for validation
}

var _ webhook.CustomValidator = &ReportCustomValidator{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type Report.
func (v *ReportCustomValidator) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	report, ok := obj.(*examplev1.Report)
	if !ok {
		return nil, fmt.Errorf("expected a Report object but got %T", obj)
	}
	reportlog.Info("Validation for Report upon creation", "name", report.GetName())

	// TODO(user): fill in your validation logic upon object creation.

	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type Report.
func (v *ReportCustomValidator) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	report, ok := newObj.(*examplev1.Report)
	if !ok {
		return nil, fmt.Errorf("expected a Report object for the newObj but got %T", newObj)
	}
	reportlog.Info("Validation for Report upon update", "name", report.GetName())

	// TODO(user): fill in your validation logic upon object update.

	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type Report.
func (v *ReportCustomValidator) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	report, ok := obj.(*examplev1.Report)
	if !ok {
		return nil, fmt.Errorf("expected a Report object but got %T", obj)
	}
	reportlog.Info("Validation for Report upon deletion", "name", report.GetName())

	// TODO(user): fill in your validation logic upon object deletion.

	return nil, nil
}
