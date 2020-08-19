package v1

import (
	"context"

	kutil "kmodules.xyz/client-go"

	promapi "github.com/coreos/prometheus-operator/pkg/apis/monitoring/v1"
	prom "github.com/coreos/prometheus-operator/pkg/client/versioned/typed/monitoring/v1"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/golang/glog"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

var json = jsoniter.ConfigFastest

func CreateOrPatchServiceMonitor(ctx context.Context, c prom.MonitoringV1Interface, meta metav1.ObjectMeta, transform func(monitor *promapi.ServiceMonitor) *promapi.ServiceMonitor, opts metav1.PatchOptions) (*promapi.ServiceMonitor, kutil.VerbType, error) {
	cur, err := c.ServiceMonitors(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		glog.V(3).Infof("Creating ServiceMonitor %s/%s.", meta.Namespace, meta.Name)
		out, err := c.ServiceMonitors(meta.Namespace).Create(ctx, transform(&promapi.ServiceMonitor{
			TypeMeta: metav1.TypeMeta{
				Kind:       promapi.PrometheusesKind,
				APIVersion: promapi.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}), metav1.CreateOptions{
			DryRun:       opts.DryRun,
			FieldManager: opts.FieldManager,
		})
		return out, kutil.VerbCreated, err
	} else if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	return PatchServiceMonitor(ctx, c, cur, transform, opts)
}

func PatchServiceMonitor(ctx context.Context, c prom.MonitoringV1Interface, cur *promapi.ServiceMonitor, transform func(monitor *promapi.ServiceMonitor) *promapi.ServiceMonitor, opts metav1.PatchOptions) (*promapi.ServiceMonitor, kutil.VerbType, error) {
	return PatchServiceMonitorObject(ctx, c, cur, transform(cur.DeepCopy()), opts)
}

func PatchServiceMonitorObject(ctx context.Context, c prom.MonitoringV1Interface, cur, mod *promapi.ServiceMonitor, opts metav1.PatchOptions) (*promapi.ServiceMonitor, kutil.VerbType, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	modJson, err := json.Marshal(mod)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	patch, err := jsonpatch.CreateMergePatch(curJson, modJson)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, kutil.VerbUnchanged, nil
	}
	glog.V(3).Infof("Patching ServiceMonitor %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.ServiceMonitors(cur.Namespace).Patch(ctx, cur.Name, types.MergePatchType, patch, opts)
	return out, kutil.VerbPatched, err
}

func TryUpdateMonitorObject(ctx context.Context, c prom.MonitoringV1Interface, meta metav1.ObjectMeta, transform func(monitor *promapi.ServiceMonitor) *promapi.ServiceMonitor, opts metav1.UpdateOptions) (result *promapi.ServiceMonitor, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.ServiceMonitors(meta.Namespace).Get(ctx, meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.ServiceMonitors(cur.Namespace).Update(ctx, transform(cur.DeepCopy()), opts)
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update ServiceMonitor %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = errors.Errorf("failed to update ServiceMonitor %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}
