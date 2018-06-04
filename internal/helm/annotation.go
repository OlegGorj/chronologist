/*
Copyright 2018 The Chronologist Authors

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

package helm

import (
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"k8s.io/helm/pkg/proto/hapi/release"

	"github.com/hypnoglow/chronologist/internal/chronologist"
)

// AnnotationFromRelease makes a chronologist annotation from helm release.
// This function always returns the same annotation for the same release.
func AnnotationFromRelease(rel *release.Release) (chronologist.Annotation, error) {
	// We are using LastDeployed field because it is relative only to a specific
	// revision and not to the release itself.
	t, err := ptypes.Timestamp(rel.Info.LastDeployed)
	if err != nil {
		return chronologist.Annotation{}, errors.Wrap(err, "unserialize timestamp from proto")
	}
	t = t.Truncate(time.Second).UTC()

	rt := chronologist.ReleaseTypeRollout
	if strings.Contains(strings.ToLower(rel.Info.Description), "rollback") {
		rt = chronologist.ReleaseTypeRollback
	}

	return chronologist.Annotation{
		GrafanaID:        0,
		Time:             t,
		ReleaseType:      rt,
		ReleaseStatus:    rel.Info.Status.Code.String(),
		ReleaseName:      rel.Name,
		ReleaseRevision:  strconv.Itoa(int(rel.Version)),
		ReleaseNamespace: rel.Namespace,
	}, nil
}

// AnnotationFromRawRelease makes a chronologist annotation from raw helm
// release data. This function always returns the same annotation for the
// same release.
func AnnotationFromRawRelease(data string) (chronologist.Annotation, error) {
	rel, err := DecodeRelease(data)
	if err != nil {
		return chronologist.Annotation{}, errors.Wrap(err, "decode raw release data")
	}

	a, err := AnnotationFromRelease(rel)
	return a, errors.Wrap(err, "create annotation from helm release")
}
