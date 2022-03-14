/*
Copyright 2021.

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

package v1alpha1

import (
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"os"
	"path/filepath"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"time"
)

const (
	WebhookCertDir  = "/apiserver.local.config/certificates"
	WebhookCertName = "apiserver.crt"
	WebhookKeyName  = "apiserver.key"
)

//minimal time durations allowed
const (
	MinDurPeerApiServerTimeout = "10ms"
	MinDurApiServerTimeout = "10ms"
	MinDurPeerDialTimeout = "10ms"
	MinDurPeerRequestTimeout = "10ms"
	MinDurApiCheckInterval = "1s"
	MinDurPeerUpdateInterval = "10s"
)

const (
	ErrPeerApiServerTimeout = "PeerApiServerTimeout " + MinDurPeerApiServerTimeout
	ErrApiServerTimeout     = "ApiServerTimeout " + MinDurApiServerTimeout
	ErrPeerDialTimeout      = "PeerDialTimeout " + MinDurPeerDialTimeout
	ErrPeerRequestTimeout   = "PeerRequestTimeout " + MinDurPeerRequestTimeout
	ErrApiCheckInterval     = "ApiCheckInterval can't be less than " + MinDurApiCheckInterval
	ErrPeerUpdateInterval   = "PeerUpdateInterval can't be less than " + MinDurPeerUpdateInterval
)

// log is for logging in this package.
var poisonpillconfiglog = logf.Log.WithName("poisonpillconfig-resource")

func (r *PoisonPillConfig) SetupWebhookWithManager(mgr ctrl.Manager) error {

	// check if OLM injected certs
	certs := []string{filepath.Join(WebhookCertDir, WebhookCertName), filepath.Join(WebhookCertDir, WebhookKeyName)}
	certsInjected := true
	for _, fname := range certs {
		if _, err := os.Stat(fname); err != nil {
			certsInjected = false
			break
		}
	}
	if certsInjected {
		server := mgr.GetWebhookServer()
		server.CertDir = WebhookCertDir
		server.CertName = WebhookCertName
		server.KeyName = WebhookKeyName
	} else {
		poisonpillconfiglog.Info("OLM injected certs for webhooks not found")
	}
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-poison-pill-medik8s-io-v1alpha1-poisonpillconfig,mutating=false,failurePolicy=fail,sideEffects=None,groups=poison-pill.medik8s.io,resources=poisonpillconfigs,verbs=create;update,versions=v1alpha1,name=vpoisonpillconfig.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &PoisonPillConfig{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *PoisonPillConfig) ValidateCreate() error {
	poisonpillconfiglog.Info("validate create", "name", r.Name)

	return r.ValidateTimes()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *PoisonPillConfig) ValidateUpdate(old runtime.Object) error {
	poisonpillconfiglog.Info("validate update", "name", r.Name)

	return r.ValidateTimes()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *PoisonPillConfig) ValidateDelete() error {
	poisonpillconfiglog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

// ValidateTimes validates each time field in the PoisonPillConfig CR doesn't go below the minimum time
// that was defined to it
func (r *PoisonPillConfig) ValidateTimes() error {
	peerApiServerTimeout := r.Spec.PeerApiServerTimeout.Milliseconds()
	apiServerTimeout := r.Spec.ApiServerTimeout.Milliseconds()
	peerDialTimeout := r.Spec.PeerDialTimeout.Milliseconds()
	peerRequestTimeout := r.Spec.PeerRequestTimeout.Milliseconds()
	apiCheckInterval := r.Spec.ApiCheckInterval.Milliseconds()
	peerUpdateInterval := r.Spec.PeerUpdateInterval.Milliseconds()
	if peerApiServerTimeout < toMS(MinDurPeerApiServerTimeout) {
		return LogAndReturnErr(ErrPeerApiServerTimeout, peerApiServerTimeout)
	} else if apiServerTimeout < toMS(MinDurApiServerTimeout) {
		return LogAndReturnErr(ErrApiServerTimeout, apiServerTimeout)
	} else if peerDialTimeout < toMS(MinDurPeerDialTimeout) {
		return LogAndReturnErr(ErrPeerDialTimeout, peerDialTimeout)
	} else if peerRequestTimeout < toMS(MinDurPeerRequestTimeout) {
		return LogAndReturnErr(ErrPeerRequestTimeout, peerRequestTimeout)
	} else if apiCheckInterval < toMS(MinDurApiCheckInterval) {
		return LogAndReturnErr(ErrApiCheckInterval, apiCheckInterval)
	} else if peerUpdateInterval < toMS(MinDurPeerUpdateInterval) {
		return LogAndReturnErr(ErrPeerUpdateInterval, peerUpdateInterval)
	}
	return nil
}

// LogAndReturnErr logs the time error with the inputTime as value for the user to see what was inserted
// and then returns the error.
func LogAndReturnErr(errMessage string, inputTime int64) error {
	err := fmt.Errorf(errMessage)
	poisonpillconfiglog.Error(err, errMessage, "time given (in milliseconds) was:", inputTime)
	return err
}

func toMS(value string) int64 {
	d, err := time.ParseDuration(value)
	if err != nil {
		//todo return error!
	}
	return d.Milliseconds()
}